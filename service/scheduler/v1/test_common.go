package v1

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/itiky/charge_scheduler/service/scheduler/testutil"
	eventsSt "github.com/itiky/charge_scheduler/storage/events/sqlite"
	"github.com/itiky/charge_scheduler/storage/sqlite_base"
)

func NewTestResource(baseSt *sqlite_base.SQLiteBase) (*testutil.SchedulerServiceTestResource, error) {
	stRes, err := eventsSt.NewTestResource(baseSt)
	if err != nil {
		return nil, fmt.Errorf("eventsSt.NewTestResource: %w", err)
	}

	schedulerSvc, err := NewScheduler(zerolog.Nop(), stRes.Storage)
	if err != nil {
		return nil, fmt.Errorf("NewScheduler: %w", err)
	}

	return &testutil.SchedulerServiceTestResource{
		Svc:        schedulerSvc,
		StorageRes: stRes,
	}, nil
}
