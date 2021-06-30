package v1

import (
	"context"
	"fmt"
	"time"

	"github.com/teambition/rrule-go"

	"github.com/itiky/charge_scheduler/common"
	"github.com/itiky/charge_scheduler/schema"
)

func (svc Scheduler) AddSingleEvent(ctx context.Context, eventType schema.SingleEventType, eventStart time.Time, endDayHours, endDayMinutes uint) error {
	// Common check
	if err := svc.validateEventInput(eventType, eventStart, endDayHours, endDayMinutes); err != nil {
		return err
	}

	newEvent := &event{
		Start: eventStart,
		End:   cloneTimeWithHourAndMinutes(eventStart, endDayHours, endDayMinutes),
	}

	// Get existing events [eventStart -1 day : eventEnd +1 day]
	rangeStart, rangeEnd := newEvent.Start.Add(-24*time.Hour), newEvent.End.Add(24*time.Hour)
	existingGreenEvents, existingRedEvents, err := svc.getGreenRedEvents(ctx, rangeStart, rangeEnd)
	if err != nil {
		return fmt.Errorf("svc.getAllRangedEvents: %w", err)
	}

	// Pick target existing events group
	var existingTargetEvents []*event
	if eventType == schema.SingleEventTypeAvailable {
		existingTargetEvents = existingGreenEvents
	} else {
		existingTargetEvents = existingRedEvents
	}

	// Check intersection
	for _, existingEvent := range existingTargetEvents {
		if svc.checkEventsIntersect(newEvent, existingEvent) {
			return fmt.Errorf("event intersects with an existing event (%d: %s): %w", existingEvent.Id, existingEvent.Type, common.ErrInvalidInput)
		}
	}

	// Create
	event := schema.SingleEvent{
		Type:          eventType,
		StartDateTime: eventStart,
		EndHours:      endDayHours,
		EndMinutes:    endDayMinutes,
		CreatedAt:     time.Now().UTC(),
	}
	if _, err := svc.eventsSt.CreateSingleEvent(ctx, event); err != nil {
		return fmt.Errorf("svc.eventsSt.CreateSingleEvent: %w", err)
	}
	svc.logger.Info().Stringer("event", event).Msgf("event created")

	return nil
}

func (svc Scheduler) AddPeriodicEvent(ctx context.Context, eventType schema.SingleEventType, eventStart time.Time, endDayHours, endDayMinutes uint) error {
	// Common check
	if err := svc.validateEventInput(eventType, eventStart, endDayHours, endDayMinutes); err != nil {
		return err
	}

	rule, err := rrule.NewRRule(rrule.ROption{
		Freq:    rrule.WEEKLY,
		Count:   0,
		Dtstart: eventStart,
	})
	if err != nil {
		return fmt.Errorf("rrule.NewRRule: %w", err)
	}

	// Get existing events [eventStart -1 day : eventEnd +1 week +1 day]
	rangeStart, rangeEnd := eventStart.Add(-24*time.Hour), cloneTimeWithHourAndMinutes(eventStart, endDayHours, endDayMinutes).Add((7*24+24)*time.Hour)
	existingGreenEvents, existingRedEvents, err := svc.getGreenRedEvents(ctx, rangeStart, rangeEnd)
	if err != nil {
		return fmt.Errorf("svc.getAllRangedEvents: %w", err)
	}

	// Pick target existing events group
	var existingTargetEvents []*event
	if eventType == schema.SingleEventTypeAvailable {
		existingTargetEvents = existingGreenEvents
	} else {
		existingTargetEvents = existingRedEvents
	}

	// Check intersection with periodic events
	for _, newEventStart := range rule.Between(rangeStart, rangeEnd, true) {
		newEvent := &event{
			Start: newEventStart,
			End:   cloneTimeWithHourAndMinutes(newEventStart, endDayHours, endDayMinutes),
		}

		for _, existingEvent := range existingTargetEvents {
			if svc.checkEventsIntersect(newEvent, existingEvent) {
				return fmt.Errorf("event intersects with an existing event (%d: %s): %w", existingEvent.Id, existingEvent.Type, common.ErrInvalidInput)
			}
		}
	}

	// Create
	event := schema.PeriodicEvent{
		Type:       eventType,
		Rrule:      *rule,
		EndHours:   endDayHours,
		EndMinutes: endDayMinutes,
		CreatedAt:  time.Now().UTC(),
	}
	if _, err := svc.eventsSt.CreatePeriodicEvent(ctx, event); err != nil {
		return fmt.Errorf("svc.eventsSt.CreatePeriodicEvent: %w", err)
	}
	svc.logger.Info().Stringer("event", event).Msgf("event created")

	return nil
}

func (svc Scheduler) validateEventInput(eventType schema.SingleEventType, eventStart time.Time, endDayHours, endDayMinutes uint) (retErr error) {
	// Input checks
	if !eventType.IsValid() {
		retErr = fmt.Errorf("%s: invalid: %w", "eventType", common.ErrInvalidInput)
		return
	}
	if eventStart.IsZero() {
		retErr = fmt.Errorf("%s: zero: %w", "eventStart", common.ErrInvalidInput)
		return
	}

	if endDayHours > 23 {
		retErr = fmt.Errorf("%s: must be LTE 23: %w", "endDayHours", common.ErrInvalidInput)
		return
	}
	if endDayMinutes > 59 {
		retErr = fmt.Errorf("%s: must be LTE 59: %w", "endDayMinutes", common.ErrInvalidInput)
		return
	}

	eventEnd := time.Date(eventStart.Year(), eventStart.Month(), eventStart.Day(), int(endDayHours), int(endDayMinutes), eventStart.Second(), eventStart.Nanosecond(), eventStart.Location())
	if eventEnd.Equal(eventStart) || eventEnd.Before(eventStart) {
		retErr = fmt.Errorf("%s: must eventEnd must be after the eventStart: %w", "endDayHours / endDayMinutes", common.ErrInvalidInput)
		return
	}

	return
}
