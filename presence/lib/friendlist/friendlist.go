package friendlist

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Evertras/events-demo/presence/lib/db"
)

type FriendStatus struct {
	Id     string
	Online bool
}

type FriendList interface {
	GetFriendStatus(ctx context.Context, id string) ([]FriendStatus, error)
	Notifications(ctx context.Context, id string) (chan FriendStatus, error)
}

type friendList struct {
	db db.Db
}

func New(db db.Db) FriendList {
	return &friendList {
		db: db,
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
			Id: id,
			Online: false,
		}

		if _, ok := activePlayers[id]; ok {
			s.Online = true
		}

		statusList = append(statusList, s)
	}

	return statusList, nil
}

func (f *friendList) Notifications(ctx context.Context, id string) (chan FriendStatus, error) {
	return nil, nil
}
