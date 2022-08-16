package ratelimiter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/edwin-gon/go-rate-limiter/apiresponse"
)

const (
	SlidingWindow = "sliding"
	FixedWindow   = "window"
)

func WindowHandler(windowType string, clientMap *ClientMap, handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clientId = r.URL.Query().Get("ClientId")

		defer apiresponse.ResponseMapper(w)

		if clientId == "" {
			panic(apiresponse.NewBadRequestError())
		} else if !clientMap.ValidClientId(clientId) {
			panic(apiresponse.NewUnauthorizedRequestError())
		}

		var now = time.Now()
		var invocationTime int64 = 0

		switch windowType {
		case SlidingWindow:
			invocationTime = now.UnixMilli()
		case FixedWindow:
			invocationTime = now.Truncate(time.Minute).UnixMilli()
		default:
			panic(apiresponse.NewBadRequestError())
		}

		err := UpdateWindowEntry(clientMap.Entries[clientId], clientId, invocationTime)
		if err != nil {
			panic(err)
		}
		handler.ServeHTTP(w, r)
	})
}

func UpdateWindowEntry(entry *Entry, clientId string, invocationTime int64) error {

	var limit = entry.subscription.RequestLimit()
	var timeFrame = entry.subscription.TimeFrame()
	var err error
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
		err = apiresponse.NewLimitExceededError()
	}
	return err
}

//Fixed Bucket — having certain number request limit R over time T. Tokens are replensihed at a rate of r over t time.
func TokenBucket(client *Entry, invocationTime int64) error {
	var err error
	if invocationTime == 0 {
		client.invocations = client.subscription.RequestLimit()
	} else if client.invocations < client.subscription.RequestLimit() {
		//Determine number of tokens to add to bucket and should not exceed limit
		rateAdded := client.subscription.TimeFrame() / int64(client.subscription.RequestLimit())
		tokensToAdd := (client.lastInvocation - invocationTime) / rateAdded
		client.invocations += int(tokensToAdd)

		if client.invocations > client.subscription.RequestLimit() {
			client.invocations = client.subscription.RequestLimit()
		}

		if client.invocations < 1 {
			err = apiresponse.NewLimitExceededError()
		} else {
			client.invocations--
		}
	}
	return err
}

//Leaky Bucket — a bucket can be thought of as a queue and r number of requests are processed over t time and once bucket is full no requests are queued
// Process requests every 12 seconds
