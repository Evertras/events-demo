package server

import "net/http"

func inviteHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span, _ := startSpan("invite", r)
		defer span.Finish()
	}
}
