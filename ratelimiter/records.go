package ratelimiter

type Entry struct {
	startTime, lastInvocation int64
	invocations               int
	subscription              Subscription
}

func NewEntry(subType Subscription) *Entry {
	return &Entry{subscription: subType}
}

type ClientMap struct {
	Entries map[string]*Entry
}

func (cm *ClientMap) ValidClientId(clientId string) bool {
	_, ok := cm.Entries[clientId]
	return ok
}
