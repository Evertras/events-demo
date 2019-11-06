package server

import (
	"log"
	"net/http"
)

func headerEchoHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span, _ := startSpan("header echo", r)
		defer span.Finish()

		for header := range r.Header {
			w.Write([]byte(header + ": " + r.Header.Get(header) + "\n"))
		}

		log.Println("headerEchoHandler end")
	}
}
