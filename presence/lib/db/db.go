package db

import "context"

type SessionData struct {
}

type Db interface {
	GetFriendList(ctx context.Context, id string) ([]string, error)

	GetSessionData(ctx context.Context, ids []string) (map[string]SessionData, error)
}
