package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/goapi/internal/config"
	handler "example.com/goapi/internal/handler/http"
)

func main() {
	c := config.New()
	r := handler.NewRouter()
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Server.Port),
		Handler:      r,
		ReadTimeout:  c.Server.TimeoutRead,
		WriteTimeout: c.Server.TimeoutWrite,
		IdleTimeout:  c.Server.TimeoutIdle,
	}

	log.Println("Starting server at port " + s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server startup failed")
	}
}
