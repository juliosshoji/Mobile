package service

import (
	"Mobile/internal/model/customer"
	"Mobile/internal/model/provider"
	"Mobile/internal/model/review"

	"github.com/labstack/echo/v4"
)

type CustomerService interface {
	Get(string) (*customer.Customer, *echo.HTTPError)
	Update(*customer.Customer) *echo.HTTPError
	Delete(string) *echo.HTTPError
	Add(*customer.Customer) *echo.HTTPError

	AddFavorite(string, string) *echo.HTTPError
}

type ProviderService interface {
	Get(string) *provider.Provider
	Update(*provider.Provider) *echo.HTTPError
	Delete(string) *echo.HTTPError
	Add(*provider.Provider) *echo.HTTPError

	AddSpecialty(provider.Specialty, string) *echo.HTTPError
}

type ReviewService interface {
	Create(*review.Review) *echo.HTTPError
	Delete(string) *echo.HTTPError
	Get(string) *review.Review
}
