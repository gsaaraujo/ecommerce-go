package gateways

import "github.com/google/uuid"

type Product struct {
	Id    uuid.UUID
	Price int64
}

type IProductGateway interface {
	FindOneById(id uuid.UUID) (*Product, error)
}
