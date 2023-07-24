package processor

import (
	"cmd/internal/version"
	"cmd/internal/version/actual/action"
	"cmd/internal/version/actual/dto"
	"cmd/internal/versionApply"
	versionApplyRepository "cmd/internal/versionApply/db"
	"cmd/pkg/client/postgresql"
	"cmd/pkg/parser"
	"cmd/pkg/utils/file"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type App struct {
	postgreSQLClient *pgxpool.Pool
}

func Run() (err error) {
	var app App
	app.postgreSQLClient, err = postgresql.NewClient()
	if err != nil {
		return err
	}

	var newJob *versionApply.Queue
	newJob = FindNewJob(&app)

	// если найдена новая запись в очереди
	if newJob.ID != 0 {
		startTime := time.Now()
		var container *dto.Container
		container, err = ParseReferenceActual(&newJob.Version)
		if err != nil {
			return err
		}
		err = ApplyVersion(container)
		if err != nil {
			return err
		}
		endTime := time.Now()
		fmt.Printf("Применение версии выполнено за %v\n", endTime.Sub(startTime))
	}
	return nil
}

func FindNewJob(app *App) *versionApply.Queue {
	repository := versionApplyRepository.NewRepository(app.postgreSQLClient)

	job, err := repository.FindNewJob()
	if err != nil {
		log.Fatalf("error create repository")
	}

	return job
}

func ParseReferenceActual(version *version.Version) (*dto.Container, error) {
	var path *string
	path = &version.VersionFilePath
	// если такой файл найден
	if file.Exists(path) {
		container, err := parser.ParseSax(path)
		return container, err
	} else {
		err := errors.New(fmt.Sprintf("ошибка применения версии, файл <%f> не найден", path))
		return nil, err
	}
}

func ApplyVersion(container *dto.Container) (err error) {
	conn, err := postgresql.NewClient()
	if err != nil {
		return err
	}

	// запускаем транзакцию
	tx, err := conn.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(context.TODO())
			if err != nil {
			}
		} else {
			err = tx.Commit(context.TODO())
			if err != nil {
			}
		}
	}()

	err = action.ApplyVersion(conn, container)
	if err != nil {
		return err
	}

	return nil
}
