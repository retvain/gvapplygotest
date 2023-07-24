package db

import (
	"cmd/internal/operator"
	"cmd/pkg/client/postgresql"
	"context"
)

type repository struct {
	client postgresql.Client
}

func (r repository) FindOByUuid(ctx context.Context, uuid string) (*operator.Operator, error) {

	q := `SELECT
    operators.operators.id
	FROM operators.operators
	WHERE uuid = $1
`
	o := operator.Operator{Uuid: uuid}
	err := r.client.QueryRow(ctx, q, uuid).Scan(&o.ID)

	if err != nil {
		return &operator.Operator{}, err
	}

	return &o, nil
}

func NewRepository(client postgresql.Client) *repository {
	return &repository{
		client: client,
	}
}
