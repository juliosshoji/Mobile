package service

import (
	"Mobile/internal/model/customer"
	"Mobile/internal/model/domain"
	"Mobile/internal/model/provider"
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const expirationMin = 30

var jwtKey = []byte("bankingapi-key")

type authenticationServiceImpl struct {
	providerRepo provider.ProviderRepository
	customerRepo customer.CustomerRepository
}

func NewAuthenticationService(providerRepo provider.ProviderRepository, customerRepo customer.CustomerRepository) AuthenticationService {
	return authenticationServiceImpl{
		providerRepo: providerRepo,
		customerRepo: customerRepo,
	}
}

func (ref authenticationServiceImpl) GenerateToken(ctx context.Context, document string) (*string, *echo.HTTPError) {
	expirationTime := time.Now().Add(time.Minute * expirationMin)
	Claim := &domain.Claims{
		Id: (document),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: expirationTime},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, *Claim)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}
	log.Info().Msg("Authenticated entrance: " + document)

	return &tokenString, nil
}

func (ref authenticationServiceImpl) AuthenticateCustomer(ctx context.Context, document, password string) (*string, *echo.HTTPError) {
	customerData, err := ref.customerRepo.Get(ctx, document)
	if err != nil {
		return nil, err
	}
	if password != customerData.Password {
		return nil, &echo.HTTPError{Message: "invalid credentials", Code: 401}
	}

	return ref.GenerateToken(ctx, document)
}

func (ref authenticationServiceImpl) AuthenticateProvider(ctx context.Context, document, password string) (*string, *echo.HTTPError) {
	providerData, err := ref.providerRepo.Get(ctx, document)
	if err != nil {
		return nil, err
	}

	encodedPassword := sha256.Sum256([]byte(providerData.Password))
	if password != fmt.Sprintf("%x", encodedPassword) {
		return nil, &echo.HTTPError{Message: "invalid credentials", Code: 401}
	}

	return ref.GenerateToken(ctx, document)
}
