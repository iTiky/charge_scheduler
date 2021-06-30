package v1

import (
	"context"
	"fmt"
	"time"

	"github.com/itiky/charge_scheduler/common"
	"github.com/itiky/charge_scheduler/schema"
)

func (svc Scheduler) GetEvents(ctx context.Context, periodStart, periodEnd time.Time) (retSingleEvents []schema.SingleEvent, retPeriodicEvents []schema.PeriodicEvent, retErr error) {
	// Input checks
	if periodStart.IsZero() {
		retErr = fmt.Errorf("%s: zero: %w", "periodStart", common.ErrInvalidInput)
		return
	}
	if periodEnd.IsZero() {
		retErr = fmt.Errorf("%s: zero: %w", "periodEnd", common.ErrInvalidInput)
		return
	}
	if periodStart.Equal(periodEnd) || periodStart.After(periodEnd) {
		retErr = fmt.Errorf("%s: periodStart must be LT periodEnd: %w", "periodStart / periodEnd", common.ErrInvalidInput)
		return
	}

	// Get
	sEvents, err := svc.eventsSt.GetSingleEventsWithinRange(ctx, periodStart, periodEnd)
	if err != nil {
		retErr = fmt.Errorf("svc.eventsSt.GetSingleEventsWithinRange(%s, %s): %w", periodStart.Format(common.TimeFmt), periodEnd.Format(common.TimeFmt), err)
		return
	}
	retSingleEvents = sEvents

	pEvents, err := svc.eventsSt.GetAllPeriodicEvents(ctx)
	if err != nil {
		retErr = fmt.Errorf("svc.eventsSt.GetAllPeriodicEvents: %w", err)
		return
	}
	retPeriodicEvents = pEvents

	return
}
