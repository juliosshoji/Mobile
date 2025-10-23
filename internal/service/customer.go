package service

import (
	"Mobile/internal/model/customer"
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type customerServiceImpl struct {
	repository customer.CustomerRepository
}

func NewCustomerService(repo customer.CustomerRepository) CustomerService {
	return customerServiceImpl{
		repository: repo,
	}
}

func (ref customerServiceImpl) Get(ctx context.Context, document string) (*customer.Customer, *echo.HTTPError) {
	customer, err := ref.repository.Get(ctx, document)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (ref customerServiceImpl) Update(ctx context.Context, customerUpdated *customer.Customer) *echo.HTTPError {

	customer, err := ref.repository.Get(ctx, customerUpdated.Document)
	if err != nil {
		return err
	}

	if customerUpdated.Name != "" {
		customer.Name = customerUpdated.Name
	}
	if customerUpdated.Birthday != "" {
		customer.Birthday = customerUpdated.Birthday
	}

	if err := ref.repository.Update(ctx, customer); err != nil {
		return err
	}
	return nil
}

func (ref customerServiceImpl) Delete(ctx context.Context, document string) *echo.HTTPError {
	if err := ref.repository.Delete(ctx, document); err != nil {
		return err
	}
	return nil
}

func (ref customerServiceImpl) Add(ctx context.Context, customer *customer.Customer) *echo.HTTPError {
	if customer.Birthday == "" {
		return &echo.HTTPError{Internal: errors.New("birthday field is missing in new customer"), Message: "birthday field is missing in new customer", Code: http.StatusBadRequest}
	}
	if customer.Document == "" {
		return &echo.HTTPError{Internal: errors.New("document field is missing in new customer"), Message: "document field is missing in new customer", Code: http.StatusBadRequest}
	}
	if customer.Name == "" {
		return &echo.HTTPError{Internal: errors.New("name field is missing in new customer"), Message: "name field is missing in new customer", Code: http.StatusBadRequest}
	}

	if err := ref.repository.Save(ctx, customer); err != nil {
		return err
	}
	return nil
}

func (ref customerServiceImpl) AddFavorite(ctx context.Context, customerDoc string, providerDoc string) *echo.HTTPError {
	customer, err := ref.repository.Get(ctx, customerDoc)
	if err != nil {
		return err
	}

	if customer.Favorites == nil {
		customer.Favorites = []string{
			providerDoc,
		}
	} else {
		customer.Favorites = append(customer.Favorites, providerDoc)
	}

	if err := ref.repository.Update(ctx, customer); err != nil {
		return err
	}
	return nil
}
