package versionApply

import (
	"cmd/internal/versionApply"
	"cmd/pkg/client/postgresql"
	"cmd/pkg/utils/formatter"
	"context"
	"github.com/jackc/pgx/v4"
	"log"
)

type repository struct {
	client postgresql.Client
}

func (r repository) FindNewJob() (queue *versionApply.Queue, err error) {
	q := formatter.Query(`SELECT 
    versions.version_apply_queue.version_id,
    versions.versions.id,
    versions.version_apply_queue.created_at,
    versions.versions.is_applied,
    versions.versions.version_file_path
	FROM versions.version_apply_queue 
	LEFT JOIN versions.versions 
	ON versions.version_apply_queue.version_id=versions.versions.id 
	ORDER BY versions.version_apply_queue.created_at
	FOR UPDATE OF version_apply_queue 
	SKIP LOCKED`)

	applyQueue := versionApply.Queue{}
	var row pgx.Row
	row = r.client.QueryRow(context.TODO(), q)
	err = row.Scan(
		&applyQueue.ID,
		&applyQueue.Version.ID,
		&applyQueue.CreatedAt,
		&applyQueue.Version.IsApplied,
		&applyQueue.Version.VersionFilePath,
	)

	if err != nil && err != pgx.ErrNoRows {
		log.Fatal("error while query")
	}

	return &applyQueue, nil
}

func NewRepository(client postgresql.Client) *repository {
	return &repository{
		client: client,
	}
}
