package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

// HTTPError is api error response
type HTTPError struct {
	Status        int    `json:"status"`
	Message       string `json:"message"`
	InternalError error  `json:"-"`
}

func (e *HTTPError) withInternalError(err error) *HTTPError {
	e.InternalError = err
	return e
}

func (e *HTTPError) Error() string {
	if e.InternalError != nil {
		return e.InternalError.Error()
	}
	return fmt.Sprintf("%d: %s", e.Status, e.Message)
}

func httpError(status int, format string, args ...any) *HTTPError {
	return &HTTPError{
		Status:  status,
		Message: fmt.Sprintf(format, args...),
	}
}

func notFoundError(format string, args ...any) *HTTPError {
	return httpError(http.StatusNotFound, format, args...)
}

func badRequestError(format string, args ...any) *HTTPError {
	return httpError(http.StatusBadRequest, format, args...)
}

func unauthorizedError(format string, args ...any) *HTTPError {
	return httpError(http.StatusUnauthorized, format, args...)
}

func internalError(format string, args ...any) *HTTPError {
	return httpError(http.StatusInternalServerError, format, args...)
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	switch e := err.(type) {
	case *HTTPError:
		if e.Status >= http.StatusInternalServerError {
			logrus.Println("ERROR", e.Error())
			e.Message = "internal server error"
		}
		w.WriteHeader(e.Status)
		render.Respond(w, r, e)

	default:
		logrus.Printf("Unknown error: %s", e.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":500,"message":"Internal server error"}`))
	}
}
