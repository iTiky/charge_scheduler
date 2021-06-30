package sqlite_base

import (
	"fmt"
	"path"

	"github.com/rs/zerolog"
)

func SetupTempSQLiteBase(tmpDir string) (retStorage *SQLiteBase, retErr error) {
	storage, err := NewSQLiteBase(zerolog.Nop(), path.Join(tmpDir, "sqlite.db"))
	if err != nil {
		retErr = fmt.Errorf("NewSQLiteBase: %w", err)
		return
	}
	retStorage = storage

	if err := storage.Migrate(); err != nil {
		retErr = fmt.Errorf("storage.Migrate: %w", err)
		return
	}

	return
}
