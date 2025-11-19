package authenticationHandler

import (
	"Mobile/internal/model/domain"
	"Mobile/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type AuthenticationHandler interface {
	Authenticate(c echo.Context) error
}

type authorizationHandlerImpl struct {
	authenticationService service.AuthenticationService
}

func NewAuthenticationHandler(authenticationService service.AuthenticationService) AuthenticationHandler {
	return authorizationHandlerImpl{
		authenticationService: authenticationService,
	}
}

func (ref authorizationHandlerImpl) Authenticate(c echo.Context) error {

	var auth domain.Authentication
	if err := c.Bind(&auth); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid authentication format")
	}

	log.Debug().Msg("document: " + auth.Document + " password: " + auth.Password)

	token, err := ref.authenticationService.AuthenticateCustomer(c.Request().Context(), auth.Document, auth.Password)
	if err != nil {
		return c.JSON(err.Code, "authentication failed")
	}

	return c.JSON(http.StatusAccepted, map[string]string{"token": *token})
}
