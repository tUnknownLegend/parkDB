package main

import (
	"fmt"
	"log"
	"net/http"

	conf "parkDB/config"
	"parkDB/delivery"
	"parkDB/middleware"
	"parkDB/repository"
	"parkDB/usecase"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// metrics

	// Create non-global registry.
	registry := prometheus.NewRegistry()

	// Add go runtime metrics and process collectors.
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	// Expose /metrics HTTP endpoint using the created custom registry.
	http.Handle(
		conf.MetricsPath,
		middleware.New(
			registry, nil).
			WrapHandler(conf.MetricsPath, promhttp.HandlerFor(
				registry,
				promhttp.HandlerOpts{}),
			))

	log.Fatalln(http.ListenAndServe(":5050", nil))

	// end metrics

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
	myRouter.Use(p.Instrument())

	metics := ginmetrics.GetMonitor()
	metics.SetMetricPath(conf.MetricsPath)
	metics.SetSlowTime(5)
	metics.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	metics.Use(myRouter)

	err = myRouter.Run(conf.ServerPort)
	if err != nil {
		log.Println("can't serve", err)
	}
}
