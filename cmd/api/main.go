package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/goapi/internal/config"
	database "example.com/goapi/internal/database"
	router "example.com/goapi/internal/router"
	"example.com/goapi/internal/utils/validator"
)

func main() {
	c := config.New()
	v := validator.New()
	db, _ := database.NewDB(c)

	r := router.NewRouter(db, v)
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Server.Port),
		Handler:      r,
		ReadTimeout:  c.Server.TimeoutRead,
		WriteTimeout: c.Server.TimeoutWrite,
		IdleTimeout:  c.Server.TimeoutIdle,
	}

	log.Println("Starting server at port " + s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server startup failed!!")
	}
}
