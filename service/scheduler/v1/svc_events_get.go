package v1

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/itiky/charge_scheduler/common"
	"github.com/itiky/charge_scheduler/schema"
)

func (svc Scheduler) getGreenRedEvents(ctx context.Context, periodStart, periodEnd time.Time) (retGreenEvents, retRedEvents []*event, retErr error) {
	events, err := svc.getAllRangedEvents(ctx, periodStart, periodEnd)
	if err != nil {
		retErr = fmt.Errorf("svc.getAllRangedEvents: %w", err)
		return
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Start.Before(events[j].Start)
	})

	var prevGreen, prevRed *event
	for i := 0; i < len(events); i++ {
		event := &events[i]

		switch event.Type {
		case schema.SingleEventTypeAvailable:
			if prevGreen != nil {
				event.Prev = prevGreen
				prevGreen.Next = event
			}
			prevGreen = event
			retGreenEvents = append(retGreenEvents, event)
		case schema.SingleEventTypeOccupied:
			if prevRed != nil {
				event.Prev = prevRed
				prevRed.Next = event
			}
			prevRed = event
			retRedEvents = append(retRedEvents, event)
		default:
			retErr = fmt.Errorf("unsupported event.Type: %s", event.Type)
			return
		}
	}

	return
}

func (svc Scheduler) getAllRangedEvents(ctx context.Context, periodStart, periodEnd time.Time) (retEvents []event, retErr error) {
	sEvents, err := svc.getSingleRangedEvents(ctx, periodStart, periodEnd)
	if err != nil {
		retErr = fmt.Errorf("svc.getSingleRangedEvents: %w", err)
		return
	}

	pEvents, err := svc.getPeriodicRangedEvents(ctx, periodStart, periodEnd)
	if err != nil {
		retErr = fmt.Errorf("svc.getPeriodicRangedEvents: %w", err)
		return
	}

	retEvents = append(sEvents, pEvents...)

	return
}

func (svc Scheduler) getSingleRangedEvents(ctx context.Context, periodStart, periodEnd time.Time) (retEvents []event, retErr error) {
	dbEvents, err := svc.eventsSt.GetSingleEventsWithinRange(ctx, periodStart, periodEnd)
	if err != nil {
		retErr = fmt.Errorf("svc.eventsSt.GetSingleEventsWithinRange(%s, %s): %w", periodStart.Format(common.TimeFmt), periodEnd.Format(common.TimeFmt), err)
		return
	}

	retEvents = make([]event, 0, len(dbEvents))
	for _, dbEvent := range dbEvents {
		retEvents = append(retEvents, event{
			Id:    dbEvent.Id,
			Type:  dbEvent.Type,
			Start: dbEvent.StartDateTime,
			End:   cloneTimeWithHourAndMinutes(dbEvent.StartDateTime, dbEvent.EndHours, dbEvent.EndMinutes),
		})
	}

	return
}

func (svc Scheduler) getPeriodicRangedEvents(ctx context.Context, periodStart, periodEnd time.Time) (retEvents []event, retErr error) {
	dbEvents, err := svc.eventsSt.GetAllPeriodicEvents(ctx)
	if err != nil {
		retErr = fmt.Errorf("svc.eventsSt.GetAllPeriodicEvents: %w", err)
		return
	}

	retEvents = make([]event, 0)
	for _, dbEvent := range dbEvents {
		for _, t := range dbEvent.Rrule.Between(periodStart, periodEnd, true) {
			retEvents = append(retEvents, event{
				Id:    dbEvent.Id,
				Type:  dbEvent.Type,
				Start: t,
				End:   cloneTimeWithHourAndMinutes(t, dbEvent.EndHours, dbEvent.EndMinutes),
			})
		}
	}

	return
}

func cloneTimeWithHourAndMinutes(t time.Time, h, m uint) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), int(h), int(m), t.Second(), t.Nanosecond(), t.Location())
}
