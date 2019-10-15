package server

import (
	"net/http"
	"time"
)

func New(addr string) *http.Server {
	router := http.NewServeMux()

	router.HandleFunc("/check", checkAuthHandler)

	return &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		Handler:      router,
	}
}
