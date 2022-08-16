package apiresponse

import (
	"encoding/json"
	"net/http"
)

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
