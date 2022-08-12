package main

import (
	"encoding/json"
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

type APIError interface {
	GetMessage() string
	GetStatusCode() int
}

type LimitExceededError struct {
	message        string
	statusCode     int
	subscribedRate string
}

func (err *LimitExceededError) GetMessage() string {
	return err.message
}

func (err *LimitExceededError) GetStatusCode() int {
	return err.statusCode
}

func NewLimitExceededError() *LimitExceededError {
	msg := "Too many requests were made. Regular service will be made available after specified rate time frame passes."
	rate := "5 requests per minute" // TODO: Custom rate limits
	return &LimitExceededError{message: msg, statusCode: http.StatusTooManyRequests, subscribedRate: rate}
}

func (err *LimitExceededError) Error() string {
	return "Too many requests were made."
}

type UnauthorizedRequestError struct {
	message    string
	statusCode int
}

func (err *UnauthorizedRequestError) GetMessage() string {
	return err.message
}

func (err *UnauthorizedRequestError) GetStatusCode() int {
	return err.statusCode
}

func NewUnauthorizedRequestError() *UnauthorizedRequestError {
	msg := "Request was denied."
	return &UnauthorizedRequestError{message: msg, statusCode: http.StatusUnauthorized}
}

func (err *UnauthorizedRequestError) Error() string {
	return "Unauthorized Request was made"
}

func ResponseMapper(w http.ResponseWriter) { // TODO: Test if I can create my own responsewriter
	if err := recover(); err != nil {

		switch err.(type) {
		case *UnauthorizedRequestError:
			v := err.(*UnauthorizedRequestError)
			WriteResponse(w, v)
		case *LimitExceededError:
			v := err.(*LimitExceededError)
			WriteResponse(w, v)
		default:
			fmt.Printf("Type %T\n", err)
			fmt.Println("Error Occured : ", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func WriteResponse(w http.ResponseWriter, payload APIError) {
	fmt.Println(payload)
	v, _ := json.Marshal(payload)
	w.WriteHeader(payload.GetStatusCode())
	w.Write(v)
}

// CreateResponseMapper for the API — based on the error type respond back to the server with the appropriate response
// Include Logger, DI (wire)
// Leaky Bucket, Fixed Bucket, Custom Rate Limits, Leveraging DynamoDB to do quick reads, What if stale data is acquired

var validClients *ClientMap = &ClientMap{map[string]*Entry{"VALID": {}}}

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
