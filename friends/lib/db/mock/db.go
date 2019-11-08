package mock

import (
	"context"
	"sync"
	"time"
)

type Db struct {
	m sync.Mutex

	Invites map[string][]string
}

func New() *Db {
	return &Db{
		Invites: make(map[string][]string),
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

func (d *Db) CreatePlayer(ctx context.Context, userID string) error {
	d.m.Lock()
	defer d.m.Unlock()

	if _, exists := d.Invites[userID]; !exists {
		d.Invites[userID] = make([]string, 0)
	}

	return nil
}

func (d *Db) SendInvite(ctx context.Context, t time.Time, fromID string, toID string) error {
	// Cheat a little here for the sake of simplicity, just ensure they exist
	d.CreatePlayer(ctx, fromID)
	d.CreatePlayer(ctx, toID)

	d.m.Lock()
	defer d.m.Unlock()

	d.Invites[toID] = append(d.Invites[toID], fromID)

	return nil
}

func (d *Db) GetPendingInvites(ctx context.Context, id string) ([]string, error) {
	d.m.Lock()
	defer d.m.Unlock()

	if pending, ok := d.Invites[id]; ok {
		return pending, nil
	}

	return []string{}, nil
}
