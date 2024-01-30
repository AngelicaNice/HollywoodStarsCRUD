package mq

import (
	"context"
	"errors"

	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type AuditPublisher struct {
	amqpChan *amqp.Channel
	qname    string
}

func CreateMQConnection(url string) (*amqp.Connection, error) {
	return amqp.Dial(url)
}

func NewAuditPublisher(cfg *config.Config, mqConn *amqp.Connection) (*AuditPublisher, error) {
	amqpChan, err := mqConn.Channel()
	if err != nil {
		return nil, errors.New("failed to open a channel")
	}

	return &AuditPublisher{amqpChan: amqpChan, qname: cfg.MQ.Name}, nil
}

func (ap *AuditPublisher) CloseChan() {
	if err := ap.amqpChan.Close(); err != nil {
		log.WithField("rabbitmq", "failed to close chan")
	}
}

func (ap *AuditPublisher) Publish(ctx context.Context, body []byte) error {
	err := ap.amqpChan.PublishWithContext(
		ctx,
		"",
		ap.qname,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	if err != nil {
		log.WithField("rabbitmq", "failed to publish log")
	}

	return err
}
