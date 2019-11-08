package mock

import (
	"context"
	"sync"
	"time"
)

type Db struct {
	m sync.Mutex

	invitesByID map[string][]string
	invitesByEmail map[string][]string
}

func New() *Db {
	return &Db{
		invitesByID: make(map[string][]string),
		invitesByEmail: make(map[string][]string),
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

	if _, exists := d.invitesByID[userID]; !exists {
		d.invitesByID[userID] = make([]string, 0)
	}

	return nil
}

func (d *Db) SendInviteByID(ctx context.Context, t time.Time, fromID string, toID string) error {
	// Cheat a little here for the sake of simplicity, just ensure they exist
	d.CreatePlayer(ctx, fromID, "f")
	d.CreatePlayer(ctx, toID, "f")

	d.m.Lock()
	defer d.m.Unlock()

	d.invitesByID[toID] = append(d.invitesByID[toID], fromID)

	return nil
}

func (d *Db) SendInviteByEmail(ctx context.Context, t time.Time, fromID string, toEmail string) error {
	// Cheat a little here for the sake of simplicity, just ensure they exist
	d.CreatePlayer(ctx, fromID, "f")
	d.CreatePlayer(ctx, "asdf", toEmail)

	d.m.Lock()
	defer d.m.Unlock()

	d.invitesByEmail[toEmail] = append(d.invitesByEmail[toEmail], fromID)

	return nil
}

func (d *Db) GetPendingInvites(ctx context.Context, id string) ([]string, error) {
	d.m.Lock()
	defer d.m.Unlock()

	if pending, ok := d.invitesByID[id]; ok {
		return pending, nil
	}

	return []string{}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Mock specific stuff

func (d *Db) MockPlayerExists(id string) bool {
	d.m.Lock()
	defer d.m.Unlock()

	_, exists := d.invitesByID[id]

	return exists
}

func (d *Db) MockInviteExists(fromID string, toID string, toEmail string) bool {
	d.m.Lock()
	defer d.m.Unlock()

	invitesByID, exists := d.invitesByID[toID]

	if exists {
		for _, i := range invitesByID {
			if i == fromID {
				return true
			}
		}
	}

	invitesByEmail, exists := d.invitesByEmail[toEmail]

	if exists {
		for _, i := range invitesByEmail {
			if i == fromID {
				return true
			}
		}
	}

	return false
}
