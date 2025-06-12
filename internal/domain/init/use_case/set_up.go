package set_up

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/tursodatabase/go-libsql"
)

type setUpUseCase struct{}

func NewSetUpUseCase() *setUpUseCase {
	return &setUpUseCase{}
}

func (uc *setUpUseCase) Execute() {
	db, err := sql.Open("libsql", "file:./shared.db")
	if err != nil {
		panic(err)
	}
	defer func() {
		if closeError := db.Close(); closeError != nil {
			fmt.Println("Error closing database", closeError)
			if err == nil {
				err = closeError
			}
		}
	}()

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to create SQLite driver: %v", err))
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///internal/infra/database/migrations/shared",
		"shared", driver,
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create migration instance: %v", err))
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(fmt.Sprintf("Failed to apply migrations: %v", err))
	}
}
