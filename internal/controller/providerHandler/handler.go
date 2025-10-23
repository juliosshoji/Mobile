package providerHandler

import (
	"Mobile/internal/controller"
	"Mobile/internal/model/provider"
	"Mobile/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ProviderHandler interface {
	controller.Handler

	GetContact(echo.Context) error
	AddSpecialty(echo.Context) error
	GetBySpecialty(echo.Context) error
}

type providerHandlerImpl struct {
	providerService service.ProviderService
}

func NewProviderHandler(service service.ProviderService) ProviderHandler {
	return providerHandlerImpl{
		providerService: service,
	}
}

func (ref providerHandlerImpl) Post(c echo.Context) error {
	var newCustomer provider.Provider
	if err := c.Bind(&newCustomer); err != nil {
		log.Err(err).Msg("error binding customer")
		if httpErr := err.(*echo.HTTPError); httpErr != nil {
			return c.NoContent(httpErr.Code)
		}
		return c.NoContent(http.StatusBadRequest)
	}

	if err := ref.providerService.Add(c.Request().Context(), &newCustomer); err != nil {
		log.Err(err.Unwrap()).Msg("error at customer service")
		return c.NoContent(err.Code)
	}

	return c.NoContent(http.StatusOK)
}

func (ref providerHandlerImpl) Get(c echo.Context) error {
	providerId := c.Param("document")
	if providerId == "" {
		log.Warn().Msg("no param provided")
		return c.NoContent(http.StatusBadRequest)
	}

	customer, err := ref.providerService.Get(c.Request().Context(), providerId)
	if err != nil {
		log.Err(err.Unwrap()).Msg("error returning provider")
		return c.NoContent(err.Code)
	}

	return c.JSON(http.StatusOK, customer)
}
func (ref providerHandlerImpl) Put(c echo.Context) error {
	var customer provider.Provider
	if err := c.Bind(&customer); err != nil {
		log.Err(err).Msg("error binding provider on update")
		if httpErr := err.(*echo.HTTPError); httpErr != nil {
			return c.NoContent(httpErr.Code)
		}
		return c.NoContent(http.StatusBadRequest)
	}

	if err := ref.providerService.Update(c.Request().Context(), &customer); err != nil {
		log.Err(err.Unwrap()).Msg("error updating provider")
		return c.NoContent(err.Code)
	}

	return c.NoContent(http.StatusOK)
}
func (ref providerHandlerImpl) Delete(c echo.Context) error {
	customerId := c.Param("document")
	if customerId == "" {
		log.Warn().Msg("no param provided")
		return c.NoContent(http.StatusBadRequest)
	}

	if err := ref.providerService.Delete(c.Request().Context(), customerId); err != nil {
		log.Err(err.Unwrap()).Msg("error deleting provider")
		return c.NoContent(err.Code)
	}

	return c.NoContent(http.StatusOK)
}

func (ref providerHandlerImpl) GetContact(c echo.Context) error {
	document := c.Param("document")
	if document == "" {
		log.Warn().Msg("no param provided")
		return c.NoContent(http.StatusBadRequest)
	}

	provider, err := ref.providerService.Get(c.Request().Context(), document)
	if err != nil {
		log.Err(err.Unwrap()).Msg("error deleting provider")
		return c.NoContent(err.Code)
	}
	return c.JSON(http.StatusOK, provider.ContactAddress)
}
func (ref providerHandlerImpl) AddSpecialty(c echo.Context) error {

	var specialty provider.Specialty
	if err := c.Bind(&specialty); err != nil {
		log.Err(err).Msg("error binding customer")
		if httpErr := err.(*echo.HTTPError); httpErr != nil {
			return c.NoContent(httpErr.Code)
		}
		return c.NoContent(http.StatusBadRequest)
	}

	providerDocument := c.Param("document")
	if providerDocument == "" {
		log.Warn().Msg("no param provided")
		return c.NoContent(http.StatusBadRequest)
	}

	if err := ref.providerService.AddSpecialty(c.Request().Context(), &specialty, providerDocument); err != nil {
		log.Err(err.Unwrap()).Msg("error at customer service")
		return c.NoContent(err.Code)
	}

	return c.NoContent(http.StatusOK)
}

func (ref providerHandlerImpl) GetBySpecialty(c echo.Context) error {
	parameter := c.Param("specialty")
	if parameter == "" {
		log.Warn().Msg("no param provided")
		return c.NoContent(http.StatusBadRequest)
	}

	specialty := provider.Specialty(parameter)
	if !specialty.IsValid() {
		log.Warn().Msg("invalid specialty provided")
		return c.NoContent(http.StatusBadRequest)
	}

	providers, err := ref.providerService.GetBySpecialty(c.Request().Context(), &specialty)
	if err != nil {
		log.Error().Err(err.Unwrap())
		return err
	}

	return c.JSON(http.StatusOK, providers)
}
