package service

import (
	"Mobile/internal/model/review"
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type reviewServiceImpl struct {
	repository review.ReviewRepository
}

func NewReviewService(repo review.ReviewRepository) ReviewService {
	return reviewServiceImpl{
		repository: repo,
	}
}

func (ref reviewServiceImpl) Create(ctx context.Context, review *review.Review) *echo.HTTPError {

	if review.Title == "" {
		return &echo.HTTPError{Internal: errors.New("title field is missing in new review"), Message: "title field is missing in new review", Code: http.StatusBadRequest}
	}
	if review.Rating == 0 {
		return &echo.HTTPError{Internal: errors.New("rating field is missing in new review"), Message: "rating field is missing in new review", Code: http.StatusBadRequest}
	}
	if review.ProviderId == "" {
		return &echo.HTTPError{Internal: errors.New("provider id field is missing in new review"), Message: "provider id field is missing in new review", Code: http.StatusBadRequest}
	}
	if review.Description == "" {
		return &echo.HTTPError{Internal: errors.New("description field is missing in new review"), Message: "description field is missing in new review", Code: http.StatusBadRequest}
	}
	if review.CustomerId == "" {
		return &echo.HTTPError{Internal: errors.New("customer id field is missing in new review"), Message: "customer id field is missing in new review", Code: http.StatusBadRequest}
	}

	if err := ref.repository.Save(ctx, review); err != nil {
		return err
	}

	return nil
}

func (ref reviewServiceImpl) Delete(ctx context.Context, id string) *echo.HTTPError {
	if err := ref.repository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (ref reviewServiceImpl) Get(ctx context.Context, id string) (*review.Review, *echo.HTTPError) {
	review, err := ref.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (ref reviewServiceImpl) GetAllBy(ctx context.Context, id string, option string) (*[]review.Review, *echo.HTTPError) {
	allReviews, err := ref.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	selectedReviews := []review.Review{}

	if option == "provider" {
		for _, review := range *allReviews {
			if review.ProviderId == id {
				selectedReviews = append(selectedReviews, review)
			}
		}
		return &selectedReviews, nil
	}
	if option == "customer" {
		for _, review := range *allReviews {
			if review.CustomerId == id {
				selectedReviews = append(selectedReviews, review)
			}
		}
		return &selectedReviews, nil
	}

	return nil, &echo.HTTPError{Internal: errors.New("no valid option"), Message: "no valid option", Code: http.StatusBadRequest}
}
