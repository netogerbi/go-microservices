package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct {
	Mailer Mail
}

const (
	webPort = "80"
)

func main() {
	app := Config{
		Mailer: NewMailer(),
	}

	if err := app.serve(); err != nil {
		log.Panic(err)
	}
}

func (app *Config) serve() error {
	log.Printf("Starting mail service...")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Println("Error starting server: ", err)
		return err
	}

	return nil
}
