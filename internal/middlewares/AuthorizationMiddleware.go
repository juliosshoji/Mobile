package middlewares

import (
	"Mobile/internal/model/domain"
	"Mobile/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type AuthMiddleware interface {
	AuthorizeCustomerMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	AuthorizeProviderMiddleware(next echo.HandlerFunc) echo.HandlerFunc
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

func (h authMiddlewareImpl) AuthorizeCustomerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get(echo.HeaderAuthorization) == "" && c.FormValue("document") != "" {
			return c.JSON(domain.ErrNotAuthenticated.Code, domain.ErrNotAuthenticated.Error())
		}
		authorizationHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		if authorizationHeader == "" {
			log.Warn().Msg("authorization header is missing")
			return c.JSON(domain.ErrNotAuthenticated.Code, domain.ErrNotAuthenticated.Error())
		}

		claims, err := h.authorizationService.Authorize(&authorizationHeader)
		if err != nil {
			log.Error().Msg("authorization failed: " + err.Error())
			return c.JSON(err.Code, err.Error())
		}

		userResponse, err := h.customerService.Get(c.Request().Context(), claims.Id)
		if err != nil {
			log.Error().Msg("failed to get user: " + err.Error())
			return c.JSON(err.Code, err.Error())
		}

		if userResponse.Name == "admin" {
			log.Info().Msg("admin user authorized")
			return next(c)
		}

		if c.Param("document") != "" {
			if c.Param("document") != claims.Id {
				log.Error().Msg("user id doesn't match with claims id")
				return c.JSON(domain.ErrUserIDNotMatch.Code, domain.ErrUserIDNotMatch.Internal)
			}
			return next(c)
		}
		if c.FormValue("document") != claims.Id {
			log.Error().Msg("user id doesn't match with claims id")
			return c.JSON(domain.ErrUserIDNotMatch.Code, domain.ErrUserIDNotMatch.Internal)
		}
		return next(c)
	}
}

func (h authMiddlewareImpl) AuthorizeProviderMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get(echo.HeaderAuthorization) == "" && c.FormValue("user_id") != "" {
			return c.JSON(domain.ErrNotAuthenticated.Code, domain.ErrNotAuthenticated.Error())
		}
		authorizationHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		if authorizationHeader == "" {
			log.Warn().Msg("authorization header is missing")
			return c.JSON(domain.ErrNotAuthenticated.Code, domain.ErrNotAuthenticated.Error())
		}

		claims, err := h.authorizationService.Authorize(&authorizationHeader)
		if err != nil {
			log.Error().Msg("authorization failed: " + err.Error())
			return c.JSON(err.Code, err.Error())
		}

		userResponse, err := h.providerService.Get(c.Request().Context(), claims.Id)
		if err != nil {
			log.Error().Msg("failed to get user: " + err.Error())
			return c.JSON(err.Code, err.Error())
		}

		if userResponse.Name == "admin" {
			log.Info().Msg("admin user authorized")
			return next(c)
		}

		if c.Param("document") != "" {
			if c.Param("document") != claims.Id {
				log.Error().Msg("user id doesn't match with claims id")
				return c.JSON(domain.ErrUserIDNotMatch.Code, domain.ErrUserIDNotMatch.Internal)
			}
			return next(c)
		}
		if c.FormValue("document") != "" {
			if c.FormValue("document") != claims.Id {
				log.Error().Msg("user id doesn't match with claims id")
				return c.JSON(domain.ErrUserIDNotMatch.Code, domain.ErrUserIDNotMatch.Internal)
			}
			return next(c)
		}
		return next(c)
	}
}
