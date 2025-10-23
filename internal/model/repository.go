package model

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

var ErrDocumentNotFound = &echo.HTTPError{Message: "could not retrive entry from repository", Code: http.StatusInternalServerError}
var ErrDocumentNotExists = &echo.HTTPError{Internal: errors.New("document does not exist in collection"), Message: "document does not exist in collection", Code: http.StatusBadRequest}

type Repository[T any] interface {
	Get(context.Context, string) (*T, *echo.HTTPError)
	Save(context.Context, *T) *echo.HTTPError
	Update(context.Context, *T) *echo.HTTPError
	Delete(context.Context, string) *echo.HTTPError
	GetAll(context.Context) (*[]T, *echo.HTTPError)
}
