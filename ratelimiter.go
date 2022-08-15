package main

import (
	"fmt"
	"net/http"
	"time"
)

type Entry struct {
	startTime, lastInvocation int64
	invocations               int
	subscription              Subscription
}

type Subscription interface {
	Name() string
	RequestLimit() int
	TimeFrame() int64 // Millisecond count
}

type BasicSubscription struct {
	name         string
	requestLimit int
	timeFrame    int64
}

type PremiumSubscription struct {
	name         string
	requestLimit int
	timeFrame    int64
}

const (
	basicName         = "Basic"
	basicRequestLimit = 5
	basicTimeFrame    = 60000

	premiumName         = "Premium"
	premiumRequestLimit = 20
	premiumTimeFrame    = 60000
)

func (sub *BasicSubscription) Name() string {
	return basicName
}

func (sub *BasicSubscription) RequestLimit() int {
	return basicRequestLimit
}

func (sub *BasicSubscription) TimeFrame() int64 {
	return basicTimeFrame
}

func NewBasicSubscription() *BasicSubscription {
	return &BasicSubscription{basicName, basicRequestLimit, basicTimeFrame}
}

func (sub *PremiumSubscription) Name() string {
	return premiumName
}

func (sub *PremiumSubscription) RequestLimit() int {
	return premiumRequestLimit
}

func (sub *PremiumSubscription) TimeFrame() int64 {
	return premiumTimeFrame
}

func NewPremiumSubscription() *PremiumSubscription {
	return &PremiumSubscription{premiumName, premiumRequestLimit, premiumTimeFrame}
}

type ClientMap struct {
	entries map[string]*Entry
}

func (cm *ClientMap) ValidClientId(clientId string) bool {
	_, ok := cm.entries[clientId]
	return ok
}

const (
	slidingWindow = "sliding"
	fixedWindow   = "window"
)

func WindowHandler(windowType string, handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clientId = r.URL.Query().Get("ClientId")

		defer ResponseMapper(w)

		if clientId == "" {
			panic(NewBadRequestError())
		} else if !validClients.ValidClientId(clientId) {
			panic(NewUnauthorizedRequestError())
		}

		var now = time.Now()
		var invocationTime int64 = 0

		switch windowType {
		case slidingWindow:
			invocationTime = now.UnixMilli()
		case fixedWindow:
			invocationTime = now.Truncate(time.Minute).UnixMilli()
		default:
			panic(NewBadRequestError())
		}

		UpdateEntry(clientId, invocationTime)
		handler.ServeHTTP(w, r)
	})
}

func UpdateEntry(clientId string, invocationTime int64) {
	var entry = validClients.entries[clientId]

	var limit = entry.subscription.RequestLimit()
	var timeFrame = entry.subscription.TimeFrame()

	if entry.startTime == 0 && entry.invocations == 0 { // new entry
		entry.startTime = invocationTime
		entry.lastInvocation = invocationTime
		entry.invocations++
		fmt.Println("New entry for: ", clientId, entry.invocations)
	} else if entry.invocations < limit && invocationTime-entry.startTime < timeFrame { // update
		entry.invocations++
		entry.lastInvocation = invocationTime
		fmt.Println("Update entry for: ", clientId, entry.invocations)
	} else if invocationTime-entry.startTime > timeFrame { // reset
		entry.invocations = 1
		entry.startTime = invocationTime
		entry.lastInvocation = invocationTime
		fmt.Println("Reset entry for: ", clientId, entry.invocations)
	} else {
		panic(NewLimitExceededError())
	}
}

//Leaky Bucket, Fixed Bucket, Custom Rate Limits
