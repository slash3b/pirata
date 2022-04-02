package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"scraper/dto"
	"scraper/service"
	"scraper/storage/repository"
	"time"

	"github.com/mailjet/mailjet-apiv3-go/v3"

	_ "github.com/mattn/go-sqlite3"
)

/*
	converter.go converter -- gets a strings parsed and transforms them to a Film Model ?
	cineplex.go service -- knows how and gets raw films and uses converter to get those

	final product of service is a batch of Film models


	I want to add those to the database
	I want to send email with films, possibly updated with more info, trailers and so on.
*/

/*
	for now
	lets have just two services:
	- this with everything inside
	- prometheus
	- grafana

	todo: figure what to do with logging
*/

func main() {

	/**
	DB stuff
	file:test.db?cache=shared&mode=memory
	*/

	db, err := sql.Open("sqlite3", "file:../../pirata.db") // todo update: pass through ENV variables
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	filmRepo := repository.NewRepository(db)

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	scraperClient := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           jar,
		Timeout:       time.Second * 30,
	}

	scraperService, err := service.NewCineplexScraper(scraperClient, filmRepo)
	if err != nil {
		log.Fatalln(err)
	}

	imdbService := service.NewIMDB(scraperClient)

	allNewFilms, err := scraperService.GetAllFilms() // todo: better naming
	if err != nil {
		log.Fatalln(err)
	}

	emailFilms := []dto.EmailFilm{}
	for _, v := range allNewFilms {
		emailFilms = append(emailFilms, dto.FromModel(v, imdbService.GetFilmData(v)))
	}

	// todo : create a static page where users might subscribe ?
	// YES and unsubscribe also

	// use static fs
	tpl, err := template.ParseFiles("static/email.html")
	if err != nil {
		panic(err)
	}

	b := bytes.NewBufferString("")
	wr := bufio.NewWriter(b)

	err = tpl.Execute(wr, emailFilms)
	if err != nil {
		panic(err)
	}

	htmlOutput := b.String()

	mjPubK := os.Getenv("MJ_APIKEY_PUBLIC")
	mjPrivK := os.Getenv("MJ_APIKEY_PRIVATE")
	fromEmail := os.Getenv("FROM_EMAIL")
	fromName := os.Getenv("FROM_NAME")

	mailjetClient := mailjet.NewMailjetClient(mjPubK, mjPrivK)
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: fromEmail,
				Name:  fromName,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: "slash3b@gmail.com",
					Name:  "Ilya",
				},
			},
			Subject:  "Yo! Pirata has found some new upcoming movies in cineplex cinema",
			HTMLPart: htmlOutput,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err = mailjetClient.SendMailV31(&messages)
	if err != nil {
		log.Fatal(err)
	}
}
