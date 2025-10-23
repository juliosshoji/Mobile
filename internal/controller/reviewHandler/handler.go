package reviewHandler

import (
	"Mobile/internal/model/review"
	"Mobile/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ReviewHandler interface {
	Post(echo.Context) error
	Get(echo.Context) error
	Delete(echo.Context) error

	GetAllBy(echo.Context) error
}

type reviewHandlerImpl struct {
	reviewService service.ReviewService
}

func NewReviewHandler(service service.ReviewService) ReviewHandler {
	return reviewHandlerImpl{
		reviewService: service,
	}
}

func (ref reviewHandlerImpl) Post(c echo.Context) error {
	var review review.Review
	if err := c.Bind(&review); err != nil {
		log.Err(err).Msg("error binding review")
		if httpErr := err.(*echo.HTTPError); httpErr != nil {
			return c.NoContent(httpErr.Code)
		}
		return c.NoContent(http.StatusBadRequest)
	}

	if err := ref.reviewService.Create(c.Request().Context(), &review); err != nil {
		log.Err(err.Unwrap()).Msg("error at review service")
		return c.NoContent(err.Code)
	}

	return c.NoContent(http.StatusAccepted)
}

func (ref reviewHandlerImpl) Get(c echo.Context) error {
	reviewID := c.Param("id")
	if reviewID == "" {
		log.Warn().Msg("no param provided")
		return c.NoContent(http.StatusBadRequest)
	}

	review, err := ref.reviewService.Get(c.Request().Context(), reviewID)
	if err != nil {
		log.Err(err.Unwrap()).Msg("error at review service")
		return c.NoContent(err.Code)
	}

	return c.JSON(http.StatusOK, review)
}

func (ref reviewHandlerImpl) Delete(c echo.Context) error {
	reviewID := c.Param("id")
	if reviewID == "" {
		log.Warn().Msg("no id param provided")
		return c.NoContent(http.StatusBadRequest)
	}

	if err := ref.reviewService.Delete(c.Request().Context(), reviewID); err != nil {
		log.Err(err.Unwrap()).Msg("error at review service")
		return c.NoContent(err.Code)
	}

	return c.NoContent(http.StatusOK)
}

func (ref reviewHandlerImpl) GetAllBy(c echo.Context) error {
	option := c.QueryParam("option")
	if option == "" {
		log.Warn().Msg("no option param provided")
		return c.NoContent(http.StatusBadRequest)
	}

	id := c.QueryParam("id")
	if id == "" {
		log.Warn().Msg("no id param provided")
		return c.NoContent(http.StatusBadRequest)
	}

	reviews, err := ref.reviewService.GetAllBy(c.Request().Context(), id, option)
	if err != nil {
		log.Err(err.Unwrap()).Msg("error at review service")
		return c.NoContent(err.Code)
	}

	return c.JSON(http.StatusOK, reviews)
}
