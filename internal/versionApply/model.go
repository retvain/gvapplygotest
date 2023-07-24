package versionApply

import (
	"cmd/internal/version"
	"time"
)

type Queue struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Version   version.Version
}
