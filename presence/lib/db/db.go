package db

import "context"

type SessionData struct {
}

type PresenceChangedEvent struct {
	PlayerID  string
	NotifyIDs []string
	Online    bool
}

type Db interface {
	GetFriendList(ctx context.Context, id string) ([]string, error)

	GetSessionData(ctx context.Context, ids []string) (map[string]SessionData, error)

	Subscribe(ctx context.Context) (chan PresenceChangedEvent, error)
}
