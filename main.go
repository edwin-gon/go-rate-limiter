package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/edwin-gon/go-rate-limiter/ratelimiter"
)

// Include Logger, DI (wire)
// Leaky Bucket, Fixed Bucket, Custom Rate Limits, Leveraging DynamoDB to do quick reads, What if stale data is acquired

var validClients *ratelimiter.ClientMap = &ratelimiter.ClientMap{Entries: map[string]ratelimiter.Entry{
	"SLIDING": ratelimiter.NewWindowEntry(ratelimiter.NewBasicSubscription()),
	"FIXED":   ratelimiter.NewWindowEntry(ratelimiter.NewBasicSubscription()),
	"BUCKET":  ratelimiter.NewTokenEntry(ratelimiter.NewBasicSubscription())}}

func main() {

	http.HandleFunc("/sliding/client", ratelimiter.WindowHandler(ratelimiter.SlidingWindow, validClients, http.HandlerFunc(getClientName)))

	http.HandleFunc("/fixed/client", ratelimiter.WindowHandler(ratelimiter.FixedWindow, validClients, http.HandlerFunc(getClientName)))

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
}
