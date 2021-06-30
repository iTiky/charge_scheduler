package scheduler

import (
	"context"
	"time"

	"github.com/itiky/charge_scheduler/schema"
)

type Scheduler interface {
	// AddPeriodicEvent creates a new non-intersecting with existing events schema.SingleEvent.
	AddSingleEvent(ctx context.Context, eventType schema.SingleEventType, eventStart time.Time, endDayHours, endDayMinutes uint) error
	// AddPeriodicEvent creates a new non-intersecting with existing events schema.PeriodicEvent with weekly period.
	AddPeriodicEvent(ctx context.Context, eventType schema.SingleEventType, eventStart time.Time, endDayHours, endDayMinutes uint) error
	// GetAvailableAgenda returns available charging slots for specified period and desired charging duration.
	GetAvailableAgenda(ctx context.Context, periodStart time.Time, periodDur, desiredDur time.Duration) (schema.AgendaResults, error)
	// GetEvents returns registered within specified range singleEvents and all available periodic events.
	GetEvents(ctx context.Context, periodStart, periodEnd time.Time) ([]schema.SingleEvent, []schema.PeriodicEvent, error)
}
