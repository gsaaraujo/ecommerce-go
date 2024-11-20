package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gsaaraujo/ecommerce-go/internal/application/usecases"
	"github.com/gsaaraujo/ecommerce-go/internal/infra"
	"github.com/gsaaraujo/ecommerce-go/internal/infra/gateways"
	"github.com/gsaaraujo/ecommerce-go/internal/infra/handlers"
	"github.com/gsaaraujo/ecommerce-go/internal/infra/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func main() {
	ctx := context.Background()

	if _, ok := os.LookupEnv("AWS_REGION"); !ok {
		panic("environment variable 'AWS_REGION' not set")
	}

	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		panic(err)
	}

	secretManager := secretsmanager.NewFromConfig(awsConfig)
	awsSecretManagerGateway := gateways.AwsSecretManagerGateway{
		SecretManager: secretManager,
	}

	dbUrl, err := awsSecretManagerGateway.Get("DATABASE_URL")
	if err != nil {
		panic(err)
	}

	dbConn, err := pgx.Connect(ctx, os.Getenv(dbUrl))
	if err != nil {
		panic(err)
	}

	validator := infra.NewValidator()

	cartRepository := repositories.CartRepository{
		Conn: dbConn,
	}

	customerGateway := gateways.CustomerGateway{
		Conn: dbConn,
	}

	productGateway := gateways.ProductGateway{
		Conn: dbConn,
	}

	addProductToCart := usecases.AddProductToCart{
		CustomerGateway: &customerGateway,
		ProductGateway:  &productGateway,
		CartRepository:  &cartRepository,
	}

	addProductToCartHandler := handlers.SecurityHandlerDecorator{
		SecretManagerGateway: &awsSecretManagerGateway,
		HttpHandler: &handlers.AddProductToCartHandler{
			Validator:        validator,
			AddProductToCart: &addProductToCart,
		},
	}

	e := echo.New()

	e.GET("/add-product-to-cart", func(c echo.Context) error {
		return addProductToCartHandler.Handle(c)
	})
}
