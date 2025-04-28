package httplib

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/kodenkai-labs/go-lib/errlib"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func HandleError(c *gin.Context, err error) {
	var statusCode int
	var message string

	var isClientError bool

	var appErr errlib.AppError
	if errors.As(err, &appErr) {
		// Map domain error to HTTP error
		switch appErr.Code() {
		case errlib.NotFoundCode:
			statusCode = http.StatusNotFound
			message = string(appErr.Slug())
			isClientError = true
		case errlib.InvalidInputCode:
			statusCode = http.StatusBadRequest
			message = string(appErr.Slug())
			isClientError = true
		case errlib.UnauthorizedCode:
			statusCode = http.StatusUnauthorized
			message = string(appErr.Slug())
			isClientError = true
		case errlib.ConflictCode:
			statusCode = http.StatusConflict
			message = string(appErr.Slug())
			isClientError = true
		case errlib.UnprocessableCode:
			statusCode = http.StatusUnprocessableEntity
			message = string(appErr.Slug())
			isClientError = true
		case errlib.InternalCode:
			statusCode = http.StatusInternalServerError
			message = string(appErr.Slug())
		default:
			statusCode = http.StatusInternalServerError
			message = string(errlib.SlugInternal)

			logrus.WithError(appErr).Error("unknown application error")
		}
	} else {
		statusCode = http.StatusInternalServerError
		message = string(errlib.SlugInternal)
	}

	if !isClientError {
		logrus.WithError(err).Error("http error")
	}

	c.JSON(statusCode, NewError(statusCode, message))
}
