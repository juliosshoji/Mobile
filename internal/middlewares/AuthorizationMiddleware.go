package middlewares

import (
	"Mobile/internal/model/domain"
	"Mobile/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type AuthMiddleware interface {
	AuthorizeMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type authMiddlewareImpl struct {
	customerService      service.CustomerService
	providerService      service.ProviderService
	authorizationService service.AuthorizationService
}

func NewUserAuthMiddleware(providerService service.ProviderService, customerService service.CustomerService, authorizationService service.AuthorizationService) AuthMiddleware {
	return authMiddlewareImpl{
		customerService:      customerService,
		providerService:      providerService,
		authorizationService: authorizationService,
	}
}

func (h authMiddlewareImpl) AuthorizeMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get(echo.HeaderAuthorization) == "" && c.Request().Header.Get("document") == "" {
			return c.JSON(domain.ErrNotAuthenticated.Code, domain.ErrNotAuthenticated.Error())
		}
		authorizationHeader := c.Request().Header.Get(echo.HeaderAuthorization)

		claims, err := h.authorizationService.Authorize(&authorizationHeader)
		if err != nil {
			log.Error().Msg("authorization failed: " + err.Error())
			return c.JSON(err.Code, err.Error())
		}

		_, err = h.customerService.Get(c.Request().Context(), claims.Id)
		if err != nil {
			log.Error().Msg("failed to get user: " + err.Error())
			return c.JSON(err.Code, err.Error())
		}

		return next(c)
	}
}
