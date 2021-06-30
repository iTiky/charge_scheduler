package sqlite

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/itiky/charge_scheduler/storage/events"
	"github.com/itiky/charge_scheduler/storage/sqlite_base"
)

var _ events.EventsStorage = (*EventsStorage)(nil)

type EventsStorage struct {
	*sqlite_base.SQLiteBase
	logger zerolog.Logger
}

// nolint:errcheck
func (s EventsStorage) DropData(ctx context.Context) error {
	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("s.Db.BeginTx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM single_events"); err != nil {
		return fmt.Errorf("tx.Exec (single_events): %w", err)
	}
	if _, err := tx.Exec("DELETE FROM periodic_events"); err != nil {
		return fmt.Errorf("tx.Exec (periodic_events): %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil
}

func NewEventsStorage(base *sqlite_base.SQLiteBase) (*EventsStorage, error) {
	if base == nil {
		return nil, fmt.Errorf("%s: nil", "base")
	}

	storage := &EventsStorage{
		SQLiteBase: base,
		logger:     base.Logger.With().Str("repository", "events").Logger(),
	}

	return storage, nil
}
