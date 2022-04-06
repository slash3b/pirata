package service

import (
	"bufio"
	"bytes"
	"embed"
	"html/template"
	"scraper/dto"
	"scraper/storage/repository"

	"github.com/mailjet/mailjet-apiv3-go/v3"

	_ "embed"
)

//go:embed template/email.html
var emailTpl embed.FS

const subject = "🍿 Yo! Check these upcoming movies in cineplex cinema"

type Sender interface {
	Send(films []dto.EmailFilm) error
}

type Mailer struct {
	cl   *mailjet.Client
	conf MailerConfig
	repo *repository.SubscriberRepository
}

type MailerConfig struct {
	FromEmail string
	FromName  string
}

func NewMailer(mjClient *mailjet.Client, conf MailerConfig, r *repository.SubscriberRepository,
) *Mailer {
	return &Mailer{
		cl:   mjClient,
		conf: conf,
		repo: r,
	}
}

func (m *Mailer) Send(films []dto.EmailFilm) error {

	tpl, err := template.ParseFS(emailTpl, "template/email.html")
	if err != nil {
		panic(err)
	}

	b := bytes.NewBufferString("")
	wr := bufio.NewWriter(b)

	err = tpl.Execute(wr, films)
	if err != nil {
		panic(err)
	}

	htmlOutput := b.String()

	fromRecipient := &mailjet.RecipientV31{
		Email: m.conf.FromEmail,
		Name:  m.conf.FromName,
	}

	allSubs, err := m.getRecipients()
	if err != nil {
		return err
	}

	for _, subscriber := range allSubs {
		messagesInfo := []mailjet.InfoMessagesV31{
			{
				From:     fromRecipient,
				To:       &mailjet.RecipientsV31{subscriber},
				Subject:  subject,
				HTMLPart: htmlOutput,
			},
		}

		messages := mailjet.MessagesV31{Info: messagesInfo}

		_, err = m.cl.SendMailV31(&messages)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Mailer) getRecipients() ([]mailjet.RecipientV31, error) {
	var subs []mailjet.RecipientV31

	allSubs, err := m.repo.GetAllActive()
	if err != nil {
		return subs, err
	}

	for _, sub := range allSubs {
		subs = append(subs, mailjet.RecipientV31{
			Email: sub.Email,
			Name:  sub.Name,
		})
	}

	return subs, nil
}
