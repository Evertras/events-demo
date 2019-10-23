package server

import (
	"net/http"
	"time"

	"github.com/Evertras/events-demo/auth/lib/auth"
)

func New(addr string, auth auth.Auth) *http.Server {
	router := http.NewServeMux()

	router.HandleFunc("/check", checkAuthHandler)
	router.HandleFunc("/login", loginHandler(auth))
	router.HandleFunc("/register", registerHandler(auth))

	return &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		Handler:      router,
	}
}
