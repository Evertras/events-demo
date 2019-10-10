package main

import (
	"log"
	"net/http"
	"time"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("helloHandler start")

	w.Write([]byte("Hello\n"))
	w.Write([]byte(r.URL.String() + "\n"))
	w.Write([]byte(r.URL.Hostname() + "\n"))
	w.Write([]byte(r.URL.RawPath + "\n"))

	log.Println("helloHandler end")
}

func headerEchoHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("headerEchoHandler start")

	for header := range r.Header {
		w.Write([]byte(header + ": " + r.Header.Get(header) + "\n"))
	}

	log.Println("headerEchoHandler end")
}

func main() {
	addr := "0.0.0.0:13030"

	router := http.NewServeMux()

	router.HandleFunc("/", helloHandler)
	router.HandleFunc("/headers", headerEchoHandler)

	server := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		Handler:      router,
	}

	log.Println("Serving", addr)

	log.Fatal(server.ListenAndServe())
}
