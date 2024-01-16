package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct {
}

const port = "4100"

func main() {
	app := Config{}

	log.Println("Starting server on port", port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
