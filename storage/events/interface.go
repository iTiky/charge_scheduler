package events

import (
	"context"
	"time"

	"github.com/itiky/charge_scheduler/schema"
)

// EventsStorage provides events repository operations.
type EventsStorage interface {
	// CreateSingleEvent creates a new schema.SingleEvent object and returns its ID.
	CreateSingleEvent(ctx context.Context, obj schema.SingleEvent) (int64, error)
	// CreatePeriodicEvent creates a new schema.PeriodicEvent object and returns its ID.
	CreatePeriodicEvent(ctx context.Context, obj schema.PeriodicEvent) (int64, error)
	// GetSingleEvent gets a schema.SingleEvent by ID (if exists).
	GetSingleEvent(ctx context.Context, id int64) (*schema.SingleEvent, error)
	// GetPeriodicEvent gets a schema.PeriodicEvent by ID (if exists).
	GetPeriodicEvent(ctx context.Context, id int64) (*schema.PeriodicEvent, error)
	// GetSingleEventsWithinRange gets a schema.SingleEvent list filtered by eventStart time range.
	GetSingleEventsWithinRange(ctx context.Context, rangeStart, rangeEnd time.Time) ([]schema.SingleEvent, error)
	// GetAllPeriodicEvents gets all schema.PeriodicEvent objects.
	GetAllPeriodicEvents(ctx context.Context) ([]schema.PeriodicEvent, error)
	// DropData removes all storage data (for debug purposes only)
	DropData(ctx context.Context) error
}
