package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

// Include Logger, DI (wire)
// Leaky Bucket, Fixed Bucket, Custom Rate Limits, Leveraging DynamoDB to do quick reads, What if stale data is acquired

var validClients *ClientMap = &ClientMap{map[string]*Entry{
	"VALID": {0, 0, 0, NewBasicSubscription()},
	"FIXED": {0, 0, 0, NewBasicSubscription()}}}

func main() {

	http.HandleFunc("/sliding/client", WindowHandler(slidingWindow, http.HandlerFunc(getClientName)))

	http.HandleFunc("/fixed/client", WindowHandler(fixedWindow, http.HandlerFunc(getClientName)))

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

func ResponseMapper(w http.ResponseWriter) {
	if err := recover(); err != nil {
		var res APIError
		switch err.(type) {
		case *BadRequestError:
			res = err.(*BadRequestError)
		case *UnauthorizedRequestError:
			res = err.(*UnauthorizedRequestError)
		case *LimitExceededError:
			res = err.(*LimitExceededError)
		default:
			res = NewInternalServerError()
		}
		WriteResponse(w, res)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func WriteResponse(w http.ResponseWriter, error APIError) {
	w.WriteHeader(error.StatusCode())
	json.NewEncoder(w).Encode(error)
}
