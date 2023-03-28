package mbaas

import (
	"context"
	"errors"

	"outbox/config"
	"outbox/internal/application/modules/categories/models"
	"outbox/internal/infrastructure/log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MBaaS interface {
	Connect() error
	Publish(context.Context, models.Outbox) error
}

type mbaas struct {
	cfg        *config.AppConfig
	connection *amqp.Connection
}

func NewMbaas(cfg *config.AppConfig) MBaaS {
	return &mbaas{cfg: cfg}
}

func (m *mbaas) Connect() error {
	ctx := context.Background()
	log.For(ctx).Infof("connecting to RabbitMQ instance....")

	connection, err := amqp.Dial(m.cfg.RabbitUrl)
	if err != nil {
		log.For(ctx).Infof("error to connect RabbitMQ. %+v", err)
		return err
	}

	log.For(ctx).Infof("successfully connected to RabbitMQ instance!")
	m.connection = connection
	return nil
}

func (m *mbaas) canConnect() bool {
	if err := m.Connect(); err != nil {
		return false
	}
	return true
}

func (m *mbaas) Publish(ctx context.Context, msg models.Outbox) error {
	if m.connection.IsClosed() && !m.canConnect() {
		return errors.New("cannot connected to rabbit client")
	}

	channel, err := m.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	// publishing a message
	err = channel.PublishWithContext(context.Background(),
		"",                // exchange
		m.cfg.OutboxQueue, // key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        []byte(msg.Message),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
