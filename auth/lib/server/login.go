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
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func loginHandler(db authdb.Db) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(400)
			log.Println("Method must be POST")
			return
		}

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

		valid, err := db.ValidateUser(login.Email, login.Password)

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to validate user:", err)
			return
		}

		if !valid {
			w.WriteHeader(400)
			log.Println("Bad credentials for", login.Email)
			return
		}

		t, err := token.New(login.Email)

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to generate token:", err)
			return
		}

		resBody, err := json.Marshal(TokenResponse{Token: t})

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to marshal token response:", err)
			return
		}

		_, err = w.Write(resBody)

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to write body:", err)
			return
		}

		log.Println("Login successful for " + login.Email)
	}
}
