package mailgun

import (
	"context"
	"fmt"
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
func (mg *Mailer) SendMessage(variables map[string]interface{}, template, subject, to string) (string, error) {

	newMessage := mg.NewMessage("lokaventour@gmail.com", subject, "", to)
	newMessage.SetTemplate(template)

	for key, val := range variables {
		fmt.Println("Key : ", key, "Value : ", val)
		newMessage.AddVariable(key, val)
	}

	fmt.Println("Email : ", to)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, id, err := mg.Send(ctx, newMessage)
	return id, err

}

//AddListMemberHost : method for adding member host
func (mg *Mailer) AddListMemberHost(member *Member, variable map[string]interface{}) error {

	memberHost := _mailgun.Member{
		Address:    member.Address,
		Name:       member.Name,
		Subscribed: _mailgun.Subscribed,
		Vars:       variable,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateMember(ctx, true, "host@mg.lokaventour.com", memberHost)
}

//AddListMemberGuest : method for adding member guest
func (mg *Mailer) AddListMemberGuest(member *Member, variable map[string]interface{}) error {

	memberGuest := _mailgun.Member{
		Address:    member.Address,
		Name:       member.Name,
		Subscribed: _mailgun.Subscribed,
		Vars:       variable,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateMember(ctx, true, "guest@mg.lokaventour.com", memberGuest)
}

//AddMemberSubscription : method for adding subscription list
func (mg *Mailer) AddMemberSubscription(member *Member, variable map[string]interface{}) error {

	memberGuest := _mailgun.Member{
		Address:    member.Address,
		Name:       member.Name,
		Subscribed: _mailgun.Subscribed,
		Vars:       variable,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.CreateMember(ctx, true, "subscription@mg.lokaventour.com", memberGuest)
}

//SendToSubs : to send subs an email
func (mg *Mailer) SendToSubs(to string) (string, error) {
	newMessage := mg.NewMessage("lokaventour.com@gmail.com",
		"Hello, %recipient.name%",
		"", to)
	newMessage.SetTemplate("lokaven-pilot-user")
	newMessage.AddTemplateVariable("title", "Hello, %recipient.name%")
	newMessage.AddVariable("fullname", "%recipient.name%")
	newMessage.AddVariable("top_message", "Terimakasih sudah bergabung dengan lokaven. Berikut ini adalah kredensial akun anda:")
	newMessage.AddVariable("email", "%recipient_email%")
	newMessage.AddVariable("password", "%recipient.password%")
	newMessage.AddVariable("bottom_message", "Silahkan unduh aplikasi kami di tautan berikut ini:")
	newMessage.AddVariable("link", "%recipient.link%")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, id, err := mg.Send(ctx, newMessage)
	return id, err
}
