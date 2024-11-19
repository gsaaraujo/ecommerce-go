package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-go/internal/domain/models"
	"github.com/gsaaraujo/ecommerce-go/internal/domain/models/cart"
	"github.com/jackc/pgx/v5"
)

type CartRepository struct {
	Conn *pgx.Conn
}

func (c *CartRepository) Create(cart cart.Cart) error {
	ctx := context.Background()

	transaction, err := c.Conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer transaction.Rollback(ctx)

	_, err = transaction.Exec(ctx, "INSERT INTO carts (id, customer_id, total_price, total_quantity) VALUES ($1, $2, $3, $4)",
		cart.Id.String(), cart.CustomerId.String(), cart.TotalPrice().Value, cart.TotalQuantity().Value)

	if err != nil {
		return err
	}

	for _, cartItem := range cart.Items {
		_, err = transaction.Exec(ctx, "INSERT INTO cart_items (id, cart_id, product_id, quantity) VALUES ($1, $2, $3, $4)",
			cartItem.Id.String(), cart.Id.String(), cartItem.ProductId.String(), cart.TotalQuantity().Value)

		if err != nil {
			return err
		}
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *CartRepository) Update(cart cart.Cart) error {
	ctx := context.Background()

	transaction, err := c.Conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer transaction.Rollback(context.Background())

	_, err = transaction.Exec(ctx, "UPDATE carts SET total_price = $1, total_quantity = $2 WHERE id = $3",
		cart.TotalPrice().Value, cart.TotalQuantity().Value, cart.Id.String())

	if err != nil {
		return err
	}

	_, err = transaction.Exec(ctx, "DELETE FROM cart_items WHERE cart_id = $1", cart.Id.String())

	if err != nil {
		return err
	}

	for _, cartItem := range cart.Items {
		_, err = transaction.Exec(ctx, "INSERT INTO cart_items (id, cart_id, product_id, quantity) VALUES ($1, $2, $3, $4)",
			cartItem.Id.String(), cart.Id.String(), cartItem.ProductId.String(), cart.TotalQuantity().Value)

		if err != nil {
			return err
		}
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *CartRepository) FindOneByCustomerId(customerId uuid.UUID) (*cart.Cart, error) {
	ctx := context.Background()

	type CartSchema struct {
		id            uuid.UUID
		customerId    uuid.UUID
		totalPrice    int64
		totalQuantity int32
		createdAt     time.Time
	}

	type CartItemSchema struct {
		id        uuid.UUID
		cartId    uuid.UUID
		productId uuid.UUID
		quantity  int32
		price     int64
		createdAt time.Time
	}

	var cartSchema CartSchema
	err := c.Conn.QueryRow(ctx,
		"SELECT id, customer_id, total_price, total_quantity, created_at FROM carts WHERE customer_id = $1", customerId).
		Scan(&cartSchema.id, &cartSchema.customerId, &cartSchema.totalPrice, &cartSchema.totalQuantity, &cartSchema.createdAt)

	if err != nil {
		return nil, err
	}

	rows, err := c.Conn.Query(ctx,
		`SELECT 
					ci.id,
					ci.cart_id,
					ci.product_id,
					ci.quantity,
					ci.created_at,
					p.price
			 FROM cart_items ci
			 JOIN products p ON ci.product_id = p.id
			 WHERE ci.cart_id = $1`, cartSchema.id)

	if err != nil {
		return nil, err
	}

	var cartItemsSchema []CartItemSchema
	for rows.Next() {
		var cartItemSchema CartItemSchema
		err := rows.Scan(&cartItemSchema.id, &cartItemSchema.cartId, &cartItemSchema.productId,
			&cartItemSchema.quantity, &cartItemSchema.createdAt, &cartItemSchema.price)

		if err != nil {
			return nil, err
		}

		cartItemsSchema = append(cartItemsSchema, cartItemSchema)
	}

	cartItems := []cart.CartItem{}
	for _, cartItemSchema := range cartItemsSchema {
		cartItem := cart.CartItem{
			Id:        cartItemSchema.id,
			ProductId: cartItemSchema.productId,
			Quantity: models.Quantity{
				Value: cartItemSchema.quantity,
			},
			Price: models.Money{
				Value: cartItemSchema.price,
			},
		}
		cartItems = append(cartItems, cartItem)
	}

	cart := cart.Cart{
		Id:         cartSchema.id,
		CustomerId: cartSchema.customerId,
		Items:      cartItems,
	}

	return &cart, nil
}
