package sqlite

import (
	"fmt"

	"github.com/itiky/charge_scheduler/storage/events/testutil"
	"github.com/itiky/charge_scheduler/storage/sqlite_base"
)

func NewTestResource(baseSt *sqlite_base.SQLiteBase) (*testutil.EventsStorageTestResource, error) {
	st, err := NewEventsStorage(baseSt)
	if err != nil {
		return nil, fmt.Errorf("NewEventsStorage: %w", err)
	}

	return &testutil.EventsStorageTestResource{
		Storage: st,
	}, nil
}
