package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = "8300"

type Config struct{}

func main() {
	app := Config{}

	log.Printf("Starting server on port %s\n", port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}
	err := server.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}
