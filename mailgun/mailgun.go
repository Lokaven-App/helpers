package mailgun

import (
	"context"
	"time"

	_mailgun "github.com/mailgun/mailgun-go/v3"
	log "github.com/sirupsen/logrus"
)

//Config : struct define config
type Config struct {
	Domain string
	Key    string
}

//Mailer : sruct that define Mailgun Implementation
type Mailer struct {
	*_mailgun.MailgunImpl
}

//Recepants : struct that define a group recepants
type Recepants struct {
	Name  string
	ID    string
	Email string
}

//Mailgun : set connection to mailgun
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

//SendMessage : to send email
func (mg *Mailer) SendMessage(subject, text, to string) (string, error) {
	newMessage := mg.NewMessage("lokaventour.com@gmail.com", subject, text, to)
	newMessage.SetTemplate("lokaven")
	newMessage.AddTemplateVariable("title", subject)
	newMessage.AddVariable("message", text)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, id, err := mg.Send(ctx, newMessage)
	return id, err
}

//SendGroup : to send email blast
func (mg *Mailer) SendGroup(subject, text string, newRecepants []*Recepants) (string, error) {
	m := make(map[string]interface{})
	for _, recepent := range newRecepants {
		newMessage := mg.NewMessage("lokaventour.com@gmail.com", subject, text, recepent.Email)
		newMessage.SetTemplate("lokaven")
		newMessage.AddTemplateVariable("title", subject)
		newMessage.AddVariable("message", text)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		m["id"] = recepent.ID
		m["name"] = recepent.Name

		newMessage.AddRecipientAndVariables(recepent.Email, m)

		_, id, err := mg.Send(ctx, newMessage)
		return id, err
	}
	return "", nil
}
