package gateways

import "github.com/google/uuid"

type ProductDTO struct {
	Id    uuid.UUID
	Price int64
}

type IProductGateway interface {
	FindOneById(id uuid.UUID) (*ProductDTO, error)
}
