package db

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type PresenceChangedEvent struct {
	PlayerID  string
	NotifyIDs []string
	Online    bool
}

type ConnectionOptions struct {
	Address string
	HostKey string
}

type Db interface {
	GetFriendList(ctx context.Context, id string) ([]string, error)
	SetFriendList(ctx context.Context, id string, friends []string) error

	GetSessionData(ctx context.Context, ids []string) (map[string]DataSession, error)
	SetSessionData(ctx context.Context, id string, data DataSession) error

	Subscribe(ctx context.Context) (chan PresenceChangedEvent, error)

	SendNotification(ctx context.Context, n PresenceChangedEvent) error
}

type db struct {
	client *redis.Client
	opts   ConnectionOptions
}

func New(opts ConnectionOptions) Db {
	return &db{
		client: redis.NewClient(&redis.Options{
			Addr:     opts.Address,
			Password: "",
			DB:       0,
		}),
		opts: opts,
	}
}

func (d *db) keyFriendList(id string) string {
	return "friendlist:" + id
}

func (d *db) keySession(id string) string {
	return "session:" + id
}

func (d *db) keyPubsubChannel(hostkey string) string {
	return "s-notify:" + hostkey
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
	fl := DataFriendList{
		Friends: friends,
	}

	raw, err := json.Marshal(fl)

	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return d.client.Set(d.keyFriendList(id), string(raw), 0).Err()
}

func (d *db) GetSessionData(ctx context.Context, ids []string) (map[string]DataSession, error) {
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = d.keySession(id)
	}
	sessions, err := d.client.MGet(keys...).Result()

	if err != nil {
		return nil, errors.Wrap(err, "failed to mget")
	}

	results := make(map[string]DataSession)

	for i, raw := range sessions {
		if raw == nil {
			continue
		}

		var s DataSession
		err = json.Unmarshal([]byte(raw.(string)), &s)

		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal JSON")
		}

		results[ids[i]] = s
	}

	return results, nil
}

func (d *db) SetSessionData(ctx context.Context, id string, data DataSession) error {
	data.Host = d.opts.HostKey

	raw, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	err = d.client.Set(d.keySession(id), string(raw), time.Minute*10).Err()

	if err != nil {
		return errors.Wrap(err, "failed to write to redis")
	}

	return nil
}

func (d *db) Subscribe(ctx context.Context) (chan PresenceChangedEvent, error) {
	changed := make(chan PresenceChangedEvent, 10)
	pubsub := d.client.Subscribe(d.keyPubsubChannel(d.opts.HostKey))

	go func() {
		defer pubsub.Close()

		msgs := pubsub.Channel()

		for {
			select {
			case msg := <-msgs:
				var n PresenceChangedEvent

				err := json.Unmarshal([]byte(msg.Payload), &n)

				if err != nil {
					log.Println("Error unmarshaling:", err)
					return
				}

				log.Printf("GOT MESSAGE: %+v", n)

				changed <- n

			case <-ctx.Done():
				return
			}
		}
	}()

	return changed, nil
}

func (d *db) SendNotification(ctx context.Context, n PresenceChangedEvent) error {
	raw, err := json.Marshal(n)

	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	// Find hosts we need to send notification too
	sessions, err := d.GetSessionData(ctx, n.NotifyIDs)

	byHost := make(map[string]PresenceChangedEvent)

	for _, session := range sessions {
		byHost[session.Host] = n
	}

	for host, _ := range byHost {
		log.Println("Publishing for", host)
		d.client.Publish(d.keyPubsubChannel(host), raw)
	}

	if err != nil {
		return errors.Wrap(err, "failed to get sessions")
	}

	return nil
}
