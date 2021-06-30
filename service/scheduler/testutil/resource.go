package testutil

import (
	"github.com/itiky/charge_scheduler/service/scheduler"
	"github.com/itiky/charge_scheduler/storage/events/testutil"
)

type SchedulerServiceTestResource struct {
	Svc        scheduler.Scheduler
	StorageRes *testutil.EventsStorageTestResource
}
