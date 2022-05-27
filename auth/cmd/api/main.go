package main

import (
	"auth/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var connectTries int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting Auth service...")

	//todo connect to db
	conn := connectToDB()
	if conn == nil {
		log.Panic("Could not connect to DB!")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	//setup config
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)

		if err != nil {
			log.Println("Connnecting to db...")
			connectTries++
		} else {
			log.Println("Connected to DB successfully!")
			return connection
		}

		if connectTries > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Waiting to try to connect again...")
		time.Sleep(time.Second * 2)
		continue
	}
}
