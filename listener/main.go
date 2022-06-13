package main

import (
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// connect to rabbitmq
	mqconn, err := connect()
	if err != nil {
		log.Panic(err)
	}
	defer mqconn.Close()
	log.Println("Successfully connected to RabbitMQ")

	// listen for msgs

	// create consumer

	// watch the queue and consume msfs
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = time.Second * 1
	var conn *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			log.Println("Trying to connect to rabbitmq...")
			counts++
		} else {
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
