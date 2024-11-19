package repositories_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-go/internal/domain/models"
	"github.com/gsaaraujo/ecommerce-go/internal/domain/models/cart"
	"github.com/gsaaraujo/ecommerce-go/internal/infra/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type CartRepositorySuite struct {
	conn              *pgx.Conn
	cartRepository    repositories.CartRepository
	postgresContainer testcontainers.Container
	suite.Suite
}

func (c *CartRepositorySuite) SetupTest() {
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

	c.Require().NoError(err)

	host, err := postgresContainer.Host(ctx)
	c.Require().NoError(err)

	port, err := postgresContainer.MappedPort(ctx, "5432")
	c.Require().NoError(err)

	postgresUrl := fmt.Sprintf("postgres://postgres:postgres@%s:%s/postgres", host, port.Port())
	conn, err := pgx.Connect(ctx, postgresUrl)
	c.Require().NoError(err)

	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS products (
			id UUID PRIMARY KEY,
			price INTEGER NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		)
	`)
	c.Require().NoError(err)

	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS carts (
			id UUID PRIMARY KEY,
			customer_id UUID NOT NULL UNIQUE,
			total_price INTEGER NOT NULL,
			total_quantity INTEGER NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		)
  `)
	c.Require().NoError(err)

	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS cart_items (
			id UUID PRIMARY KEY,
			cart_id UUID NOT NULL,
			product_id UUID NOT NULL,
			quantity INTEGER NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (cart_id) REFERENCES carts (id),
			FOREIGN KEY (product_id) REFERENCES products (id)
		)
  `)
	c.Require().NoError(err)

	c.conn = conn
	c.postgresContainer = postgresContainer
	c.cartRepository = repositories.CartRepository{
		Conn: conn,
	}
}

func (p *CartRepositorySuite) TearDownTest() {
	p.postgresContainer.Terminate(context.Background())
}

func (p *CartRepositorySuite) TestCartRepository_Create_OnSuccess_ReturnsNil() {
	ctx := context.Background()
	cartId := uuid.New()
	cartItemId := uuid.New()
	productId := uuid.New()
	customerId := uuid.New()
	cartItem := cart.CartItem{
		Id:        cartItemId,
		ProductId: productId,
		Quantity: models.Quantity{
			Value: 2,
		},
		Price: models.Money{
			Value: 2550,
		},
	}
	cart := cart.Cart{
		Id:         cartId,
		CustomerId: customerId,
		Items:      []cart.CartItem{cartItem},
	}
	_, err := p.conn.Exec(ctx, "INSERT INTO products (id, price) VALUES ($1, $2)", productId, 2550)
	p.Require().NoError(err)

	err = p.cartRepository.Create(cart)
	p.Require().NoError(err)

	cartScheme := struct {
		id            uuid.UUID
		customerId    uuid.UUID
		totalPrice    int64
		totalQuantity int32
		createdAt     time.Time
	}{}

	err = p.conn.QueryRow(ctx, "SELECT id, customer_id, total_price, total_quantity, created_at FROM carts WHERE customer_id = $1", customerId).
		Scan(&cartScheme.id, &cartScheme.customerId, &cartScheme.totalPrice, &cartScheme.totalQuantity, &cartScheme.createdAt)
	p.Require().NoError(err)

	p.Equal(cartId, cartScheme.id)
	p.Equal(customerId, cartScheme.customerId)
	p.Equal(int64(5100), cartScheme.totalPrice)
	p.Equal(int32(2), cartScheme.totalQuantity)

	rows, err := p.conn.Query(ctx, "SELECT id, cart_id, product_id, quantity, created_at FROM cart_items WHERE cart_id = $1", cartScheme.id)
	p.NoError(err)

	for rows.Next() {
		cartItemSchema := struct {
			id        uuid.UUID
			cartId    uuid.UUID
			productId uuid.UUID
			quantity  int32
			createdAt time.Time
		}{}

		err := rows.Scan(&cartItemSchema.id, &cartItemSchema.cartId, &cartItemSchema.productId, &cartItemSchema.quantity, &cartItemSchema.createdAt)
		p.NoError(err)

		p.Equal(cartItemId, cartItemSchema.id)
		p.Equal(cartId, cartItemSchema.cartId)
		p.Equal(productId, cartItemSchema.productId)
	}
}

func (p *CartRepositorySuite) TestCartRepository_Update_OnSuccess_ReturnsNil() {
	ctx := context.Background()
	cartId := uuid.New()
	cartItemId := uuid.New()
	productId := uuid.New()
	customerId := uuid.New()
	cartItem := cart.CartItem{
		Id:        cartItemId,
		ProductId: productId,
		Quantity: models.Quantity{
			Value: 5,
		},
		Price: models.Money{
			Value: 1299,
		},
	}
	cart := cart.Cart{
		Id:         cartId,
		CustomerId: customerId,
		Items:      []cart.CartItem{cartItem},
	}

	_, err := p.conn.Exec(ctx, "INSERT INTO products (id, price) VALUES ($1, $2)", productId, 2550)
	p.Require().NoError(err)
	_, err = p.conn.Exec(ctx, "INSERT INTO carts (id, customer_id, total_price, total_quantity) VALUES ($1, $2, $3, $4)",
		cartId, customerId, 10, 5)
	p.NoError(err)
	_, err = p.conn.Exec(ctx, "INSERT INTO cart_items (id, cart_id, product_id, quantity) VALUES ($1, $2, $3, $4)",
		cartItemId, cartId, productId, 5)
	p.NoError(err)

	err = p.cartRepository.Update(cart)
	p.NoError(err)

	cartScheme := struct {
		id            uuid.UUID
		customerId    uuid.UUID
		totalPrice    int64
		totalQuantity int32
		createdAt     time.Time
	}{}
	err = p.conn.QueryRow(ctx, "SELECT id, customer_id, total_price, total_quantity, created_at FROM carts WHERE customer_id = $1", customerId).
		Scan(&cartScheme.id, &cartScheme.customerId, &cartScheme.totalPrice, &cartScheme.totalQuantity, &cartScheme.createdAt)
	p.NoError(err)

	p.Equal(cartId, cartScheme.id)
	p.Equal(customerId, cartScheme.customerId)
	p.Equal(int64(6495), cartScheme.totalPrice)
	p.Equal(int32(5), cartScheme.totalQuantity)

	rows, err := p.conn.Query(ctx, "SELECT id, cart_id, product_id, quantity, created_at FROM cart_items WHERE cart_id = $1", cartScheme.id)
	p.NoError(err)

	for rows.Next() {
		cartItemSchema := struct {
			id        uuid.UUID
			cartId    uuid.UUID
			productId uuid.UUID
			quantity  int32
			createdAt time.Time
		}{}

		err := rows.Scan(&cartItemSchema.id, &cartItemSchema.cartId, &cartItemSchema.productId, &cartItemSchema.quantity, &cartItemSchema.createdAt)
		p.NoError(err)

		p.Equal(cartItemId, cartItemSchema.id)
		p.Equal(cartId, cartItemSchema.cartId)
		p.Equal(productId, cartItemSchema.productId)
	}
}

func (p *CartRepositorySuite) TestCartRepository_FindOneByCustomerId_OnSuccess_ReturnsCart() {
	ctx := context.Background()
	cartId := uuid.New()
	cartItemId := uuid.New()
	productId := uuid.New()
	customerId := uuid.New()
	cartItem := cart.CartItem{
		Id:        cartItemId,
		ProductId: productId,
		Quantity: models.Quantity{
			Value: 7,
		},
		Price: models.Money{
			Value: 4720,
		},
	}
	cart := cart.Cart{
		Id:         cartId,
		CustomerId: customerId,
		Items:      []cart.CartItem{cartItem},
	}

	_, err := p.conn.Exec(ctx, "INSERT INTO products (id, price) VALUES ($1, $2)", productId, 4720)
	p.NoError(err)
	_, err = p.conn.Exec(ctx, "INSERT INTO carts (id, customer_id, total_price, total_quantity) VALUES ($1, $2, $3, $4)",
		cartId, customerId, 1640, 7)
	p.NoError(err)
	_, err = p.conn.Exec(ctx, "INSERT INTO cart_items (id, cart_id, product_id, quantity) VALUES ($1, $2, $3, $4)",
		cartItemId, cartId, productId, 7)
	p.NoError(err)

	sut, err := p.cartRepository.FindOneByCustomerId(customerId)
	p.NoError(err)

	p.Equal(cart, *sut)
}

func TestCartRepository(t *testing.T) {
	suite.Run(t, new(CartRepositorySuite))
}
