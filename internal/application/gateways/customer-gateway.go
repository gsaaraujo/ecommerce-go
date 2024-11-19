package gateways

import "github.com/google/uuid"

type ICustomerGateway interface {
	ExistsById(customerId uuid.UUID) (bool, error)
}
