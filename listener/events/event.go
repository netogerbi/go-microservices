package events

import amqp "github.com/rabbitmq/amqp091-go"

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable?
		false,        // auto-detect?
		false,        // internal?
		false,        // no-wait?
		nil,          // arguments?
	)
}

func declareQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name
		true,  // durable?
		false, // auto-delete?
		true,  // exclusive?
		false, // no-wait?
		nil,   // arguments
	)
}
