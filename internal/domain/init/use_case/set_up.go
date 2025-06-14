package set_up

import (
	"database/sql"
	"fmt"
	"io/fs"
	"strings"

	_ "github.com/tursodatabase/go-libsql"
)

const FILE_UP_MIGRATION_INDEX = 1

type setUpUseCase struct {
	Migrations fs.FS
}

func NewSetUpUseCase(migrations fs.FS) *setUpUseCase {
	return &setUpUseCase{
		Migrations: migrations,
	}
}

func (uc *setUpUseCase) Execute() {
	db, err := sql.Open("libsql", "file:./shared.db")
	if err != nil {
		panic(fmt.Sprintf("Failed to open database: %v", err))
	}
	defer func() {
		if closeError := db.Close(); closeError != nil {
			fmt.Println("Error closing database", closeError)
			if err == nil {
				err = closeError
			}
		}
	}()

	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("Failed to connect to the database: %v", err))
	}

	err = uc.MigrateUp(db)
	if err != nil {
		panic(fmt.Sprintf("Failed to run migrations: %v", err))
	}
}

func (uc *setUpUseCase) MigrateUp(db *sql.DB) error {
	files, err := fs.ReadDir(uc.Migrations, "migrations/shared")
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}

		fileSections := strings.Split(file.Name(), ".")
		if fileSections[FILE_UP_MIGRATION_INDEX] != "up" {
			continue // skip files that are not "up" migrations
		}

		migration, err := fs.ReadFile(uc.Migrations, "migrations/shared/"+file.Name())
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		fmt.Println(string(migration))

		_, err = db.Exec(string(migration))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
		}
	}

	return nil
}
