package service

import (
	"Mobile/internal/model/customer"
	"Mobile/internal/model/provider"
	"Mobile/internal/model/review"
	"context"

	"github.com/labstack/echo/v4"
)

type CustomerService interface {
	Get(context.Context, string) (*customer.Customer, *echo.HTTPError)
	Update(context.Context, *customer.Customer) *echo.HTTPError
	Delete(context.Context, string) *echo.HTTPError
	Add(context.Context, *customer.Customer) *echo.HTTPError

	AddFavorite(context.Context, string, string) *echo.HTTPError
}

type ProviderService interface {
	Get(context.Context, string) (*provider.Provider, *echo.HTTPError)
	Update(context.Context, *provider.Provider) *echo.HTTPError
	Delete(context.Context, string) *echo.HTTPError
	Add(context.Context, *provider.Provider) *echo.HTTPError

	AddSpecialty(context.Context, *provider.Specialty, string) *echo.HTTPError
	GetBySpecialty(context.Context, *provider.Specialty) (*[]provider.Provider, *echo.HTTPError)
}

type ReviewService interface {
	Create(context.Context, *review.Review) *echo.HTTPError
	Delete(context.Context, string) *echo.HTTPError
	Get(context.Context, string) (*review.Review, *echo.HTTPError)

	GetAllBy(context.Context, string, string) (*[]review.Review, *echo.HTTPError)
}
