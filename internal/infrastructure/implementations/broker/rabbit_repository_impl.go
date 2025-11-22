package broker

import (
	"log"

	"authService/internal/config"
	"authService/internal/domain/repositories"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitRepositoryImpl struct {
	cfg *config.Config
}

func (r *RabbitRepositoryImpl) NewConnection() *amqp.Connection {
	conn, err := amqp.Dial(r.cfg.RabbitMQ.RabbitMQUrl())
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
	}
	return conn
}

func NewRabbitRepositoryImpl(cfg *config.Config) repositories.RabbitRepository {
	return &RabbitRepositoryImpl{
		cfg: cfg,
	}
}

func (r *RabbitRepositoryImpl) CreateEmailMSG(email string) error {
	conn := r.NewConnection()
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()

	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer func() {
		if err := channel.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
	}()

	que, err := channel.QueueDeclare(r.cfg.BrokerConstants.EmailConfirm, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = channel.Publish("", que.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(email),
	})
	if err != nil {
		return err
	}
	return nil
}
