package rabbitmq

import (

	"github.com/streadway/amqp"
)

func connectToRabbitMQ() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqps://aetsgqau:ud_T9i0rUzPfJp2heU_9XkVayu8PjTEg@lionfish.rmq.cloudamqp.com/aetsgqau")
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func publishMessage(conn *amqp.Connection, routingKey string, message string) error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	err = channel.Publish(
		"",          // exchange
		routingKey,  // routing key (queue name)
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
