package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	RabbitMQConn *amqp.Connection
}

func main() {
	mqconn, err := connect()
	if err != nil {
		log.Panic(err)
	}
	defer mqconn.Close()

	app := Config{
		RabbitMQConn: mqconn,
	}

	log.Printf("Starting broker service on port %s", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = time.Second * 1
	var conn *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Println("Trying to connect to rabbitmq...")
			counts++
		} else {
			log.Println("Successfully connected to RabbitMQ")
			conn = c
			break
		}

		if counts > 5 {
			log.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2))
		log.Printf("Backing off %d seconds until try again...", backOff)
		time.Sleep(backOff)
		continue
	}

	return conn, nil
}
