package v1

import (
	"time"

	"github.com/itiky/charge_scheduler/schema"
)

type event struct {
	Id    int64
	Type  schema.SingleEventType
	Start time.Time
	End   time.Time
	Prev  *event
	Next  *event
}
