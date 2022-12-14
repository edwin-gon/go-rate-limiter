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
		var err error
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

		windowEntry := clientMap.Entries[clientId].(*WindowEntry)
		err = UpdateWindowEntry(windowEntry, clientId, invocationTime)

		if err != nil {
			panic(err)
		}
		handler.ServeHTTP(w, r)
	})
}

func UpdateWindowEntry(entry *WindowEntry, clientId string, invocationTime int64) error {

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

//Fixed Bucket ??? having certain number request limit R over time T. Tokens are replensihed at a rate of r over t time.
func TokenBucket(client TokenEntry, invocationTime int64) error {
	// Initial state
	if invocationTime == 0 {
		client.invocations = client.subscription.RequestLimit()
		return nil
	}

	rateAdded := client.subscription.TimeFrame() / int64(client.subscription.RequestLimit()) // time / token
	tokensToAdd := int((invocationTime - client.lastInvocation) / rateAdded)                 // token
	combinedTokens := client.invocations + int(tokensToAdd)

	//Limit Exceeded
	if combinedTokens < 1 {
		return apiresponse.NewLimitExceededError()
	}

	//Update Tokens ??? Reset to limit or increment
	if combinedTokens > client.subscription.RequestLimit() {
		client.invocations = client.subscription.RequestLimit()
	} else {
		client.invocations = combinedTokens
	}

	client.invocations-- // Token to be used for request

	return nil
}

//Leaky Bucket ??? a bucket can be thought of as a queue and r number of requests are processed over t time and once bucket is full no requests are queued
func LeakyBucket(client *TokenEntry, invocationTime int64) error {
	// Check if queue value can be dequeued
	// Attempt to add to the queue if unable return error
	if client.queue.Count() > 0 {
		rateAdded := client.subscription.TimeFrame() / int64(client.subscription.RequestLimit()) // time / token
		requestToDequeue := int((invocationTime - client.lastInvocation) / rateAdded)

		for requestToDequeue > 0 {
			packet := client.queue.Dequeue()
			fmt.Printf("Packet %d dequeued.\n", packet)
			requestToDequeue--
		}
	}

	if wasAdded := client.queue.Enqueue(int(invocationTime)); !wasAdded {
		return apiresponse.NewLimitExceededError()
	}

	fmt.Println("Packet was recieved.")
	return nil
}
