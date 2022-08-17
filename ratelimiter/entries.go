package ratelimiter

type Entry interface {
	StartTime() int64
	LastInvocation() int64
	Invocations() int
	Subscription() Subscription
}

type WindowEntry struct {
	startTime, lastInvocation int64
	invocations               int
	subscription              Subscription
}

type TokenEntry struct {
	startTime, lastInvocation int64
	invocations               int
	subscription              Subscription
	queue                     *BasicQueue[int]
}

func NewWindowEntry(subType Subscription) *WindowEntry {
	return &WindowEntry{subscription: subType}
}

func (entry *WindowEntry) StartTime() int64 {
	return entry.startTime
}

func (entry *WindowEntry) LastInvocation() int64 {
	return entry.lastInvocation
}

func (entry *WindowEntry) Invocations() int {
	return entry.invocations
}

func (entry *WindowEntry) Subscription() Subscription {
	return entry.subscription
}

func NewTokenEntry(subType Subscription) *TokenEntry {
	return &TokenEntry{subscription: subType, queue: NewBasicQueue[int](subType.RequestLimit())}
}

func (entry *TokenEntry) StartTime() int64 {
	return entry.startTime
}

func (entry *TokenEntry) LastInvocation() int64 {
	return entry.lastInvocation
}

func (entry *TokenEntry) Invocations() int {
	return entry.invocations
}

func (entry *TokenEntry) Subscription() Subscription {
	return entry.subscription
}

type ClientMap struct {
	Entries map[string]Entry
}

func (cm *ClientMap) ValidClientId(clientId string) bool {
	_, ok := cm.Entries[clientId]
	return ok
}
