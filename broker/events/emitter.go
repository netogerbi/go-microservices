package events

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	conn *amqp.Connection
}

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	e := Emitter{
		conn: conn,
	}

	if err := e.setup(); err != nil {
		return Emitter{}, err
	}

	return e, nil
}

func (e *Emitter) setup() error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return declareExchange(ch)
}

func (e *Emitter) Push(event string, severity string) error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	log.Println("Publishing to an exchange...")

	err = ch.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)

	if err != nil {
		return err
	}

	return nil
}
