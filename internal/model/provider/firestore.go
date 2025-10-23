package provider

import (
	"Mobile/internal/model"
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ProviderRepository interface {
	model.Repository[Provider]
}

type providerRepositoryImpl struct {
	database *firestore.Client
}

const collection = "providers"

func NewProviderRepository(client *firestore.Client) ProviderRepository {
	return providerRepositoryImpl{
		database: client,
	}
}

func (ref providerRepositoryImpl) Get(ctx context.Context, document string) (*Provider, *echo.HTTPError) {
	docSnapshot, err := ref.database.Collection(collection).Doc(document).Get(ctx)
	if err != nil {
		log.Error().Err(err).Msg("could not retrieve document from firestore")
		return nil, model.ErrDocumentNotFound.SetInternal(err)
	}
	if !docSnapshot.Exists() {
		log.Warn().Msg("document does not exist in collection")
		return nil, model.ErrDocumentNotExists
	}
	var provider Provider
	if err := docSnapshot.DataTo(&provider); err != nil {
		return nil, &echo.HTTPError{Internal: err, Message: "error binding info to provider", Code: http.StatusInternalServerError}
	}

	return &provider, nil
}

func (ref providerRepositoryImpl) Save(ctx context.Context, provider *Provider) *echo.HTTPError {
	docRef := ref.database.Collection(collection).Doc(provider.Document)
	if _, err := docRef.Create(ctx, provider); err != nil {
		log.Error().Err(err).Msg("error creating document in firestore")
		return &echo.HTTPError{Internal: err, Message: "error creating document", Code: http.StatusInternalServerError}
	}

	return nil
}

func (ref providerRepositoryImpl) Update(ctx context.Context, provider *Provider) *echo.HTTPError {
	docRef := ref.database.Collection(collection).Doc(provider.Document)
	if _, err := docRef.Set(ctx, provider); err != nil {
		log.Error().Err(err).Msg("error updating document in firestore")
		return &echo.HTTPError{Internal: err, Message: "error updating document", Code: http.StatusInternalServerError}
	}

	return nil
}

func (ref providerRepositoryImpl) Delete(ctx context.Context, document string) *echo.HTTPError {
	docRef := ref.database.Collection(collection).Doc(document)
	if _, err := docRef.Delete(ctx); err != nil {
		log.Error().Err(err).Msg("error deleting document in firestore")
		return &echo.HTTPError{Internal: err, Message: "error deleting document", Code: http.StatusInternalServerError}
	}

	return nil
}
func (ref providerRepositoryImpl) GetAll(ctx context.Context) (*[]Provider, *echo.HTTPError) {
	docIterator := ref.database.Collection(collection).DocumentRefs(ctx)

	docRefs, err := docIterator.GetAll()
	if err != nil {
		log.Error().Err(err).Msg("error getting documents from firestore")
		return nil, &echo.HTTPError{Internal: err, Message: "error getting documents from firestore", Code: http.StatusInternalServerError}
	}

	allProviders := []Provider{}
	provider := Provider{}
	for _, doc := range docRefs {
		snapshot, err := doc.Get(ctx)
		if err != nil {
			log.Error().Err(err).Msg("error getting snapshot from firestore")
			return nil, &echo.HTTPError{Internal: err, Message: "error getting snapshot from firestore", Code: http.StatusInternalServerError}
		}

		if err := snapshot.DataTo(&provider); err != nil {
			log.Error().Err(err).Msg("error binding data of provider from firestore")
			return nil, &echo.HTTPError{Internal: err, Message: "error binding data of provider from firestore", Code: http.StatusInternalServerError}
		}

		allProviders = append(allProviders, provider)
	}

	return &allProviders, nil
}
