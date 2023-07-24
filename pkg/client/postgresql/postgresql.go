package postgresql

import (
	pgConf "cmd/config/database"
	"cmd/pkg/utils"
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient() (pool *pgxpool.Pool, err error) {
	config := pgConf.Get()
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	var timeout int
	timeout = config.Timeout

	err = repeatable.DoWithTries(func() error {
		_, cancel := context.WithTimeout(context.TODO(), time.Duration(timeout))
		defer cancel()

		pool, err = pgxpool.Connect(context.TODO(), dsn)
		if err != nil {
			panic(err)
		}
		return nil
	},
		config.MaxAttempts,
		time.Duration(timeout),
	)

	if err != nil {
		log.Fatalf("error do with tries postgresql")
	}

	return pool, nil
}
