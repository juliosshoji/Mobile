package authorizationhandler

import (
	"Mobile/internal/model/domain"
	"Mobile/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthotenticationHandler interface {
	Authenticate(c echo.Context) error
}

type authorizationHandlerImpl struct {
	authenticationService service.AuthenticationService
}

func NewAuthenticationHandler(authenticationService service.AuthenticationService) AuthotenticationHandler {
	return authorizationHandlerImpl{
		authenticationService: authenticationService,
	}
}

func (ref authorizationHandlerImpl) Authenticate(c echo.Context) error {

	if c.Request().Header.Get("User-Type") == "customer" {

		var auth domain.Authentication
		if err := c.Bind(&auth); err != nil {
			return c.JSON(http.StatusBadRequest, "invalid authentication format")
		}

		token, err := ref.authenticationService.AuthenticateCustomer(c.Request().Context(), auth.Document, auth.Password)
		if err != nil {
			return c.JSON(err.Code, "authentication failed")
		}

		return c.JSON(http.StatusAccepted, map[string]string{"token": *token})

	}
	if c.Request().Header.Get("User-Type") == "provider" {

		var auth domain.Authentication
		if err := c.Bind(&auth); err != nil {
			return c.JSON(http.StatusBadRequest, "invalid authentication format")
		}

		token, err := ref.authenticationService.AuthenticateProvider(c.Request().Context(), auth.Document, auth.Password)
		if err != nil {
			return c.JSON(err.Code, "authentication failed")
		}

		return c.JSON(http.StatusAccepted, map[string]string{"token": *token})
	}

	return c.JSON(http.StatusBadRequest, "invalid user type")
}
