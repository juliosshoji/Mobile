package review

import (
	"Mobile/internal/model"
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ReviewRepository interface {
	model.Repository[Review]
}

type reviewRepositoryImpl struct {
	database *firestore.Client
}

const collection = "reviews"

func NewReviewRepository(client *firestore.Client) ReviewRepository {
	return reviewRepositoryImpl{
		database: client,
	}
}

func (ref reviewRepositoryImpl) Get(ctx context.Context, document string) (*Review, *echo.HTTPError) {
	docSnapshot, err := ref.database.Collection(collection).Doc(document).Get(ctx)
	if err != nil {
		log.Error().Err(err).Msg("could not retrieve document from firestore")
		return nil, model.ErrDocumentNotFound.SetInternal(err)
	}
	if !docSnapshot.Exists() {
		log.Warn().Msg("document does not exist in collection")
		return nil, model.ErrDocumentNotExists
	}
	var review Review
	if err := docSnapshot.DataTo(&review); err != nil {
		return nil, &echo.HTTPError{Internal: err, Message: "error binding info to review", Code: http.StatusInternalServerError}
	}

	return &review, nil
}
func (ref reviewRepositoryImpl) Save(ctx context.Context, review *Review) *echo.HTTPError {
	docRef := ref.database.Collection(collection).NewDoc()

	review.Id = docRef.ID

	if _, err := docRef.Set(ctx, review); err != nil {
		log.Error().Err(err).Msg("error creating document in firestore")
		return &echo.HTTPError{Internal: err, Message: "error creating document", Code: http.StatusInternalServerError}
	}

	log.Info().Msg("new review created")

	return nil
}

func (ref reviewRepositoryImpl) Update(ctx context.Context, review *Review) *echo.HTTPError {
	return nil
}

func (ref reviewRepositoryImpl) Delete(ctx context.Context, document string) *echo.HTTPError {
	docRef := ref.database.Collection(collection).Doc(document)
	if _, err := docRef.Delete(ctx); err != nil {
		log.Error().Err(err).Msg("error deleting document in firestore")
		return &echo.HTTPError{Internal: err, Message: "error deleting document", Code: http.StatusInternalServerError}
	}

	return nil
}

func (ref reviewRepositoryImpl) GetAll(ctx context.Context) (*[]Review, *echo.HTTPError) {
	docIterator := ref.database.Collection(collection).DocumentRefs(ctx)

	docRefs, err := docIterator.GetAll()
	if err != nil {
		log.Error().Err(err).Msg("error getting documents from firestore")
		return nil, &echo.HTTPError{Internal: err, Message: "error getting documents from firestore", Code: http.StatusInternalServerError}
	}

	allReviews := []Review{}
	review := Review{}
	for _, doc := range docRefs {
		snapshot, err := doc.Get(ctx)
		if err != nil {
			log.Error().Err(err).Msg("error getting snapshot from firestore")
			return nil, &echo.HTTPError{Internal: err, Message: "error getting snapshot from firestore", Code: http.StatusInternalServerError}
		}

		if err := snapshot.DataTo(&review); err != nil {
			log.Error().Err(err).Msg("error binding data of provider from firestore")
			return nil, &echo.HTTPError{Internal: err, Message: "error binding data of provider from firestore", Code: http.StatusInternalServerError}
		}

		allReviews = append(allReviews, review)
	}

	return &allReviews, nil
}
