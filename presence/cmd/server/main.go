package main

import (
	"context"
	"log"

	"github.com/Evertras/events-demo/presence/lib/db"
	"github.com/Evertras/events-demo/presence/lib/friendlist"
	"github.com/Evertras/events-demo/presence/lib/server"
	"github.com/Evertras/events-demo/presence/lib/server/wslistener"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wsl := wslistener.New("0.0.0.0:13337", "/")
	d := db.New(db.ConnectionOptions{
		Address: "presence-db:6379",
	})

	fl := friendlist.New(d)

	go fl.ListenForChanges(ctx)

	if err := d.SetFriendList(ctx, "A", []string { "B", "C" }); err != nil {
		log.Fatal(err)
	}

	if err := d.SetFriendList(ctx, "B", []string { "A" }); err != nil {
		log.Fatal(err)
	}

	if err := d.SetFriendList(ctx, "C", []string { "A" }); err != nil {
		log.Fatal(err)
	}

	s := server.New(wsl, d, fl)

	log.Fatal(s.Run(ctx))
}
