package db

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type SessionData struct {
}

type PresenceChangedEvent struct {
	PlayerID  string
	NotifyIDs []string
	Online    bool
}

type ConnectionOptions struct {
	Address string
}

type Db interface {
	GetFriendList(ctx context.Context, id string) ([]string, error)
	SetFriendList(ctx context.Context, id string, friends []string) error

	GetSessionData(ctx context.Context, ids []string) (map[string]SessionData, error)

	Subscribe(ctx context.Context) (chan PresenceChangedEvent, error)

	SendNotification(ctx context.Context, n PresenceChangedEvent) error
}

type db struct {
	client *redis.Client
}

func New(opts ConnectionOptions) Db {
	return &db{
		client: redis.NewClient(&redis.Options{
			Addr:     opts.Address,
			Password: "",
			DB:       0,
		}),
	}
}

func (d *db) keyFriendList(id string) string {
	return "friendlist:" + id
}

func (d *db) keySession(id string) string {
	return "session:" + id
}

func (d *db) GetFriendList(ctx context.Context, id string) ([]string, error) {
	raw, err := d.client.Get(d.keyFriendList(id)).Result()

	if err != nil {
		return nil, errors.Wrap(err, "failed to get from redis")
	}

	var fl DataFriendList

	err = json.Unmarshal([]byte(raw), &fl)

	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal result")
	}

	return fl.Friends, nil
}

func (d *db) SetFriendList(ctx context.Context, id string, friends []string) error {
	fl := DataFriendList {
		Friends: friends,
	}

	raw, err := json.Marshal(fl)

	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return d.client.Set(d.keyFriendList(id), string(raw), 0).Err()
}

func (d *db) GetSessionData(ctx context.Context, ids []string) (map[string]SessionData, error) {
	return nil, nil
}

func (d *db) Subscribe(ctx context.Context) (chan PresenceChangedEvent, error) {
	return nil, nil
}

func (d *db) SendNotification(ctx context.Context, n PresenceChangedEvent) error {
	return nil
}
