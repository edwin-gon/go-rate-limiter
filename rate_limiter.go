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
	GetName() string
	GetRequestLimit() int
	GetTimeFrame() int64
}

type BasicSubscription struct {
	name                    string
	requestLimit, timeFrame int64
}

type PremiumSubscription struct {
	name                    string
	requestLimit, timeFrame int64
}

func (sub *BasicSubscription) GetName() string {
	return "Basic"
}

func (sub *BasicSubscription) GetRequestLimit() int {
	return 5
}

func (sub *BasicSubscription) GetTimeFrame() int64 {
	return 60000
}

func NewBasicSubscription() *BasicSubscription {
	return &BasicSubscription{"BASIC", 5, 60000}
}

func (sub *PremiumSubscription) GetName() string {
	return "Premium"
}

func (sub *PremiumSubscription) GetRequestLimit() int {
	return 20
}

func (sub *PremiumSubscription) GetTimeFrame() int64 {
	return 60000
}

func NewPremiumSubscription() *PremiumSubscription {
	return &PremiumSubscription{"PREMIUM", 20, 60000}
}

type ClientMap struct {
	entries map[string]*Entry
}

func (cm *ClientMap) ValidClientId(clientId string) bool {
	_, ok := cm.entries[clientId]
	return ok
}

func SlidingWindow(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clientId = r.URL.Query().Get("ClientId")

		defer ResponseMapper(w)

		if clientId == "" {
			panic(NewBadRequestError())
		} else if !validClients.ValidClientId(clientId) {
			panic(NewUnauthorizedRequestError())
		}

		var now = time.Now()
		var invocationTime = now.UnixMilli()

		var entry = validClients.entries[clientId]

		var limit = entry.subscription.GetRequestLimit()
		var timeFrame = entry.subscription.GetTimeFrame()

		if entry.startTime == 0 && entry.invocations == 0 { // new entry
			entry.startTime = invocationTime
			entry.lastInvocation = invocationTime
			entry.invocations++
			fmt.Println("New entry for: ", clientId, entry.invocations)
		} else if entry.invocations < limit && invocationTime-entry.startTime < timeFrame { // update
			entry.invocations++
			entry.lastInvocation = invocationTime
			fmt.Println("Update entry for: ", clientId, entry.invocations)
		} else if invocationTime-entry.startTime > 60000 { // reset
			entry.invocations = 1
			entry.startTime = invocationTime
			entry.lastInvocation = invocationTime
			fmt.Println("Reset entry for: ", clientId, entry.invocations)
		} else {
			panic(NewLimitExceededError())
		}
		handler.ServeHTTP(w, r)

	})
}

// Fixed Window Leaky Bucket, Fixed Bucket, Custom Rate Limits
