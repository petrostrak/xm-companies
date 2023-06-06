package utils

import (
	"errors"
	"log"
	"net/http"
	"os"
)

var (
	ErrInvalidIDParam = errors.New("invalid id parameter")
	ErrEmptyBody      = errors.New("body must not be empty")
	ErrBadJSON        = errors.New("body contains badly-formed JSON")
	ErrSingleJSON     = errors.New("body must only contain a single JSON value")
)

func LogError(err error) {
	logger := log.New(os.Stdout, "[ERROR] ", log.Ldate|log.Ltime)
	logger.Println(err)
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := Envelope{"error": message}

	err := WriteJSON(w, status, env, nil)
	if err != nil {
		LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	LogError(err)
	message := "the server encountered a problem and could not process your request"
	ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}
