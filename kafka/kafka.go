package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Lokaven-App/proto/notification"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Username string
	Password string
	Url      string
	Topic    string
	Group    string
}

type Message struct {
	config       kafka.WriterConfig
	Defaults     []byte
	Notification notification.Notification
}

type Header struct {
	Key   string
	Value string
}

type Code int32

const (
	MessageBody         Code = 0
	MessageNotification Code = 1
)

func Init(config Config) Message {
	dialer := kafka.Dialer{
		Timeout: 10 * time.Second,
		SASLMechanism: plain.Mechanism{
			Username: config.Username,
			Password: config.Password,
		},
	}

	msg := Message{}
	msg.config = kafka.WriterConfig{
		Brokers:     []string{config.Url},
		Topic:       config.Topic,
		Balancer:    &kafka.Hash{},
		Dialer:      &dialer,
		BatchSize:   1,
		MaxAttempts: 20,
		Async:       false,
	}

	return msg
}

func (msg *Message) Publish(ctx context.Context, header []Header, code Code) error {
	body, err := getBody(msg, &code)
	if err != nil {
		return err
	}

	headers := []kafka.Header{}
	for _, val := range header {
		headers = append(headers, kafka.Header{Key: val.Key, Value: []byte(val.Value)})
	}

	errPublish := kafka.NewWriter(msg.config).WriteMessages(ctx, kafka.Message{
		Value:   body,
		Headers: headers,
	})
	if errPublish != nil {
		return errPublish
	}
	return nil
}

func Consumer(config Config) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.Url},
		Dialer: &kafka.Dialer{
			SASLMechanism: plain.Mechanism{
				Username: config.Username,
				Password: config.Password,
			},
		},
		GroupID:  config.Group,
		Topic:    config.Topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return reader
}

func getBody(msg *Message, code *Code) (body []byte, err error) {
	switch *code {
	case 0:
		log.Infof("Publish message default with code : ", *code)
		body = msg.Defaults
	case 1:
		log.Infof("Publish message notification with code : ", *code)
		body, err = json.Marshal(&msg.Notification)
		if err != nil {
			return nil, err
		}
	default:
		//Action
	}
	return body, nil
}
