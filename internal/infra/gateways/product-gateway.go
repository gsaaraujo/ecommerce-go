package gateways

import (
	"context"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-go/internal/application/gateways"
	"github.com/jackc/pgx/v5"
)

type ProductGateway struct {
	Conn *pgx.Conn
}

func (p *ProductGateway) FindOneById(id uuid.UUID) (*gateways.ProductDTO, error) {
	productSchema := struct {
		id    uuid.UUID
		price int64
	}{}

	err := p.Conn.QueryRow(context.Background(), "SELECT id, price FROM products WHERE id = $1", id).
		Scan(&productSchema.id, &productSchema.price)

	if err == nil {
		return &gateways.ProductDTO{
			Id:    productSchema.id,
			Price: productSchema.price,
		}, nil
	}

	if err.Error() == "no rows in result set" {
		return nil, nil
	}

	return nil, err
}
