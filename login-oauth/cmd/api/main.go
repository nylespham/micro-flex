package main

import (
	"database/sql"
	"fmt"
	"log"
	"login-oauth/data"
	"net/http"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const port = "13780"

var count int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service on port: %s\n", port)

	// TODO connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Could not connect to database")
	}

	// set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := "host=localhost port=5432 dbname=postgres user=postgres password=P@ssw0rd sslmode=disable timezone=UTC connect_timeout=5"

	for {

		connection, err := openDB(dsn)

		if err != nil {
			log.Println("Postgres not available ...")
			count++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if count > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
