package v1

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/itiky/charge_scheduler/service/scheduler"
	"github.com/itiky/charge_scheduler/storage/events"
)

var _ scheduler.Scheduler = (*Scheduler)(nil)

type Scheduler struct {
	logger   zerolog.Logger
	eventsSt events.EventsStorage
}

func NewScheduler(logger zerolog.Logger, eventsSt events.EventsStorage) (*Scheduler, error) {
	if eventsSt == nil {
		return nil, fmt.Errorf("%s: nil", "eventsSt")
	}

	return &Scheduler{
		logger:   logger.With().Str("component", "Scheduler service").Logger(),
		eventsSt: eventsSt,
	}, nil
}
