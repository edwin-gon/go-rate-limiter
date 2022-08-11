package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Entry struct {
	startTime, lastInvocation int64
	invocations               int64
}

type ClientMap struct {
	entries map[string]*Entry
}

var validClients *ClientMap = &ClientMap{map[string]*Entry{"VALID": {}}}

func (cm *ClientMap) ValidClientId(clientId string) bool {
	_, ok := cm.entries[clientId]
	return ok
}

func RateLimiter(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clientId = r.URL.Query().Get("ClientId")

		if !validClients.ValidClientId(clientId) {
			w.WriteHeader(http.StatusUnauthorized)
			panic(errors.New("Invalid Client."))
		}

		var now = time.Now()
		var invocationTime = now.UnixMilli()

		var entry = validClients.entries[clientId]
		// Invalid client provided
		// Check start time and invocations — if start time is 0 allow through (start state),
		// (Mid state update) if invocations is less than 5 and start time has not passed 1 minute time span update last invocation and counter
		// (Mid state reset) if invocations is less than 5 and start time has passed 1 minute time span
		// (Mid state error) if invocation is equal than 5 and start time has is not passed the 1 minute time span

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
			panic(errors.New("You have hit your limity of 5 requests per minute."))
		}
		handler.ServeHTTP(w, r)
	})
}

func main() {

	http.HandleFunc("/client-name", RateLimiter(http.HandlerFunc(getClientName)))

	err := http.ListenAndServe(":5050", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error occured: %s\n", err)
		os.Exit(1)
	}
}

func getClientName(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /client-name request\n")
	w.WriteHeader(http.StatusOK)
}
