package service

import (
	"bufio"
	"bytes"
	"embed"
	"html/template"
	"log"
	"scraper/dto"
	"scraper/metrics"
	"scraper/storage/repository"

	"github.com/mailjet/mailjet-apiv3-go/v3"

	_ "embed"
)

//go:embed template/email.html
var emailTpl embed.FS

const subject = "🍿 Yo! Check these upcoming movies in cineplex cinema"

type Sender interface {
	Send(films <-chan dto.EmailFilm) error
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

func NewMailer(mjClient *mailjet.Client, conf MailerConfig, r *repository.SubscriberRepository) *Mailer {
	return &Mailer{
		cl:   mjClient,
		conf: conf,
		repo: r,
	}
}

func (m *Mailer) Send(emailFilms <-chan dto.EmailFilm) error {
	tpl, err := template.ParseFS(emailTpl, "template/email.html")
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("could_not_parse_email_template").Inc()
		return err
	}

	b := bytes.NewBufferString("")
	wr := bufio.NewWriter(b)

	var films []dto.EmailFilm
	for ef := range emailFilms {
		films = append(films, ef)
	}

	if len(films) == 0 {
		return nil
	}

	err = tpl.Execute(wr, films)
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("could_not_execute_email_template").Inc()
		return err
	}

	htmlOutput := b.String()

	fromRecipient := &mailjet.RecipientV31{
		Email: m.conf.FromEmail,
		Name:  m.conf.FromName,
	}

	allSubs, err := m.getRecipients()
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("unable_to_get_recipients").Inc()
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
			// todo: implement retry with exponential backoff just for fun
			log.Println(err)
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
