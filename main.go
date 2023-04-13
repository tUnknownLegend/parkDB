package main

import (
	"fmt"
	"log"

	conf "parkDB/config"
	"parkDB/delivery"
	"parkDB/middleware"
	"parkDB/repository"
	"parkDB/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type CustomAspect struct {
	CustomValue int
}

func (a *CustomAspect) GetStats() interface{} {
	return a.CustomValue
}

func (a *CustomAspect) Name() string {
	return "Custom"
}

func (a *CustomAspect) InRoot() bool {
	return false
}

func main() {
	myRouter := gin.New()

	dbConf := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		conf.DBhost, conf.DBuser, conf.DBpwd, conf.DBname, conf.DBport)
	connStr, err := pgx.ParseConnectionString(dbConf)
	if err != nil {
		log.Println(err)
	}
	db, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connStr,
		MaxConnections: 200,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	})
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	userStore := repository.NewUserRepository(db)
	forumStore := repository.NewForumRepository(db)
	postStore := repository.NewPostRepository(db)
	serviceStore := repository.NewServiceRepository(db)
	threadStore := repository.NewThreadRepository(db)

	userUsecase := usecase.NewUserUsecase(userStore)
	forumUsecase := usecase.NewForumUsecase(forumStore, threadStore, userStore)
	postUsecase := usecase.NewPostUsecase(postStore, userStore, threadStore, forumStore)
	serviceUsecase := usecase.NewServiceUsecase(serviceStore)
	threadUsecase := usecase.NewThreadUsecase(threadStore, postStore, userStore)

	routerGroup := myRouter.Group(conf.BasePath)
	delivery.NewUserHandler(routerGroup, conf.BaseUserPath, userUsecase)
	delivery.NewForumHandler(routerGroup, conf.BaseForumPath, forumUsecase)
	delivery.NewPostHandler(routerGroup, conf.BasePostPath, postUsecase)
	delivery.NewServiceHandler(routerGroup, conf.BaseServicePath, serviceUsecase)
	delivery.NewThreadHandler(routerGroup, conf.BaseThreadPath, threadUsecase)

	prometheus.Register(middleware.HitsCounter)
	myRouter.Use(middleware.IncCounter)

	myRouter.GET(conf.MetricsPath, func(c *gin.Context) {
		handler := promhttp.Handler()
		handler.ServeHTTP(c.Writer, c.Request)
	})

	err = myRouter.Run(conf.ServerPort)
	if err != nil {
		log.Println("can't serve", err)
	}
}
