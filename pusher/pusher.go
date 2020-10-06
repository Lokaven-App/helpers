package pusher

import "github.com/pusher/pusher-http-go"

type Config struct {
	AppID   string
	Key     string
	Secret  string
	Cluster string
}

type Client struct {
	*pusher.Client
}

func Init(config Config) Client {
	client := pusher.Client{
		AppID:   config.AppID,
		Key:     config.Key,
		Secret:  config.Secret,
		Cluster: config.Cluster,
	}
	return Client{&client}
}

func (client *Client) PushTrigger(data []byte, event, channel string) error {
	if err := client.Trigger(channel, event, data); err != nil {
		return err
	}
	return nil
}
