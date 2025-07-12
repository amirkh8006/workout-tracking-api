package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx" , "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")

	if err != nil {
		return nil, fmt.Errorf("db open %w", err)
	}

	fmt.Println("connected to database...")

	return db, nil
}

func MigrateFs(db *sql.DB, migrationFs fs.FS, dir string) error {
	goose.SetBaseFS(migrationFs)
	defer func() {
		goose.SetBaseFS(nil)
	}()

	return migrate(db, dir)
}

func migrate(db *sql.DB, directory string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, directory)
	if err != nil  {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}