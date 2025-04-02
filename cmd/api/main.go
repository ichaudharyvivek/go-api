package main

import (
	"log"

	"example.com/go-api/internal/env"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	app := &application{
		config: cfg,
	}

	r := app.mount()
	log.Fatal(app.run(r))
}
