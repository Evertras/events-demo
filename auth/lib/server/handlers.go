package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Evertras/events-demo/auth/lib/authdb"
	"github.com/Evertras/events-demo/auth/lib/token"
)

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func checkAuthHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("X-Auth-Token")

	if authToken == "" {
		w.WriteHeader(401)
		log.Println("No header")
		return
	}

	claim, err := token.Parse(authToken)

	if err != nil {
		w.WriteHeader(401)
		log.Println("Could not validate header:", err)
		return
	}

	w.Header().Set("X-User-ID", claim.Username)

	w.WriteHeader(200)
}

func loginHandler(db authdb.Db) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to parse:", err)
			return
		}

		var login LoginBody
		err = json.Unmarshal(body, &login)

		if err != nil {
			w.WriteHeader(400)
			log.Println("Could not parse body:", err)
			return
		}

		valid, err := db.ValidateUser(login.Username, login.Password)

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to validate user:", err)
			return
		}

		if !valid {
			w.WriteHeader(400)
			log.Println("Bad credentials for", login.Username)
			return
		}

		header, err := token.New(login.Username)

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to generate token:", err)
			return
		}

		_, err = w.Write([]byte(header))

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to write body:", err)
			return
		}
	}
}
