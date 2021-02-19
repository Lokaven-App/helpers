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

//Member : struct that define member mailing
type Member struct {
	Address string
	Name    string
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

//AddListMemberHost : method for adding member host
func (mg *Mailer) AddListMemberHost(member *Member) error {

	memberHost := _mailgun.Member{
		Address:    member.Address,
		Name:       member.Name,
		Subscribed: _mailgun.Subscribed,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateMember(ctx, true, "host@mg.lokaventour.com", memberHost)
}

//AddListMemberGuest : method for adding member guest
func (mg *Mailer) AddListMemberGuest(member *Member) error {

	memberGuest := _mailgun.Member{
		Address:    member.Address,
		Name:       member.Name,
		Subscribed: _mailgun.Subscribed,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateMember(ctx, true, "guest@mg.lokaventour.com", memberGuest)
}
