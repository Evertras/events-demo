package friendlist

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/Evertras/events-demo/presence/lib/db"
)

type FriendStatus struct {
	ID     string
	Online bool
}

type FriendList interface {
	ListenForChanges(ctx context.Context) error
	GetFriendStatus(ctx context.Context, id string) ([]FriendStatus, error)
	Subscribe(ctx context.Context, id string) (chan FriendStatus, error)
	Unsubscribe(ctx context.Context, id string)
}

type friendList struct {
	db db.Db

	subLock       sync.RWMutex
	subscriptions map[string]chan FriendStatus
}

func New(db db.Db) FriendList {
	return &friendList{
		db:            db,
		subscriptions: make(map[string]chan FriendStatus),
	}
}

func (f *friendList) ListenForChanges(ctx context.Context) error {
	events, err := f.db.Subscribe(ctx)

	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case ev, open := <-events:
			if !open {
				return nil
			}

			f.subLock.RLock()

			for _, id := range ev.NotifyIDs {
				// If no one is listening, just drop it on the floor
				if s, ok := f.subscriptions[id]; ok {
					s <- FriendStatus{
						ID:     ev.PlayerID,
						Online: ev.Online,
					}
				}
			}

			f.subLock.RUnlock()
		}
	}
}

func (f *friendList) GetFriendStatus(ctx context.Context, id string) ([]FriendStatus, error) {
	friends, err := f.db.GetFriendList(ctx, id)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get friend list")
	}

	activePlayers, err := f.db.GetSessionData(ctx, friends)

	statusList := make([]FriendStatus, 0, len(friends))

	for _, id := range friends {
		s := FriendStatus{
			ID:     id,
			Online: false,
		}

		if _, ok := activePlayers[id]; ok {
			s.Online = true
		}

		statusList = append(statusList, s)
	}

	return statusList, nil
}

func (f *friendList) Subscribe(ctx context.Context, id string) (chan FriendStatus, error) {
	f.subLock.Lock()
	defer f.subLock.Unlock()

	if _, ok := f.subscriptions[id]; ok {
		return nil, errors.New("already subscribed for this ID")
	}

	c := make(chan FriendStatus, 10)

	f.subscriptions[id] = c

	return c, nil
}

func (f *friendList) Unsubscribe(ctx context.Context, id string) {
	f.subLock.Lock()
	defer f.subLock.Unlock()

	delete(f.subscriptions, id)
}
