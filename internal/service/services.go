package service

import (
	"Mobile/internal/model/customer"
	"Mobile/internal/model/domain"
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
	GetFavorite(context.Context, string) (*[]provider.Provider, *echo.HTTPError)
	AddService(context.Context, string, customer.ServicesDone) *echo.HTTPError
}

type ProviderService interface {
	Get(context.Context, string) (*provider.Provider, *echo.HTTPError)
	Update(context.Context, *provider.Provider) *echo.HTTPError
	Delete(context.Context, string) *echo.HTTPError
	Add(context.Context, *provider.Provider) *echo.HTTPError

	AddSpecialty(context.Context, *provider.Specialty, string) *echo.HTTPError
	GetBySpecialty(context.Context, *provider.Specialty) (*[]provider.Provider, *echo.HTTPError)
	AddProfilePhoto(context.Context, string, string) *echo.HTTPError
}

type ReviewService interface {
	Create(context.Context, *review.Review) *echo.HTTPError
	Delete(context.Context, string) *echo.HTTPError
	Get(context.Context, string) (*review.Review, *echo.HTTPError)

	GetAllBy(context.Context, string, string) (*[]review.Review, *echo.HTTPError)
}

type AuthenticationService interface {
	GenerateToken(context.Context, string) (*string, *echo.HTTPError)
	AuthenticateProvider(context.Context, string, string) (*string, *echo.HTTPError)
	AuthenticateCustomer(context.Context, string, string) (*string, *echo.HTTPError)
}

type AuthorizationService interface {
	Authorize(authHeader *string) (*domain.Claims, *echo.HTTPError)
}
