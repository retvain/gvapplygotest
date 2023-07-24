package classification

import "context"

type Repository interface {
	FindByClassificationId(ctx context.Context, classificationId string) (*Classification, error)
	GetAll(ctx context.Context) (*[]Classification, error)
}
