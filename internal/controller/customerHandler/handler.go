package customerhandler

import (
	"Mobile/internal/model/customer"
	"Mobile/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type CustomerHandler interface {
	Register(echo.Context) error
	Delete(echo.Context) error
	Update(echo.Context) error
	Get(echo.Context) error

	AddFavorite(echo.Context) error
}

type handlerImpl struct {
	customerService service.CustomerService
}

func NewCustomerHandler(customerService service.CustomerService) CustomerHandler {
	return handlerImpl{
		customerService: customerService,
	}
}

func (ref handlerImpl) Register(c echo.Context) error {
	var newCustomer customer.Customer
	if err := c.Bind(&newCustomer); err != nil {
		log.Err(err).Msg("error binding customer")
		if httpErr := err.(*echo.HTTPError); httpErr != nil {
			return c.NoContent(httpErr.Code)
		}
		return c.NoContent(http.StatusBadRequest)
	}

	if err := ref.customerService.Add(&newCustomer); err != nil {
		log.Err(err.Unwrap()).Msg("error at customer service")
		return c.NoContent(err.Code)
	}

	return c.NoContent(http.StatusOK)
}

func (ref handlerImpl) Delete(c echo.Context) error {

	customerId := c.Param("document")
	if customerId == "" {
		log.Warn().Msg("no param provided")
		return c.NoContent(http.StatusBadRequest)
	}

	if err := ref.customerService.Delete(customerId); err != nil {
		log.Err(err.Unwrap()).Msg("error deleting customer")
		return c.NoContent(err.Code)
	}

	return c.NoContent(http.StatusOK)
}
func (ref handlerImpl) Update(c echo.Context) error {

	var customer customer.Customer
	if err := c.Bind(&customer); err != nil {
		log.Err(err).Msg("error binding customer on update")
		if httpErr := err.(*echo.HTTPError); httpErr != nil {
			return c.NoContent(httpErr.Code)
		}
		return c.NoContent(http.StatusBadRequest)
	}

	if err := ref.customerService.Update(&customer); err != nil {
		log.Err(err.Unwrap()).Msg("error updating customer")
		return c.NoContent(err.Code)
	}

	return c.NoContent(http.StatusOK)
}
func (ref handlerImpl) Get(c echo.Context) error {

	customerId := c.Param("document")
	if customerId == "" {
		log.Warn().Msg("no param provided")
		return c.NoContent(http.StatusBadRequest)
	}

	customer, err := ref.customerService.Get(customerId)
	if err != nil {
		log.Err(err.Unwrap()).Msg("error returning customer")
		return c.NoContent(err.Code)
	}

	return c.JSON(http.StatusOK, customer)
}

func (ref handlerImpl) AddFavorite(c echo.Context) error {

	customerId := c.FormValue("customer_id")
	providerId := c.FormValue("provider_id")

	if providerId == "" || customerId == "" {
		log.Error().Msg("parameters are incomplete (customer: " + customerId + ") (provider: " + providerId + ")")
		return c.NoContent(http.StatusBadRequest)
	}

	if err := ref.customerService.AddFavorite(customerId, providerId); err != nil {
		log.Err(err.Unwrap()).Msg("error adding provider to customer's favorite")
		return c.NoContent(err.Code)
	}

	return c.NoContent(http.StatusOK)
}
