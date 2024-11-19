package gateways_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-go/internal/infra/gateways"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ProductGatewaySuite struct {
	conn              *pgx.Conn
	productGateway    gateways.ProductGateway
	postgresContainer testcontainers.Container
	suite.Suite
}

func (p *ProductGatewaySuite) SetupTest() {
	ctx := context.Background()
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:latest",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
				"POSTGRES_DB":       "postgres",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp"),
		},
	})

	p.Require().NoError(err)

	host, err := postgresContainer.Host(ctx)
	p.Require().NoError(err)

	port, err := postgresContainer.MappedPort(ctx, "5432")
	p.Require().NoError(err)

	postgresUrl := fmt.Sprintf("postgres://postgres:postgres@%s:%s/postgres", host, port.Port())
	conn, err := pgx.Connect(ctx, postgresUrl)
	p.Require().NoError(err)

	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS products (
			id UUID PRIMARY KEY,
			price INTEGER NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		)
  `)
	p.Require().NoError(err)

	p.conn = conn
	p.postgresContainer = postgresContainer
	p.productGateway = gateways.ProductGateway{
		Conn: conn,
	}
}

func (p *ProductGatewaySuite) TearDownTest() {
	p.postgresContainer.Terminate(context.Background())
}

func (p *ProductGatewaySuite) TestProductGateway_FindOneById_OnProductExists_ReturnsProduct() {
	productId := uuid.New()
	_, err := p.conn.Exec(context.Background(), "INSERT INTO products (id, price) VALUES ($1, $2)", productId, 2550)
	p.Require().NoError(err)

	sut, err := p.productGateway.FindOneById(productId)

	p.NoError(err)
	p.Equal(productId, sut.Id)
	p.Equal(int64(2550), sut.Price)
}

func (p *ProductGatewaySuite) TestProductGateway_FindOneById_OnProductNotExists_ReturnsNil() {
	productId := uuid.New()

	sut, err := p.productGateway.FindOneById(productId)

	p.NoError(err)
	p.Nil(sut)
}

func TestProductGateway(t *testing.T) {
	suite.Run(t, new(ProductGatewaySuite))
}
