package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

// Consumer - структура для потребителя сообщений
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// NewConsumer - функция для создания нового потребителя сообщений
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// Setup - функция для настройки потребителя сообщений
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

// Payload - структура для полезной нагрузки сообщения
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		log.Printf("Failed to create RabbitMQ channel: %s", err)
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		log.Printf("Failed to declare a queue: %s", err)
		return err
	}

	for _, topic := range topics {
		err := ch.QueueBind(
			q.Name,
			topic,
			"log_topic",
			false,
			nil,
		)
		if err != nil {
			log.Printf("Failed to bind a queue: %s", err)
			return err
		}
	}

	messages, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to register a consumer: %s", err)
		return err
	}

	forever := make(chan bool)

	go func() {
		for msg := range messages {
			// process message
			log.Printf("Received message: %s", msg.Body)

			var payload Payload
			err := json.Unmarshal(msg.Body, &payload)
			if err != nil {
				log.Printf("Error unmarshalling message from RabbitMQ: %s", err)
				continue
			}

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [log_topics, %s]\n", q.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Printf("Error logging event: %s", err)
		}
	case "auth":
	// authenticate

	// you can have as many cases as you want, as long as you write the logic

	default:
		err := logEvent(payload)
		if err != nil {
			log.Printf("Error logging event: %s", err)
		}
	}
}

func logEvent(payload Payload) error {
	jsonData, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return fmt.Errorf("[logEevent] -> error marshalling payload: %s", err)
	}

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest(http.MethodPost, logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("[logEevent] -> error creating request: %s", err)
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("[logEevent] -> error executing request: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("[logEevent] -> error returned status code %d", response.StatusCode)
	}

	return nil
}
