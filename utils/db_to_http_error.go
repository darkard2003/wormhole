package utils

import (
	"log"
	"net/http"

	"github.com/darkard2003/wormhole/services/db"
)

type HttpError struct {
	Code     int    `json:"code"`
	Response string `json:"response"`
}

func DBToHttpError(err error) *HttpError {
	switch err.(type) {
	case *db.AlreadyExistsError:
		return &HttpError{
			Code:     http.StatusConflict,
			Response: err.Error(),
		}
	case *db.NotFoundError:
		return &HttpError{
			Code:     http.StatusNotFound,
			Response: err.Error(),
		}
	case *db.ValidationError:
		return &HttpError{
			Code:     http.StatusBadRequest,
			Response: err.Error(),
		}
	default:
		log.Println("Unknown error:", err)
		return &HttpError{
			Code:     http.StatusInternalServerError,
			Response: "Internal server error",
		}
	}
}
