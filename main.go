package main

import (
	"fmt"
	"log"

	conf "parkDB/config"
	"parkDB/delivery"
	"parkDB/middleware"
	"parkDB/repository"
	"parkDB/usecase"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/penglongli/gin-metrics/ginmetrics"
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

	p := ginprom.New(
		ginprom.Engine(myRouter),
		ginprom.Subsystem("gin"),
		ginprom.Path(conf.MetricsPath),
	)
	p.Use(myRouter)
	myRouter.Use(p.Instrument())
	myRouter.Use(middleware.IncCounter)

	metics := ginmetrics.GetMonitor()
	metics.SetMetricPath("/dev" + conf.MetricsPath)
	metics.SetSlowTime(5)
	metics.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	metics.Use(myRouter)

	err = myRouter.Run(conf.ServerPort)
	if err != nil {
		log.Println("can't serve", err)
	}
}
