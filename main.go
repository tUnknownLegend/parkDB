package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	conf "github.com/tUnknownLegend/parkDB/conf"

	"github.com/gorilla/mux"
)

func main() {
	myRouter := mux.NewRouter()
	urlDB := "postgres://" + os.Getenv("TEST_POSTGRES_USER") + ":" + os.Getenv("TEST_POSTGRES_PASSWORD") + "@" + os.Getenv("TEST_DATABASE_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("TEST_POSTGRES_DB")
	log.Println("conn: ", urlDB)
	db, err := sql.Open("pgx", urlDB)
	if err != nil {
		log.Println("could not connect to database")
	} else {
		log.Println("database is reachable")
	}
	defer db.Close()

	// userStore := repository.NewUserStore(db)
	// productStore := repository.NewProductStore(db)

	// userUsecase := usecase.NewUserUsecase(userStore, sessManager, mailManager)
	// productUsecase := usecase.NewProductUsecase(productStore, ordersManager, mailManager)

	// userHandler := deliv.NewUserHandler(userUsecase)
	// sessionHandler := deliv.NewSessionHandler(userUsecase)
	// productHandler := deliv.NewProductHandler(productUsecase, userUsecase)

	// orderHandler := deliv.NewOrderHandler(userHandler, productHandler)

	// userRouter := myRouter.PathPrefix("/api/v1/user").Subrouter()
	// cartRouter := myRouter.PathPrefix("/api/v1/cart").Subrouter()

	// myRouter.HandleFunc(conf.PathLogin, sessionHandler.Login).Methods(http.MethodPost, http.MethodOptions)
	// myRouter.HandleFunc(conf.PathLogOut, sessionHandler.Logout).Methods(http.MethodDelete, http.MethodOptions)
	// myRouter.HandleFunc(conf.PathSignUp, sessionHandler.SignUp).Methods(http.MethodPost, http.MethodOptions)
	// myRouter.HandleFunc(conf.PathSessions, sessionHandler.GetSession).Methods(http.MethodGet, http.MethodOptions)

	// myRouter.HandleFunc(conf.PathProductByID, productHandler.GetProductByID).Methods(http.MethodGet, http.MethodOptions)
	// myRouter.HandleFunc(conf.PathMain, productHandler.GetHomePage).Methods(http.MethodGet, http.MethodOptions)
	// myRouter.HandleFunc(conf.PathCategory, productHandler.GetProductsByCategory).Methods(http.MethodGet, http.MethodOptions)
	// myRouter.HandleFunc(conf.PathSeacrh, productHandler.GetProductsBySearch).Methods(http.MethodPost, http.MethodOptions)
	// myRouter.HandleFunc(conf.PathSuggestions, productHandler.GetSuggestions).Methods(http.MethodPost, http.MethodOptions)
	// myRouter.HandleFunc(conf.PathRecommendations, productHandler.GetRecommendations).Methods(http.MethodGet, http.MethodOptions)
	// myRouter.HandleFunc(conf.PathProductsWithDiscount, productHandler.GetProductsWithBiggestDiscount).Methods(http.MethodGet, http.MethodOptions)
	// myRouter.HandleFunc(conf.PathBestProductCategory, productHandler.GetBestProductInCategory).Methods(http.MethodGet, http.MethodOptions)
	// myRouter.HandleFunc(conf.PathRecalculateRatings, productHandler.RecalculateRatingsForInitscriptProducts).Methods(http.MethodPost, http.MethodOptions)

	// userRouter.HandleFunc(conf.PathProfile, userHandler.GetUser).Methods(http.MethodGet, http.MethodOptions)
	// userRouter.HandleFunc(conf.PathProfile, userHandler.ChangeProfile).Methods(http.MethodPost, http.MethodOptions)
	// userRouter.HandleFunc(conf.PathAvatar, userHandler.SetAvatar).Methods(http.MethodPost, http.MethodOptions)
	// userRouter.HandleFunc(conf.PathPassword, userHandler.ChangePassword).Methods(http.MethodPost, http.MethodOptions)
	// userRouter.HandleFunc(conf.PathFavorites, productHandler.GetFavorites).Methods(http.MethodGet, http.MethodOptions)
	// userRouter.HandleFunc(conf.PathInsertIntoFavorites, productHandler.InsertItemIntoFavorites).Methods(http.MethodPost, http.MethodOptions)
	// userRouter.HandleFunc(conf.PathDeleteFromFavorites, productHandler.DeleteItemFromFavorites).Methods(http.MethodPost, http.MethodOptions)

	// myRouter.HandleFunc(conf.PathComments, orderHandler.GetComments).Methods(http.MethodGet, http.MethodOptions)
	// userRouter.HandleFunc(conf.PathMakeComment, orderHandler.CreateComment).Methods(http.MethodPost, http.MethodOptions)

	// cartRouter.HandleFunc("", orderHandler.GetCart).Methods(http.MethodGet, http.MethodOptions)
	// cartRouter.HandleFunc("", orderHandler.UpdateCart).Methods(http.MethodPost, http.MethodOptions)
	// cartRouter.HandleFunc(conf.PathAddItemToCart, orderHandler.AddItemToCart).Methods(http.MethodPost, http.MethodOptions)
	// cartRouter.HandleFunc(conf.PathDeleteItemFromCart, orderHandler.DeleteItemFromCart).Methods(http.MethodPost, http.MethodOptions)
	// cartRouter.HandleFunc(conf.PathMakeOrder, orderHandler.MakeOrder).Methods(http.MethodPost, http.MethodOptions)
	// cartRouter.HandleFunc(conf.PathGetOrders, orderHandler.GetOrders).Methods(http.MethodGet, http.MethodOptions)
	// cartRouter.HandleFunc(conf.PathPromo, orderHandler.SetPromocode).Methods(http.MethodPost, http.MethodOptions)
	// cartRouter.HandleFunc(conf.PathChangeOrderStatus, orderHandler.ChangeOrderStatus).Methods(http.MethodPost, http.MethodOptions)

	// myRouter.PathPrefix(conf.PathDocs).Handler(httpSwagger.WrapHandler)
	// myRouter.Use(loggingAndCORSHeadersMiddleware)

	// instrumentation := muxprom.NewDefaultInstrumentation()
	// myRouter.Use(instrumentation.Middleware)
	// myRouter.Path("/metrics").Handler(promhttp.Handler())

	// amw := deliv.NewAuthMiddleware(userUsecase)

	// userRouter.Use(amw.CheckAuthMiddleware)
	// cartRouter.Use(amw.CheckAuthMiddleware)

	err = http.ListenAndServe(conf.Port, myRouter)
	if err != nil {
		log.Println("can't serve", err)
	}
}
