package mock

import (
	"context"
	"sync"
	"time"
)

type Db struct {
	m sync.Mutex

	invites map[string][]string
}

func New() *Db {
	return &Db{
		invites: make(map[string][]string),
	}
}

func (d *Db) Connect(ctx context.Context) error {
	return nil
}

func (d *Db) Close() error {
	return nil
}

func (d *Db) GetSharedValue(ctx context.Context, key string, ifNotSet string) (string, error) {
	return ifNotSet, nil
}

func (d *Db) CreatePlayer(ctx context.Context, userID string, email string) error {
	d.m.Lock()
	defer d.m.Unlock()

	if _, exists := d.invites[userID]; !exists {
		d.invites[userID] = make([]string, 0)
	}

	return nil
}

func (d *Db) SendInvite(ctx context.Context, t time.Time, fromID string, toID string) error {
	// Cheat a little here for the sake of simplicity, just ensure they exist
	d.CreatePlayer(ctx, fromID, "f")
	d.CreatePlayer(ctx, toID, "f")

	d.m.Lock()
	defer d.m.Unlock()

	d.invites[toID] = append(d.invites[toID], fromID)

	return nil
}

func (d *Db) GetPendingInvites(ctx context.Context, id string) ([]string, error) {
	d.m.Lock()
	defer d.m.Unlock()

	if pending, ok := d.invites[id]; ok {
		return pending, nil
	}

	return []string{}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Mock specific stuff

func (d *Db) MockPlayerExists(id string) bool {
	d.m.Lock()
	defer d.m.Unlock()

	_, exists := d.invites[id]

	return exists
}

func (d *Db) MockInviteExists(fromID string, toID string) bool {
	d.m.Lock()
	defer d.m.Unlock()

	invites, exists := d.invites[toID]

	if !exists {
		return false
	}

	for _, i := range invites {
		if i == fromID {
			return true
		}
	}

	return false
}
