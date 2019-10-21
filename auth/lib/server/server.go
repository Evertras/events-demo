package server

import (
	"net/http"
	"time"

	"github.com/Evertras/events-demo/auth/lib/authdb"
)

func New(addr string, db authdb.Db) *http.Server {
	router := http.NewServeMux()

	router.HandleFunc("/check", checkAuthHandler)
	router.HandleFunc("/login", loginHandler(db))
	router.HandleFunc("/register", registerHandler(db))

	return &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		Handler:      router,
	}
}
