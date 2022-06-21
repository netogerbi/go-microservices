package main

import (
	"listener/events"
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

	// listen for msgs
	log.Println("Start listening and consumming messages from RqbbitMQ...")

	// create consumer
	c, err := events.NewConsumer(mqconn)
	if err != nil {
		panic(err)
	}

	// watch the queue and consume msfs
	if err = c.Listen([]string{"log.INFO", "log.ERROR", "log.WARN"}); err != nil {
		log.Panicln(err)
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
