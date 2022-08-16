package apiresponse

import "net/http"

type APIError interface {
	StatusCode() int
}

// BadRequestError (400)
type BadRequestError struct {
	Message    string
	statusCode int
}

func (err *BadRequestError) StatusCode() int {
	return err.statusCode
}

func NewBadRequestError() *BadRequestError {
	msg := "Invalid parameters provided."
	return &BadRequestError{Message: msg, statusCode: http.StatusBadRequest}
}

func (err *BadRequestError) Error() string {
	return "Invalid parameters provided."
}

// UnauthorizedRequestError
type UnauthorizedRequestError struct {
	Message    string
	statusCode int
}

func (err *UnauthorizedRequestError) StatusCode() int {
	return err.statusCode
}

func NewUnauthorizedRequestError() *UnauthorizedRequestError {
	msg := "Request was denied."
	return &UnauthorizedRequestError{Message: msg, statusCode: http.StatusUnauthorized}
}

func (err *UnauthorizedRequestError) Error() string {
	return "Unauthorized Request was made."
}

// LimitExceededError (429)
type LimitExceededError struct {
	Message        string
	statusCode     int
	SubscribedRate string
}

func (err *LimitExceededError) StatusCode() int {
	return err.statusCode
}

func NewLimitExceededError() *LimitExceededError {
	msg := "Too many requests. Service will be made available per subscribed rate."
	rate := "5 requests per minute" // TODO: Custom rate limits
	return &LimitExceededError{Message: msg, statusCode: http.StatusTooManyRequests, SubscribedRate: rate}
}

func (err *LimitExceededError) Error() string {
	return "Too many requests were made."
}

// InternalServerError (500)
type InternalServerError struct {
	Message    string
	statusCode int
}

const internalErrorMessage = "Internal Server Error."

func (err *InternalServerError) StatusCode() int {
	return err.statusCode
}

func NewInternalServerError() *InternalServerError {
	msg := internalErrorMessage
	return &InternalServerError{Message: msg, statusCode: http.StatusInternalServerError}
}

func (err *InternalServerError) Error() string {
	return internalErrorMessage
}
