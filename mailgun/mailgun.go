package mailgun

import (
	"context"
	"time"

	_mailgun "github.com/mailgun/mailgun-go/v3"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Domain string
	Key    string
}

type Mailer struct {
	*_mailgun.MailgunImpl
}

func Mailgun(config Config) *Mailer {
	mg := _mailgun.NewMailgun(config.Domain, config.Key)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, err := mg.GetDomainTracking(ctx, mg.Domain())
	if err != nil {
		log.Infof("Mailgun : %s", err.Error())
	}
	return &Mailer{mg}
}

func (mg *Mailer) SendMessage(subject, text, to string) (string, error) {
	newMessage := mg.NewMessage("aniqma@aniqma.com", subject, text, to)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, id, err := mg.Send(ctx, newMessage)
	return id, err
}
