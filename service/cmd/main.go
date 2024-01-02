package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/Nahbox/streamed-order-viewer/service/internal/cached"
	"github.com/Nahbox/streamed-order-viewer/service/internal/config"
	"github.com/Nahbox/streamed-order-viewer/service/internal/db"
	"github.com/Nahbox/streamed-order-viewer/service/internal/gui"
	"github.com/Nahbox/streamed-order-viewer/service/internal/handler"
	"github.com/Nahbox/streamed-order-viewer/service/internal/handler/order"
	"github.com/Nahbox/streamed-order-viewer/service/internal/nats-subscriber"
)

func main() {
	godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatal("read config from env", err)
	}

	// Подключение к базе данных
	database, err := db.Initialize(cfg.PgConfig)
	if err != nil {
		log.Fatal("init db", err)
	}
	defer database.Close()

	// Cache init
	cacheManager := cached.NewCache(1*time.Hour, 24*time.Hour)

	wg.Add(1)
	go func(cacheManager *cached.CacheManager, err error) {
		defer wg.Done()
		data, err := database.GetAllOrders()
		if err != nil {
			return
		}

		for _, val := range data {
			key := val.Order.Id
			cacheManager.Cache.Set(strconv.Itoa(key), val, -1)
		}

		log.Println("cache ready")
	}(cacheManager, err)

	orderHandler := order.NewHandler(database, cacheManager)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = nats_subscriber.NatsStreamingSubscribe(ctx, orderHandler)
		if err != nil {
			log.Fatal("func: NatsStreamingSubscribe ", err)
		}
	}()

	router := chi.NewRouter()
	router.MethodNotAllowed(handler.MethodNotAllowed)
	router.NotFound(handler.NotFound)

	router.Get("/", gui.FindOrderByIDPage(""))
	router.Get("/order", gui.FindOrderByID(orderHandler))

	addr := fmt.Sprintf(":%d", cfg.AppPort)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.WithError(err).Fatal("run http server")
		}
	}()
	defer Stop(server)

	log.Infof("started API server on %s", addr)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch

	cancel()
	wg.Wait()

	log.Infoln("stopping API server")
}

func Stop(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Error("shutdown server")
		os.Exit(1)
	}
}
