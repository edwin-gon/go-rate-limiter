package main

import (
	"fmt"
	"net/http"
	"time"
)

type Entry struct {
	startTime, lastInvocation int64
	invocations               int64
}

type ClientMap struct {
	entries map[string]*Entry
}

func (cm *ClientMap) ValidClientId(clientId string) bool {
	_, ok := cm.entries[clientId]
	return ok
}

func RateLimiter(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clientId = r.URL.Query().Get("ClientId")

		defer ResponseMapper(w)

		if !validClients.ValidClientId(clientId) {
			panic(NewUnauthorizedRequestError())
		}

		var now = time.Now()
		var invocationTime = now.UnixMilli()

		var entry = validClients.entries[clientId]

		if entry.startTime == 0 && entry.invocations == 0 { // new entry
			entry.startTime = invocationTime
			entry.lastInvocation = invocationTime
			entry.invocations++
			fmt.Println("New entry for: ", clientId, entry.invocations)
		} else if entry.invocations < 5 && invocationTime-entry.startTime < 60000 { // update
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
