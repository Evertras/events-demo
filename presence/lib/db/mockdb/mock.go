package mockdb

import (
	"context"
	"sync"

	"github.com/Evertras/events-demo/presence/lib/db"
)

type MockDb struct {
	m sync.RWMutex

	friendLists map[string][]string
	sessions map[string]db.SessionData
}

func New() *MockDb {
	return &MockDb{
		friendLists: make(map[string][]string),
		sessions: make(map[string]db.SessionData),
	}
}

// Quick type enforcement check
var _ db.Db = &MockDb{}

///////////////////////////////////////////////////////////////////////////////
// Interface matching stuff

func (d *MockDb) GetFriendList(ctx context.Context, id string) ([]string, error) {
	d.m.RLock()
	defer d.m.RUnlock()

	if friends, ok := d.friendLists[id]; ok {
		return friends, nil
	}

	return nil, nil
}

func (d *MockDb) GetSessionData(ctx context.Context, ids []string) (map[string]db.SessionData, error) {
	d.m.RLock()
	defer d.m.RUnlock()

	data := make(map[string]db.SessionData)

	for _, id := range ids {
		if session, ok := d.sessions[id]; ok {
			data[id] = session
		}
	}

	return data, nil
}

///////////////////////////////////////////////////////////////////////////////
// Mock-specific methods

func (d *MockDb) SetFriendList(id string, friends []string) {
	d.m.Lock()
	defer d.m.Unlock()

	d.friendLists[id] = friends
}

func (d *MockDb) Connect(id string) {
	d.m.Lock()
	defer d.m.Unlock()

	d.sessions[id] = db.SessionData{}
}

func (d *MockDb) Disconnect(id string) {
	d.m.Lock()
	defer d.m.Unlock()

	delete(d.sessions, id)
}
