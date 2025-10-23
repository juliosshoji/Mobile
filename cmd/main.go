package main

import (
	"Mobile/internal/controller"
	customerhandler "Mobile/internal/controller/customerHandler"
	"Mobile/internal/controller/providerHandler"
	"Mobile/internal/controller/reviewHandler"
	"Mobile/internal/model/customer"
	"Mobile/internal/model/provider"
	"Mobile/internal/model/review"
	"Mobile/internal/service"
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	logFileName := os.Getenv("LOG_FILE")
	if logFileName == "" {
		log.Fatal().Msg("LOG_FILE environment variable not set")
		return
	}

	port := os.Getenv("HOST")
	if port == "" {
		log.Fatal().Msg("HOST environment variable not set")
		return
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

	gcproject := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if gcproject == "" {
		log.Fatal().Msg("GOOGLE_CLOUD_PROJECT environment variable not set")
		return
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

	customerHandler := customerhandler.NewCustomerHandler(customerService)
	providerHandler := providerHandler.NewProviderHandler(providerService)
	reviewHandler := reviewHandler.NewReviewHandler(reviewService)

	echo := echo.New()

	apiGroup := echo.Group("/v1")

	DefineRoutes(apiGroup, customerHandler, providerHandler, reviewHandler)

	return controller.NewServer(echo, port)
}

func DefineRoutes(group *echo.Group, customerHandler customerhandler.CustomerHandler, providerHandler providerHandler.ProviderHandler, reviewHandler reviewHandler.ReviewHandler) {
	group.POST("/customers", customerHandler.Post)
	group.GET("/customers/:document", customerHandler.Get)
	group.PUT("/customers/:document", customerHandler.Put)
	group.DELETE("/customers/:document", customerHandler.Delete)
	group.PUT("/customers/:customer_id", customerHandler.AddFavorite)

	group.POST("/providers", providerHandler.Post)
	group.GET("/providers/:document", providerHandler.Get)
	group.PUT("/providers/:document", providerHandler.Put)
	group.DELETE("/providers/:document", providerHandler.Delete)
	group.PUT("/providers/specialty/:document", providerHandler.AddSpecialty)
	group.GET("/providers/specialty/:specialty", providerHandler.GetBySpecialty)

	group.POST("/reviews", reviewHandler.Post)
	group.GET("/reviews/:id", reviewHandler.Get)
	group.DELETE("/reviews/:id", reviewHandler.Delete)
	group.GET("/reviews/:option/:id", reviewHandler.GetAllBy)
}
