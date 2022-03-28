package middleware

import "net/http"

func RateLimiter(maxRunning int, f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	limitingBuffer := make(chan struct{}, maxRunning)
	return func(w http.ResponseWriter, r *http.Request) {

		limitingBuffer <- struct{}{}
		defer func() { <-limitingBuffer }()

		f(w, r)
	}
}
