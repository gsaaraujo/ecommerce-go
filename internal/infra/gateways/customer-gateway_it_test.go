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

type CustomerGatewaySuite struct {
	conn              *pgx.Conn
	customerGateway   gateways.CustomerGateway
	postgresContainer testcontainers.Container
	suite.Suite
}

func (p *CustomerGatewaySuite) SetupTest() {
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
		CREATE TABLE IF NOT EXISTS customers (
			id UUID PRIMARY KEY,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		)
  `)
	p.Require().NoError(err)

	p.conn = conn
	p.postgresContainer = postgresContainer
	p.customerGateway = gateways.CustomerGateway{
		Conn: conn,
	}
}

func (p *CustomerGatewaySuite) TearDownTest() {
	p.postgresContainer.Terminate(context.Background())
}

func (p *CustomerGatewaySuite) TestCustomerGateway_ExistsById_OnCustomerExists_ReturnsCustomer() {
	customerId := uuid.New()
	_, err := p.conn.Exec(context.Background(), "INSERT INTO customers (id) VALUES ($1)", customerId)
	p.Require().NoError(err)

	exists, err := p.customerGateway.ExistsById(customerId)

	p.NoError(err)
	p.Equal(exists, true)
}

func (suite *CustomerGatewaySuite) TestCustomerGateway_ExistsById_OnCustomerNotExists_ReturnsCustomer() {
	customerId := uuid.New()

	exists, err := suite.customerGateway.ExistsById(customerId)

	suite.NoError(err)
	suite.Equal(exists, false)
}

func TestCustomerGateway(t *testing.T) {
	suite.Run(t, new(CustomerGatewaySuite))
}
