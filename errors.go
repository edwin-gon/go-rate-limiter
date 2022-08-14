package main

import "net/http"

type APIError interface {
	GetStatusCode() int
}

// LimitExceededError
type LimitExceededError struct {
	Message        string
	StatusCode     int
	SubscribedRate string
}

func (err *LimitExceededError) GetStatusCode() int {
	return err.StatusCode
}

func NewLimitExceededError() *LimitExceededError {
	msg := "Too many requests. Service will be made available per subscribed rate."
	rate := "5 requests per minute" // TODO: Custom rate limits
	return &LimitExceededError{Message: msg, StatusCode: http.StatusTooManyRequests, SubscribedRate: rate}
}

func (err *LimitExceededError) Error() string {
	return "Too many requests were made."
}

// UnauthorizedRequestError
type UnauthorizedRequestError struct {
	Message    string
	StatusCode int
}

func (err *UnauthorizedRequestError) GetStatusCode() int {
	return err.StatusCode
}

func NewUnauthorizedRequestError() *UnauthorizedRequestError {
	msg := "Request was denied."
	return &UnauthorizedRequestError{Message: msg, StatusCode: http.StatusUnauthorized}
}

func (err *UnauthorizedRequestError) Error() string {
	return "Unauthorized Request was made."
}

// InternalServerError
type InternalServerError struct {
	Message    string
	StatusCode int
}

func (err *InternalServerError) GetStatusCode() int {
	return err.StatusCode
}

func NewInternalServerError() *InternalServerError {
	msg := "Internal Server Error."
	return &InternalServerError{Message: msg, StatusCode: http.StatusInternalServerError}
}

func (err *InternalServerError) Error() string {
	return "Internal Server Error."
}
