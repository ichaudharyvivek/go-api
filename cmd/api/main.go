package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"example.com/goapi/docs"
	"example.com/goapi/internal/config"
	"example.com/goapi/internal/config/env"
	"example.com/goapi/internal/database"
	"example.com/goapi/internal/database/cache"
	"example.com/goapi/internal/router"
	"example.com/goapi/internal/utils/logger"
	"example.com/goapi/internal/utils/validator"
)

var appEnv = env.GetString("APP_ENV", "dev")
var isProd = appEnv == "prod"

func main() {
	logger.Setup(isProd)

	c := config.New()
	v := validator.New()
	db, _ := database.NewDB(c)
	rd := cache.NewClient(c)

	r := router.NewRouter(db, v, rd)
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Server.Port),
		Handler:      r,
		ReadTimeout:  c.Server.TimeoutRead,
		WriteTimeout: c.Server.TimeoutWrite,
		IdleTimeout:  c.Server.TimeoutIdle,
	}

	swaggerInit()

	log.Println("Connecting to redis client at " + rd.Options().Addr)
	err := rd.Ping(context.Background())
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Starting server at port" + s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server startup failed!!")
	}
}

func swaggerInit() {
	docs.SwaggerInfo.Version = "0.0.1"
	docs.SwaggerInfo.Host = env.GetString("APP_HOST", "localhost:8080")
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Title = "GopherSocial Media App API"
	docs.SwaggerInfo.Description = "My awesome API"
}
