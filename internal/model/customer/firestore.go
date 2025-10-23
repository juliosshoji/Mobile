package customer

import (
	"Mobile/internal/model"
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type CustomerRepository interface {
	model.Repository[Customer]
}

type customerRepositoryImpl struct {
	databaseClient *firestore.Client
}

const collection = "customers"

func NewCustomerRepository(client *firestore.Client) CustomerRepository {
	return customerRepositoryImpl{
		databaseClient: client,
	}
}

func (ref customerRepositoryImpl) Get(ctx context.Context, document string) (*Customer, *echo.HTTPError) {
	docSnapshot, err := ref.databaseClient.Collection(collection).Doc(document).Get(ctx)
	if err != nil {
		log.Error().Err(err).Msg("could not retrieve document from firestore")
		return nil, model.ErrDocumentNotFound.SetInternal(err)
	}
	if !docSnapshot.Exists() {
		log.Warn().Msg("document does not exist in collection")
		return nil, model.ErrDocumentNotExists
	}
	var customer Customer
	if err := docSnapshot.DataTo(&customer); err != nil {
		return nil, &echo.HTTPError{Internal: err, Message: "error binding info to customer", Code: http.StatusInternalServerError}
	}

	return &customer, nil
}

func (ref customerRepositoryImpl) Save(ctx context.Context, customer *Customer) *echo.HTTPError {
	docRef := ref.databaseClient.Collection(collection).Doc(customer.Document)
	if _, err := docRef.Create(ctx, customer); err != nil {
		log.Error().Err(err).Msg("error creating document in firestore")
		return &echo.HTTPError{Internal: err, Message: "error creating document", Code: http.StatusInternalServerError}
	}

	return nil
}

func (ref customerRepositoryImpl) Update(ctx context.Context, customer *Customer) *echo.HTTPError {
	docRef := ref.databaseClient.Collection(collection).Doc(customer.Document)
	if _, err := docRef.Set(ctx, customer); err != nil {
		log.Error().Err(err).Msg("error updating document in firestore")
		return &echo.HTTPError{Internal: err, Message: "error updating document", Code: http.StatusInternalServerError}
	}

	return nil
}
func (ref customerRepositoryImpl) Delete(ctx context.Context, document string) *echo.HTTPError {
	docRef := ref.databaseClient.Collection(collection).Doc(document)
	if _, err := docRef.Delete(ctx); err != nil {
		log.Error().Err(err).Msg("error deleting document in firestore")
		return &echo.HTTPError{Internal: err, Message: "error deleting document", Code: http.StatusInternalServerError}
	}

	return nil
}

func (ref customerRepositoryImpl) GetAll(ctx context.Context) (*[]Customer, *echo.HTTPError) {
	docIterator := ref.databaseClient.Collection(collection).DocumentRefs(ctx)

	docRefs, err := docIterator.GetAll()
	if err != nil {
		log.Error().Err(err).Msg("error getting documents from firestore")
		return nil, &echo.HTTPError{Internal: err, Message: "error getting documents from firestore", Code: http.StatusInternalServerError}
	}

	allCustomers := []Customer{}
	customer := Customer{}
	for _, doc := range docRefs {
		snapshot, err := doc.Get(ctx)
		if err != nil {
			log.Error().Err(err).Msg("error getting snapshot from firestore")
			return nil, &echo.HTTPError{Internal: err, Message: "error getting snapshot from firestore", Code: http.StatusInternalServerError}
		}

		if err := snapshot.DataTo(&customer); err != nil {
			log.Error().Err(err).Msg("error binding data of customer from firestore")
			return nil, &echo.HTTPError{Internal: err, Message: "error binding data of customer from firestore", Code: http.StatusInternalServerError}
		}

		allCustomers = append(allCustomers, customer)
	}

	return &allCustomers, nil
}
