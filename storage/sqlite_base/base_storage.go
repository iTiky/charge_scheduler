package sqlite_base

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	//_ "github.com/golang-migrate/migrate/v4/source/file"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"

	"github.com/itiky/charge_scheduler/storage/sqlite_base/resources"
)

const DefMigrationPath = "file://${GOPATH}/src/github.com/itiky/charge_scheduler/storage/sqlite_base/migrations"

type SQLiteBase struct {
	Db     *sqlx.DB
	Logger zerolog.Logger
}

func (s SQLiteBase) Close() error {
	if err := s.Db.Close(); err != nil {
		return fmt.Errorf("s.db.Close: %w", err)
	}

	return nil
}

func (s SQLiteBase) Migrate() error {
	dbDriver, err := sqlite3.WithInstance(s.Db.DB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("driver init: sqlite3.WithInstance: %w", err)
	}

	migrationsRes := bindata.Resource(resources.AssetNames(),
		func(name string) ([]byte, error) {
			return resources.Asset(name)
		},
	)
	resDriver, err := bindata.WithInstance(migrationsRes)
	if err != nil {
		return fmt.Errorf("bindata.WithInstance: %w", err)
	}

	migrateManager, err := migrate.NewWithInstance("go-bindata", resDriver, "sqlite3", dbDriver)
	if err != nil {
		return fmt.Errorf("migration manager init: migrate.NewWithInstance: %w", err)
	}

	// //Files source version
	// //import _ "github.com/golang-migrate/migrate/v4/source/file"
	//migrateManager, err := migrate.NewWithDatabaseInstance(migrationsPath, "sqlite3", driver)
	//if err != nil {
	//	return fmt.Errorf("migration manager init: migrate.NewWithDatabaseInstance(%s): %w", migrationsPath, err)
	//}

	prevVersion, prevMigrationFailed, err := migrateManager.Version()
	if err != nil {
		if !errors.Is(err, migrate.ErrNilVersion) {
			return fmt.Errorf("reading previous migration version: migrateManager.Version: %w", err)
		}
	}
	if prevMigrationFailed {
		return fmt.Errorf("previous migration (%d) failed, can't continue", prevVersion)
	}

	if err := migrateManager.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migration failed: migrateManager.Up(): %w", err)
		}
	}

	curVersion, _, err := migrateManager.Version()
	if err != nil {
		if !errors.Is(err, migrate.ErrNilVersion) {
			return fmt.Errorf("reading current migration version: migrateManager.Version: %w", err)
		}
	}

	if curVersion != prevVersion {
		s.Logger.Info().Msgf("Migrated %d -> %d", prevVersion, curVersion)
	} else {
		s.Logger.Debug().Msg("Migration skipped")
	}

	return nil
}

func NewSQLiteBase(logger zerolog.Logger, filePath string) (*SQLiteBase, error) {
	db, err := sqlx.Open("sqlite3", filePath)
	if err != nil {
		return nil, fmt.Errorf("sql.Open(%s): %w", filePath, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return &SQLiteBase{
		Db:     db,
		Logger: logger.With().Str("component", "SQLite storage").Logger(),
	}, nil
}
