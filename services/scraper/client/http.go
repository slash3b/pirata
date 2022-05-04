package client

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

func NewHttpClientWithCookies() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           jar,
		Timeout:       time.Second * 30,
	}, nil
}
