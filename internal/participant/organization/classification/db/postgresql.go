package db

import (
	"cmd/internal/participant/organization/classification"
	"cmd/pkg/client/postgresql"
	"context"
)

type repository struct {
	client postgresql.Client
}

func (r repository) FindByClassificationId(
	ctx context.Context,
	classificationId string,
) (*classification.Classification, error) {

	q := `SELECT
    participants.classifications.id
	FROM participants.classifications
	WHERE classification_id = $1
`
	c := classification.Classification{ClassificationId: classificationId}
	err := r.client.QueryRow(ctx, q, classificationId).Scan(&c.ID)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r repository) GetAll(
	ctx context.Context,
) (*[]classification.Classification, error) {

	q := `SELECT
    participants.classifications.id,
    participants.classifications.classification_id
	FROM participants.classifications
`
	var c classification.Classification
	collect := make([]classification.Classification, 0)
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&c.ID, &c.ClassificationId)
		if err != nil {
			return nil, err
		}
		collect = append(collect, c)
	}

	return &collect, nil
}

func NewRepository(client postgresql.Client) *repository {
	return &repository{
		client: client,
	}
}
