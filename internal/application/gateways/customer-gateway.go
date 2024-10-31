package gateways

import "github.com/google/uuid"

type ICustomerGateway interface {
	ExistsByCustomerId(customerId uuid.UUID) (bool, error)
}
