package main

import (
	"context"
	"log"

	"github.com/Evertras/events-demo/presence/lib/db/mockdb"
	"github.com/Evertras/events-demo/presence/lib/friendlist"
	"github.com/Evertras/events-demo/presence/lib/server"
	"github.com/Evertras/events-demo/presence/lib/server/wslistener"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wsl := wslistener.New("0.0.0.0:8111", "/")
	d := mockdb.New()
	fl := friendlist.New(d)

	go fl.ListenForChanges(ctx)

	d.MakeFriends("A", "B")
	d.MakeFriends("A", "C")

	s := server.New(wsl, d, fl)

	log.Fatal(s.Run(ctx))
}
