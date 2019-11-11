package mock

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type Db struct {
	m sync.Mutex

	invitesByID map[string][]string
	idFromEmail map[string]string
	emailFromID map[string]string
}

func New() *Db {
	return &Db{
		invitesByID: make(map[string][]string),
		idFromEmail: make(map[string]string),
		emailFromID: make(map[string]string),
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

	d.idFromEmail[email] = userID

	return nil
}

func (d *Db) SendInviteByID(ctx context.Context, t time.Time, fromID string, toID string) error {
	d.m.Lock()
	defer d.m.Unlock()

	if _, exists := d.invitesByID[toID]; !exists {
		d.invitesByID[toID] = make([]string, 0)
	}

	d.invitesByID[toID] = append(d.invitesByID[toID], fromID)

	return nil
}

func (d *Db) GetIDFromEmail(ctx context.Context, email string) (string, error) {
	d.m.Lock()
	defer d.m.Unlock()

	id, exists := d.idFromEmail[email]

	if !exists {
		return "", errors.New("does not exist")
	}

	return id, nil
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

	id := toID

	if len(toID) == 0 {
		id = d.idFromEmail[toEmail]
	}

	invitesByID, exists := d.invitesByID[id]

	if exists {
		for _, i := range invitesByID {
			if i == fromID {
				return true
			}
		}
	}

	return false
}
