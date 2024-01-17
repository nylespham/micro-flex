package main

import (
	"fmt"
	"log"
	"net/http"
	// "os"
	// "strconv"
)

type Config struct {
	Mailer Mail
}

const port = "4100"

func main() {
	app := Config{
		Mailer: createMail(),
	}

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

func createMail() Mail {
	// port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	// m := Mail{
	// 	Domain:      os.Getenv("MAIL_DOMAIN"),
	// 	Host:        os.Getenv("MAIL_HOST"),
	// 	Port:        port,
	// 	Username:    os.Getenv("MAIL_USERNAME"),
	// 	Password:    os.Getenv("MAIL_PASSWORD"),
	// 	Encryption:  os.Getenv("MAIL_ENCRYPTION"),
	// 	FromName:    os.Getenv("FROM_NAME"),
	// 	FromAddress: os.Getenv("FROM_ADDRESS"),
	// }
	m := Mail{
		// Domain:      "localhost",
		Host:        "35.198.205.242",
		Port:        1025,
		Username:    "",
		Password:    "",
		Encryption:  "none",
		FromName:    "Test",
		FromAddress: "user@example.com",
	}
	return m
}
