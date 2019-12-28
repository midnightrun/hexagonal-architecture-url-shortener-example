package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	h "github.com/midnightrun/hexagonal-architecture-url-shortener-example/api"
	ll "github.com/midnightrun/hexagonal-architecture-url-shortener-example/logger/logrus"
	zl "github.com/midnightrun/hexagonal-architecture-url-shortener-example/logger/zap"
	mr "github.com/midnightrun/hexagonal-architecture-url-shortener-example/repository/mongo"
	rr "github.com/midnightrun/hexagonal-architecture-url-shortener-example/repository/redis"
	"github.com/midnightrun/hexagonal-architecture-url-shortener-example/shortener"
)

func httpPort() string {
	port := "8000"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	return fmt.Sprintf(":%s", port)
}

func chooseRepo() shortener.RedirectRepository {
	switch strings.ToLower(os.Getenv("URL_DB")) {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")

		repo, err := rr.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}

		return repo
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongoDB := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))

		repo, err := mr.NewMongoRepository(mongoURL, mongoDB, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}

		return repo
	}

	return nil
}

func chooseLogger() shortener.Logger {
	config := shortener.Configuration{
		EnableConsole:     true,
		ConsoleLevel:      shortener.Debug,
		ConsoleJSONFormat: true,
		EnableFile:        true,
		FileLevel:         shortener.Info,
		FileJSONFormat:    true,
		FileLocation:      "log.log",
	}

	var logger shortener.Logger

	var err error

	switch strings.ToLower(os.Getenv("LOGGER")) {
	case "zap":
		logger, err = zl.NewZapLogger(config)
		if err != nil {
			return nil
		}
	case "logrus":
		logger, err = ll.NewLogrusLogger(config)
		if err != nil {
			return nil
		}
	}

	contextLogger := logger.WithFields(shortener.Fields{"key1": "value1"})

	return contextLogger
}

func main() {
	fmt.Println("Hexagonal URL Shortener")

	log := chooseLogger()
	repo := chooseRepo()
	service := shortener.NewRedirectService(repo)
	handler := h.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	errs := make(chan error, 2)

	go func() {
		log.WithFields().Infof()
		fmt.Println("Listening on port: ", httpPort())
		errs <- http.ListenAndServe(httpPort(), r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}
