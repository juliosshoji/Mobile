package main

import (
	"Mobile/internal/controller"
	customerhandler "Mobile/internal/controller/customerHandler"
	"Mobile/internal/controller/providerHandler"
	"Mobile/internal/controller/reviewHandler"
	"Mobile/internal/middlewares"
	"Mobile/internal/model/customer"
	"Mobile/internal/model/provider"
	"Mobile/internal/model/review"
	"Mobile/internal/service"
	"context"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func main() {
	// logFileName := os.Getenv("LOG_FILE")
	// if logFileName == "" {
	// 	log.Fatal().Msg("LOG_FILE environment variable not set")
	// 	return
	// }

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal().Msg("HOST environment variable not set")
		return
	}

	if port[0] != ':' {
		port = ":" + port
	}
	log.Info().Msg("port assigned is: " + port)

	// logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	// if err != nil {
	// 	panic(err)
	// }
	// defer logFile.Close()

	// multiWriter := io.MultiWriter(os.Stdout, logFile)

	// dw := diode.NewWriter(multiWriter, 1000, 5*time.Millisecond, func(missed int) {
	// 	fmt.Printf("Global Logger dropped %d messages\n", missed)
	// })

	// log.Logger = zerolog.New(dw).With().Timestamp().Logger()
	// zerolog.TimeFieldFormat = time.RFC3339Nano
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)

	gcproject := os.Getenv("GOOGLE_CLOUD_PROJECT_ID")
	if gcproject == "" {
		// Fallback for local testing (you can set GOOGLE_CLOUD_PROJECT locally)
		gcproject = os.Getenv("GOOGLE_CLOUD_PROJECT")
		if gcproject == "" {
			log.Warn().Msg("GOOGLE_CLOUD_PROJECT_ID or GOOGLE_CLOUD_PROJECT env var not set")
		}
	}
	ctx := context.Background()
	defer ctx.Done()

	client, err := firestore.NewClient(ctx, gcproject)
	if err != nil {
		log.Error().Msg(err.Error() + ", terminating...")
		return
	}

	if err := InitializeServer(client, port); err != nil {
		log.Fatal().Err(err)
		return
	}
}

func InitializeServer(client *firestore.Client, port string) error {
	customerRepository := customer.NewCustomerRepository(client)
	providerRepository := provider.NewProviderRepository(client)
	reviewRepository := review.NewReviewRepository(client)

	customerService := service.NewCustomerService(customerRepository)
	providerService := service.NewProviderService(providerRepository)
	reviewService := service.NewReviewService(reviewRepository)
	authorizationService := service.NewAuthorizationService()

	customerHandler := customerhandler.NewCustomerHandler(customerService)
	providerHandler := providerHandler.NewProviderHandler(providerService)
	reviewHandler := reviewHandler.NewReviewHandler(reviewService)

	echo := echo.New()

	echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	authMiddleware := middlewares.NewUserAuthMiddleware(providerService, customerService, authorizationService)

	apiGroup := echo.Group("/v1")

	DefineRoutes(apiGroup, customerHandler, providerHandler, reviewHandler, authMiddleware)

	return controller.NewServer(echo, port)
}

func DefineRoutes(group *echo.Group, customerHandler customerhandler.CustomerHandler, providerHandler providerHandler.ProviderHandler, reviewHandler reviewHandler.ReviewHandler, authMiddleware middlewares.AuthMiddleware) {
	customersGroup := group.Group("/customers")
	customersGroup.Use(authMiddleware.AuthorizeCustomerMiddleware)
	customersGroup.POST("", customerHandler.Post)
	customersGroup.GET("/:document", customerHandler.Get)
	customersGroup.PUT("/:document", customerHandler.Put)
	customersGroup.DELETE("/:document", customerHandler.Delete)
	customersGroup.PUT("/:document", customerHandler.AddFavorite)

	providersGroup := group.Group("/providers")
	providersGroup.Use(authMiddleware.AuthorizeProviderMiddleware)
	providersGroup.POST("", providerHandler.Post)
	providersGroup.GET("/:document", providerHandler.Get)
	providersGroup.PUT("/:document", providerHandler.Put)
	providersGroup.DELETE("/:document", providerHandler.Delete)
	providersGroup.PUT("/:document", providerHandler.AddSpecialty)
	group.GET("/specialty/:specialty", providerHandler.GetBySpecialty)

	group.POST("/reviews", reviewHandler.Post)
	group.GET("/reviews/:id", reviewHandler.Get)
	group.DELETE("/reviews/:id", reviewHandler.Delete)
	group.GET("/reviews/:option/:id", reviewHandler.GetAllBy)
}
