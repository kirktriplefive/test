package main

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kirktriplefive/test/pkg/cache"
	"github.com/kirktriplefive/test/pkg/nats_sub"
	"github.com/kirktriplefive/test/pkg/repository"
	"github.com/kirktriplefive/test/pkg/service"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}




func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err!=nil{
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err!=nil{
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host: viper.GetString("db.host"),
		Port: viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName: viper.GetString("db.dbname"),
		SSLMode: viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil{
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos:= repository.NewRepository(db)
    cache := cache.NewCache()

	services := service.NewService(repos, cache)
    
    if err := services.CacheFromPQ();err != nil {
		logrus.Fatal(err, "Failed to recover cache")
	}


    nc := nats_sub.NewSubscriber(nats_sub.Client{
        M:         &sync.Mutex{},
		Host:      stan.DefaultNatsURL,
		ClusterID: viper.GetString("client.ClusterID"),
		ClientID:  viper.GetString("client.ClientID"),
		Subject:   viper.GetString("client.Subject"),
		Service:   services,
    })
    err = nc.ConnectToStan()
    if err != nil {
		logrus.Fatal(err)
	} else {
		logrus.Info("connection successful")
	}
    defer nc.Close()

    listenErr := make(chan error, 1)
	osSignals := make(chan os.Signal, 1)

	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		rout := gin.New()
		rout.LoadHTMLGlob("web-interface/*")
		api := rout.Group("/api")
	{
		order:=api.Group("/order")
		{
			//order.POST("/", h.createOrder)
			order.GET("/:id", services.GetOrderById)
		}

	}
	http.Handle("/", rout)
	listenErr <- rout.Run()
	}()

	logrus.Info("service started")

    select {
	case err = <-listenErr:
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
	case <-osSignals:
		services.Close()
	}

}

func cors(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type,Accept")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "OPTIONS,POST,GET")
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		if string(ctx.Method()) == http.MethodOptions {
			ctx.Response.SetStatusCode(200)
			return
		}
		next(ctx)
	}
}