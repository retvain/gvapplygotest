package operator

import "context"

type Repository interface {
	FindByUuid(ctx context.Context, uuid string) (*Operator, error)
}
