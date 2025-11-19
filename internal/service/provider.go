package service

import (
	"Mobile/internal/model/provider"
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type providerServiceImpl struct {
	repository provider.ProviderRepository
}

func NewProviderService(repo provider.ProviderRepository) ProviderService {
	return providerServiceImpl{
		repository: repo,
	}
}

func (ref providerServiceImpl) Get(ctx context.Context, document string) (*provider.Provider, *echo.HTTPError) {
	provider, err := ref.repository.Get(ctx, document)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (ref providerServiceImpl) Update(ctx context.Context, providerUpdated *provider.Provider) *echo.HTTPError {
	provider, err := ref.repository.Get(ctx, providerUpdated.Document)
	if err != nil {
		return err
	}

	if providerUpdated.Birthday != "" {
		provider.Birthday = providerUpdated.Birthday
	}
	if providerUpdated.Name != "" {
		provider.Name = providerUpdated.Name
	}
	if providerUpdated.ContactType != "" {
		if providerUpdated.ContactAddress == "" {
			return &echo.HTTPError{Internal: errors.New("contact address did not change with type"), Message: "contact address did not change with type", Code: http.StatusBadRequest}
		}
		if !providerUpdated.ContactType.IsValid() {
			return &echo.HTTPError{Internal: errors.New("contact type is invalid"), Message: "contact type is invalid", Code: http.StatusBadRequest}
		}
		provider.ContactType = providerUpdated.ContactType
		provider.ContactAddress = providerUpdated.ContactAddress

	}

	if err := ref.repository.Update(ctx, provider); err != nil {
		return err
	}

	return nil
}

func (ref providerServiceImpl) Delete(ctx context.Context, document string) *echo.HTTPError {
	if err := ref.repository.Delete(ctx, document); err != nil {
		return err
	}
	return nil
}

func (ref providerServiceImpl) Add(ctx context.Context, provider *provider.Provider) *echo.HTTPError {
	if provider.Birthday == "" {
		return &echo.HTTPError{Internal: errors.New("birthday field is missing in new provider"), Message: "birthday field is missing in new provider", Code: http.StatusBadRequest}
	}
	if provider.ContactAddress == "" {
		return &echo.HTTPError{Internal: errors.New("contact field is missing in new provider"), Message: "contact field is missing in new provider", Code: http.StatusBadRequest}
	}
	if !provider.ContactType.IsValid() {
		return &echo.HTTPError{Internal: errors.New("contact type field is missing or is invalid in new provider"), Message: "contact field is missing or is invalid in new provider", Code: http.StatusBadRequest}
	}
	if provider.Document == "" {
		return &echo.HTTPError{Internal: errors.New("document field is missing in new provider"), Message: "document field is missing in new provider", Code: http.StatusBadRequest}
	}
	if provider.Name == "" {
		return &echo.HTTPError{Internal: errors.New("name field is missing in new provider"), Message: "name field is missing in new provider", Code: http.StatusBadRequest}
	}
	if provider.Specialties == nil {
		return &echo.HTTPError{Internal: errors.New("specialties field is missing in new provider"), Message: "specialties field is missing in new provider", Code: http.StatusBadRequest}
	}

	if err := ref.repository.Save(ctx, provider); err != nil {
		return err
	}

	return nil
}

func (ref providerServiceImpl) AddSpecialty(ctx context.Context, specialty *provider.Specialty, providerDocument string) *echo.HTTPError {
	providerUpdate, err := ref.repository.Get(ctx, providerDocument)
	if err != nil {
		return err
	}

	if *specialty == "" {
		return &echo.HTTPError{Internal: errors.New("specialty was not provided"), Message: "specialty was not provided", Code: http.StatusBadRequest}
	}

	if providerUpdate.Specialties == nil {
		providerUpdate.Specialties = []provider.Specialty{
			*specialty,
		}
	} else {
		providerUpdate.Specialties = append(providerUpdate.Specialties, *specialty)
	}

	if err := ref.repository.Update(ctx, providerUpdate); err != nil {
		return err
	}

	return nil
}

func (ref providerServiceImpl) GetBySpecialty(ctx context.Context, specialty *provider.Specialty) (*[]provider.Provider, *echo.HTTPError) {

	unfilteredProviders, err := ref.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	filteredProviders := []provider.Provider{}

	for _, provider := range *unfilteredProviders {
		for _, providerSpecialty := range provider.Specialties {
			if providerSpecialty == *specialty {
				filteredProviders = append(filteredProviders, provider)
				break
			}
		}
	}

	return &filteredProviders, nil
}

func (ref providerServiceImpl) AddProfilePhoto(ctx context.Context, providerDocument string, photoData string) *echo.HTTPError {
	provider, err := ref.repository.Get(ctx, providerDocument)
	if err != nil {
		return nil
	}
	provider.ProfilePhoto = photoData

	if err := ref.repository.Update(ctx, provider); err != nil {
		return nil
	}

	return nil
}
