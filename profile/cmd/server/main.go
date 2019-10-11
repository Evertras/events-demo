package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Profile struct {
	Intro string `json:"intro"`
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("profileHandler start")
	defer log.Println("profileHandler end")

	p := Profile{
		Intro: "This is a profile returned from the profile service!",
	}

	data, err := json.Marshal(p)

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {
	addr := "0.0.0.0:13083"

	router := http.NewServeMux()

	router.HandleFunc("/", profileHandler)

	server := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		Handler:      router,
	}

	log.Println("Serving", addr)

	log.Fatal(server.ListenAndServe())
}
