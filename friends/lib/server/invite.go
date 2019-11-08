package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Evertras/events-demo/friends/lib/events"
	"github.com/Evertras/events-demo/friends/lib/events/friendevents"
	"github.com/Evertras/events-demo/shared/stream"
)

type InviteBody struct {
	ToID string `json:"id"`
	ToEmail string `json:"email"`
}

func inviteHandler(streamWriter stream.Writer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span, ctx := startSpan("invite", r)
		defer span.Finish()

		if r.Method != "POST" {
			w.WriteHeader(400)
			log.Println("Method must be POST")
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to read body:", err)
			return
		}

		var inviteBody InviteBody

		err = json.Unmarshal(body, &inviteBody)

		from := r.Header.Get("X-User-ID")

		if from == "" {
			w.WriteHeader(400)
			log.Println("Failed to get user ID from header")
			return
		}

		ev := friendevents.NewInviteSent()

		ev.FromID = from
		ev.ToID = inviteBody.ToID
		ev.ToEmail = inviteBody.ToEmail
		ev.TimeUnixMs = time.Now().Unix()

		err = streamWriter.Write(ctx, []byte(from), events.EventIDInviteSent, ev)

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to write to stream:", err)
			return
		}
	}
}
