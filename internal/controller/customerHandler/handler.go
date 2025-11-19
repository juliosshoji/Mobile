package customerhandler

import (
	"Mobile/internal/controller"
	"Mobile/internal/model/customer"
	"Mobile/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type CustomerHandler interface {
	controller.Handler

	AddFavorite(echo.Context) error
	GetFavorite(echo.Context) error
}

type handlerImpl struct {
	customerService service.CustomerService
}

func NewCustomerHandler(customerService service.CustomerService) CustomerHandler {
	return handlerImpl{
		customerService: customerService,
	}
}

func (ref handlerImpl) Post(c echo.Context) error {
	var newCustomer customer.Customer
	if err := c.Bind(&newCustomer); err != nil {
		log.Err(err).Msg("error binding customer")
		if httpErr := err.(*echo.HTTPError); httpErr != nil {
			return c.NoContent(httpErr.Code)
		}
		return c.NoContent(http.StatusBadRequest)
	}

	if err := ref.customerService.Add(c.Request().Context(), &newCustomer); err != nil {
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

	if err := ref.customerService.Delete(c.Request().Context(), customerId); err != nil {
		log.Err(err.Unwrap()).Msg("error deleting customer")
		return c.NoContent(err.Code)
	}

	return c.NoContent(http.StatusOK)
}
func (ref handlerImpl) Put(c echo.Context) error {

	var customer customer.Customer

	customer.Document = c.Param("document")
	if customer.Document == "" {
		log.Warn().Msg("no param provided")
		return c.NoContent(http.StatusBadRequest)
	}

	if err := c.Bind(&customer); err != nil {
		log.Err(err).Msg("error binding customer on update")
		if httpErr := err.(*echo.HTTPError); httpErr != nil {
			return c.NoContent(httpErr.Code)
		}
		return c.NoContent(http.StatusBadRequest)
	}

	if err := ref.customerService.Update(c.Request().Context(), &customer); err != nil {
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

	customer, err := ref.customerService.Get(c.Request().Context(), customerId)
	if err != nil {
		log.Err(err.Unwrap()).Msg("error returning customer")
		return c.NoContent(err.Code)
	}

	return c.JSON(http.StatusOK, customer)
}

func (ref handlerImpl) AddFavorite(c echo.Context) error {

	type AddFavoriteRequest struct {
		ProviderID string `json:"provider_id"`
	}

	customerId := c.Param("document")
	var request AddFavoriteRequest
	if err := c.Bind(&request); err != nil {
		log.Error().Err(err).Msg("error binding add favorite request")
		if httpErr := err.(*echo.HTTPError); httpErr != nil {
			return c.NoContent(httpErr.Code)
		}
		return c.NoContent(http.StatusBadRequest)
	}

	if request.ProviderID == "" || customerId == "" {
		log.Error().Msg("parameters are incomplete (customer: " + customerId + ") (provider: " + request.ProviderID + ")")
		return c.NoContent(http.StatusBadRequest)
	}

	if err := ref.customerService.AddFavorite(c.Request().Context(), customerId, request.ProviderID); err != nil {
		log.Err(err.Unwrap()).Msg("error adding provider to customer's favorite")
		return c.NoContent(err.Code)
	}

	return c.NoContent(http.StatusOK)
}

func (ref handlerImpl) GetFavorite(c echo.Context) error {
	customerId := c.Param("document")
	if customerId == "" {
		log.Warn().Msg("no param provided")
		return c.NoContent(http.StatusBadRequest)
	}

	favorites, err := ref.customerService.GetFavorite(c.Request().Context(), customerId)
	if err != nil {
		log.Err(err.Unwrap()).Msg("error getting customer's favorite")
		return c.NoContent(err.Code)

	}

	return c.JSON(http.StatusOK, *favorites)
}
