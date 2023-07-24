package version

import "context"

type Repository interface {
	FindAll(ctx context.Context) (u []Version, err error)
	FindOne(ctx context.Context, id string) (*Version, error)
	//Update(ctx context.Context, version Version) error
}
