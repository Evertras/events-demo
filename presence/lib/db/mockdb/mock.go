package mockdb

import (
	"context"
	"sync"

	"github.com/Evertras/events-demo/presence/lib/db"
)

type MockDb struct {
	m sync.RWMutex

	friendLists   map[string][]string
	sessions      map[string]db.SessionData
	notifications chan db.PresenceChangedEvent
}

func New() *MockDb {
	return &MockDb{
		friendLists:   make(map[string][]string),
		sessions:      make(map[string]db.SessionData),
		notifications: make(chan db.PresenceChangedEvent, 100),
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

func (d *MockDb) Subscribe(ctx context.Context) (chan db.PresenceChangedEvent, error) {
	return d.notifications, nil
}

func (d *MockDb) SendNotification(ctx context.Context, n db.PresenceChangedEvent) error {
	// Just feed it back to ourselves
	d.notifications <- n

	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Mock-specific methods

func (d *MockDb) MakeFriends(a string, b string) {
	d.m.Lock()
	defer d.m.Unlock()

	if aList, ok := d.friendLists[a]; ok {
		d.friendLists[a] = append(aList, b)
	} else {
		d.friendLists[a] = []string{b}
	}

	if bList, ok := d.friendLists[b]; ok {
		d.friendLists[b] = append(bList, a)
	} else {
		d.friendLists[b] = []string{a}
	}
}

func (d *MockDb) SetFriendList(id string, friends []string) {
	for _, friendID := range friends {
		d.MakeFriends(id, friendID)
	}
}

func (d *MockDb) Connect(id string) {
	d.m.Lock()
	defer d.m.Unlock()

	d.notifications <- db.PresenceChangedEvent{
		PlayerID:  id,
		NotifyIDs: d.friendLists[id],
		Online:    true,
	}

	d.sessions[id] = db.SessionData{}
}

func (d *MockDb) Disconnect(id string) {
	d.m.Lock()
	defer d.m.Unlock()

	d.notifications <- db.PresenceChangedEvent{
		PlayerID:  id,
		NotifyIDs: d.friendLists[id],
		Online:    false,
	}

	delete(d.sessions, id)
}

func (d *MockDb) Close() {
	close(d.notifications)
}
