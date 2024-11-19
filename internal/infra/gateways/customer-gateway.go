package gateways

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CustomerGateway struct {
	Conn *pgx.Conn
}

func (c *CustomerGateway) ExistsById(id uuid.UUID) (bool, error) {
	var customerId string
	err := c.Conn.QueryRow(context.Background(), "SELECT id FROM customers WHERE id = $1", id).Scan(&customerId)

	if err == nil {
		return true, nil
	}

	if err.Error() == "no rows in result set" {
		return false, nil
	}

	return false, err
}
