package version

import (
	"cmd/internal/version"
	"cmd/pkg/client/postgresql"
	"context"
)

type repository struct {
	client postgresql.Client
}

func (r repository) FindOne(ctx context.Context, id string) (version.Version, error) {

	q := `SELECT
    versions.versions.id,
	versions.versions.version_file_path,
	versions.versions.is_applied
	FROM versions.versions	
	WHERE id = $1
`
	var v version.Version

	err := r.client.QueryRow(ctx, q, id).Scan(
		&v.ID,
		&v.VersionFilePath,
		&v.IsApplied)

	if err != nil {
		return version.Version{}, err
	}

	return v, nil
}

func (r repository) FindAll(ctx context.Context) (u []version.Version, err error) {
	q := "SELECT id, version_file_path FROM versions.versions"

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	versions := make([]version.Version, 0)

	for rows.Next() {
		var v version.Version

		err = rows.Scan(&v.ID, &v.VersionFilePath)

		if err != nil {
			return nil, err
		}

		versions = append(versions, v)

		if err = rows.Err(); err != nil {
			return nil, err
		}

	}

	return versions, nil
}

func NewRepository(client postgresql.Client) *repository {
	return &repository{
		client: client,
	}
}
