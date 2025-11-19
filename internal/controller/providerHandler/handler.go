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
	AddProfilePhoto(c echo.Context) error
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
	var provider provider.Provider
	provider.Document = c.Param("document")

	if err := c.Bind(&provider); err != nil {
		log.Err(err).Msg("error binding provider on update")
		if httpErr := err.(*echo.HTTPError); httpErr != nil {
			return c.NoContent(httpErr.Code)
		}
		return c.NoContent(http.StatusBadRequest)
	}

	if err := ref.providerService.Update(c.Request().Context(), &provider); err != nil {
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

	type AddSpecialtyObject struct {
		Specialty string `json:"specialty"`
	}

	var requestBody AddSpecialtyObject
	if err := c.Bind(&requestBody); err != nil {
		log.Err(err).Msg("error binding speialty")
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

	specialty := provider.Specialty(requestBody.Specialty)

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

func (ref providerHandlerImpl) AddProfilePhoto(c echo.Context) error {
	providerDoc := c.Param(":document")
	if providerDoc == "" {
		log.Warn().Msg("no param provided")
		return c.NoContent(http.StatusBadRequest)
	}

	file := c.FormValue("profile_photo")
	if err := ref.providerService.AddProfilePhoto(c.Request().Context(), providerDoc, file); err != nil {
		log.Err(err.Unwrap()).Msg("error adding profile photo")
		return c.NoContent(err.Code)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "profile photo added successfully"})
}
