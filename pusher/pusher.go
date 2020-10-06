package pusher

import (
	pushnotifications "github.com/pusher/push-notifications-go"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	InstanceId string
	SecretKey  string
}

type Client struct {
	pushnotifications.PushNotifications
}

func Init(config Config) Client {
	client, err := pushnotifications.New(config.InstanceId, config.SecretKey)
	if err != nil {
		log.Error(err.Error())
	}
	return Client{client}
}

func (client *Client) PublishNotification(users []string, publishRequest map[string]interface{}) (string, error) {
	publishId, err := client.PublishToUsers(users, publishRequest)
	if err != nil {
		return "", err
	}
	return publishId, nil
}
