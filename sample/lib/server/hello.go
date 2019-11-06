package server

import (
	"log"
	"net/http"
)

func helloHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span, _ := startSpan("hello", r)
		defer span.Finish()

		w.Write([]byte("Hello\n"))
		w.Write([]byte(r.URL.String() + "\n"))
		w.Write([]byte(r.URL.Hostname() + "\n"))
		w.Write([]byte(r.URL.RawPath + "\n"))

		log.Println("helloHandler end")
	}
}
