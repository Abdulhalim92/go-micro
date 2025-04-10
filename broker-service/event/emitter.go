package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

// Emitter - это структура, которая отвечает за отправку событий в RabbitMQ
type Emitter struct {
	connection *amqp.Connection
}

// NewEventEmitter - это конструктор для создания нового экземпляра Emitter
func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}

// Setup - это метод, который устанавливает соединение с RabbitMQ и настраивает обменник
func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %s", err)
		return err
	}
	defer channel.Close()

	return declareExchange(channel)
}

// Push - это метод, который отправляет событие в RabbitMQ
func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %s", err)
		return err
	}
	defer channel.Close()

	log.Println("Pushing event to RabbitMQ")

	err = channel.Publish(
		"log_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		log.Printf("Failed to publish event: %s", err)
		return err
	}

	return nil
}
