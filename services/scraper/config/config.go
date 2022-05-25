package config

import (
	"fmt"
	"os"
)

type Env struct {
	MailjetPubKey     string
	MailJetPrivateKey string
	FromEmail         string
	FromName          string
}

func GetEnv() (*Env, error) {
	mjPubKey, exists := os.LookupEnv("MJ_APIKEY_PUBLIC")
	if !exists {
		return nil, fmt.Errorf("could not lookup MJ_APIKEY_PUBLIC variable")
	}

	mjPrivKey, exists := os.LookupEnv("MJ_APIKEY_PRIVATE")
	if !exists {
		return nil, fmt.Errorf("could not lookup MJ_APIKEY_PRIVATE variable")
	}

	fromEmail, exists := os.LookupEnv("FROM_EMAIL")
	if !exists {
		return nil, fmt.Errorf("could not lookup FROM_EMAIL variable")
	}

	fromName, exists := os.LookupEnv("FROM_NAME")
	if !exists {
		return nil, fmt.Errorf("could not lookup FROM_NAME variable")
	}

	return &Env{
		MailjetPubKey:     mjPubKey,
		MailJetPrivateKey: mjPrivKey,
		FromName:          fromName,
		FromEmail:         fromEmail,
	}, nil
}
