package v1

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/itiky/charge_scheduler/common"
	"github.com/itiky/charge_scheduler/schema"
)

func (s *ServiceTestSuite) Test_Create_NoOverlap() {
	t := s.T()
	ctx := s.ctx
	targetSvc := s.r.Svc.(*Scheduler)
	require.NoError(t, s.r.StorageRes.Storage.DropData(ctx))

	// fail: wrong inputs
	{
		now := time.Date(2000, 1, 1, 15, 30, 0, 0, time.UTC)

		// EventType
		{
			err := targetSvc.AddSingleEvent(ctx, schema.SingleEventType(""), now, 0, 0)
			require.Error(t, err)
			require.True(t, errors.Is(err, common.ErrInvalidInput))
		}
		// endDayHours
		{
			err := targetSvc.AddSingleEvent(ctx, schema.SingleEventTypeAvailable, now, 24, 0)
			require.Error(t, err)
			require.True(t, errors.Is(err, common.ErrInvalidInput))
		}
		// endDayMinutes
		{
			err := targetSvc.AddSingleEvent(ctx, schema.SingleEventTypeAvailable, now, 0, 60)
			require.Error(t, err)
			require.True(t, errors.Is(err, common.ErrInvalidInput))
		}
		// endDayHours < eventStart
		{
			err := targetSvc.AddSingleEvent(ctx, schema.SingleEventTypeAvailable, now, uint(now.Hour()-1), uint(now.Minute()))
			require.Error(t, err)
			require.True(t, errors.Is(err, common.ErrInvalidInput))
		}
		// endDayMinutes < eventStart
		{
			now := now.Add(10 * time.Minute)
			err := targetSvc.AddSingleEvent(ctx, schema.SingleEventTypeAvailable, now, uint(now.Hour()), uint(now.Minute()-1))
			require.Error(t, err)
			require.True(t, errors.Is(err, common.ErrInvalidInput))
		}
	}

	// ok: AddSingleEvent
	// 15.01.2000 (SAT) 09:00 - 12:00
	{
		eventStart := time.Date(2000, 1, 15, 9, 0, 0, 0, time.UTC)
		require.NoError(t, targetSvc.AddSingleEvent(ctx, schema.SingleEventTypeAvailable, eventStart, 12, 0))
	}

	// fail: AddSingleEvent: intersect
	// 15.01.2000 (SAT) 11:30 - 13:00 -> collide the same day
	{
		eventStart := time.Date(2000, 1, 15, 11, 30, 0, 0, time.UTC)
		err := targetSvc.AddSingleEvent(ctx, schema.SingleEventTypeAvailable, eventStart, 13, 0)
		require.Error(t, err)
		require.True(t, errors.Is(err, common.ErrInvalidInput))
	}

	// ok: AddSingleEvent
	// 15.01.2000 (SAT) 18:00 - 19:30
	{
		eventStart := time.Date(2000, 1, 15, 18, 0, 0, 0, time.UTC)
		require.NoError(t, targetSvc.AddSingleEvent(ctx, schema.SingleEventTypeOccupied, eventStart, 19, 30))
	}

	// fail: AddPeriodicEvent: intersect wint single
	// 08.01.2000 (SAT) 12:00 - 13:00 -> collide the next week
	{
		eventStart := time.Date(2000, 1, 8, 12, 0, 0, 0, time.UTC)
		err := targetSvc.AddPeriodicEvent(ctx, schema.SingleEventTypeAvailable, eventStart, 13, 0)
		require.Error(t, err)
		require.True(t, errors.Is(err, common.ErrInvalidInput))
	}

	// ok: AddPeriodicEvent
	// 09.01.2000 (SUN) 12:00 - 13:00
	{
		eventStart := time.Date(2000, 1, 9, 12, 0, 0, 0, time.UTC)
		require.NoError(t, targetSvc.AddPeriodicEvent(ctx, schema.SingleEventTypeAvailable, eventStart, 13, 0))
	}

	// fail: AddPeriodicEvent: intersect wint periodic
	// 23.01.2000 (SAN) 11:00 - 12:15 -> collide the week before
	{
		eventStart := time.Date(2000, 1, 23, 11, 0, 0, 0, time.UTC)
		err := targetSvc.AddPeriodicEvent(ctx, schema.SingleEventTypeAvailable, eventStart, 12, 15)
		require.Error(t, err)
		require.True(t, errors.Is(err, common.ErrInvalidInput))
	}

	// ok: AddPeriodicEvent
	// 18.01.2000 (TUE) 00:00 - 23:59
	{
		eventStart := time.Date(2000, 1, 18, 0, 0, 0, 0, time.UTC)
		require.NoError(t, targetSvc.AddPeriodicEvent(ctx, schema.SingleEventTypeOccupied, eventStart, 23, 59))
	}

	// check the resulting green / red events for the month of January
	{
		greenEvents, redEvents, err := targetSvc.getGreenRedEvents(ctx, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 31, 23, 59, 0, 0, time.UTC))
		require.NoError(t, err)

		require.Len(t, greenEvents, 5)
		require.Len(t, redEvents, 3)

		checkEventTimeRange(t,
			2000, 1, 9, 12, 0,
			2000, 1, 9, 13, 0,
			greenEvents[0],
		)
		checkEventTimeRange(t,
			2000, 1, 15, 9, 0,
			2000, 1, 15, 12, 0,
			greenEvents[1],
		)
		checkEventTimeRange(t,
			2000, 1, 16, 12, 0,
			2000, 1, 16, 13, 0,
			greenEvents[2],
		)
		checkEventTimeRange(t,
			2000, 1, 23, 12, 0,
			2000, 1, 23, 13, 0,
			greenEvents[3],
		)
		checkEventTimeRange(t,
			2000, 1, 30, 12, 0,
			2000, 1, 30, 13, 0,
			greenEvents[4],
		)

		checkEventTimeRange(t,
			2000, 1, 15, 18, 0,
			2000, 1, 15, 19, 30,
			redEvents[0],
		)
		checkEventTimeRange(t,
			2000, 1, 18, 0, 0,
			2000, 1, 18, 23, 59,
			redEvents[1],
		)
		checkEventTimeRange(t,
			2000, 1, 25, 0, 0,
			2000, 1, 25, 23, 59,
			redEvents[2],
		)

		checkEventsAreSortedLinkedList(t, greenEvents)
		checkEventsAreSortedLinkedList(t, redEvents)
	}
}

func (s *ServiceTestSuite) Test_Create_Overlap() {
	t := s.T()
	ctx := s.ctx
	targetSvc := s.r.Svc.(*Scheduler)
	require.NoError(t, s.r.StorageRes.Storage.DropData(ctx))

	// ok / fail: Green: AddPeriodicEvent
	// 21.05.2001 (MON) 09:00 - 11:59
	{
		eventStart := time.Date(2001, 5, 21, 9, 0, 0, 0, time.UTC)
		require.NoError(t, targetSvc.AddPeriodicEvent(ctx, schema.SingleEventTypeAvailable, eventStart, 11, 59))

		// fail: Green: AddPeriodicEvent
		// 28.05.2001 (MON) 10:00 - 13:00
		{
			eventStart := time.Date(2001, 5, 28, 10, 0, 0, 0, time.UTC)
			err := targetSvc.AddPeriodicEvent(ctx, schema.SingleEventTypeAvailable, eventStart, 13, 0)
			require.Error(t, err)
			require.True(t, errors.Is(err, common.ErrInvalidInput))
		}
	}

	// ok / fail: Green: AddSingleEvent
	// 21.05.2001 (MON) 12:00 - 13:00
	{
		eventStart := time.Date(2001, 5, 21, 12, 0, 0, 0, time.UTC)
		require.NoError(t, targetSvc.AddSingleEvent(ctx, schema.SingleEventTypeAvailable, eventStart, 13, 0))

		// fail: Green: AddSingleEvent
		// 21.05.2001 (MON) 13:00 - 14:00
		{
			eventStart := time.Date(2001, 5, 21, 13, 0, 0, 0, time.UTC)
			err := targetSvc.AddSingleEvent(ctx, schema.SingleEventTypeAvailable, eventStart, 14, 0)
			require.Error(t, err)
			require.True(t, errors.Is(err, common.ErrInvalidInput))
		}
	}

	// ok: Green: AddSingleEvent
	// 22.05.2001 (TUE) 15:00 - 18:00
	{
		eventStart := time.Date(2001, 5, 22, 15, 0, 0, 0, time.UTC)
		require.NoError(t, targetSvc.AddSingleEvent(ctx, schema.SingleEventTypeAvailable, eventStart, 18, 0))
	}

	// ok / fail: Red: AddPeriodicEvent (overlaps periodic and single Green)
	// 21.05.2001 (MON) 08:00 - 14:00
	{
		eventStart := time.Date(2001, 5, 21, 8, 0, 0, 0, time.UTC)
		require.NoError(t, targetSvc.AddPeriodicEvent(ctx, schema.SingleEventTypeOccupied, eventStart, 14, 0))

		// fail: Red: AddPeriodicEvent
		// 14.05.2001 (MON) 09:00 - 10:00
		{
			eventStart := time.Date(2001, 5, 14, 9, 0, 0, 0, time.UTC)
			err := targetSvc.AddPeriodicEvent(ctx, schema.SingleEventTypeOccupied, eventStart, 10, 0)
			require.Error(t, err)
			require.True(t, errors.Is(err, common.ErrInvalidInput))
		}
	}

	// ok / fail: Red: AddSingleEvent (overlaps single Green)
	// 22.05.2001 (TUE) 16:00 - 17:00
	{
		eventStart := time.Date(2001, 5, 22, 16, 0, 0, 0, time.UTC)
		require.NoError(t, targetSvc.AddSingleEvent(ctx, schema.SingleEventTypeOccupied, eventStart, 17, 0))

		// fail: Red: AddSingleEvent
		// 22.05.2001 (TUE) 16:00 - 16:30
		{
			eventStart := time.Date(2001, 5, 22, 16, 0, 0, 0, time.UTC)
			err := targetSvc.AddSingleEvent(ctx, schema.SingleEventTypeOccupied, eventStart, 16, 30)
			require.Error(t, err)
			require.True(t, errors.Is(err, common.ErrInvalidInput))
		}
	}

	// check the resulting green / red events for two weeks
	{
		greenEvents, redEvents, err := targetSvc.getGreenRedEvents(ctx, time.Date(2001, 5, 20, 0, 0, 0, 0, time.UTC), time.Date(2001, 6, 2, 23, 59, 0, 0, time.UTC))
		require.NoError(t, err)

		require.Len(t, greenEvents, 4)
		require.Len(t, redEvents, 3)

		checkEventTimeRange(t,
			2001, 5, 21, 9, 0,
			2001, 5, 21, 11, 59,
			greenEvents[0],
		)
		checkEventTimeRange(t,
			2001, 5, 21, 12, 0,
			2001, 5, 21, 13, 0,
			greenEvents[1],
		)
		checkEventTimeRange(t,
			2001, 5, 22, 15, 0,
			2001, 5, 22, 18, 0,
			greenEvents[2],
		)
		checkEventTimeRange(t,
			2001, 5, 28, 9, 0,
			2001, 5, 28, 11, 59,
			greenEvents[3],
		)
		//
		checkEventTimeRange(t,
			2001, 5, 21, 8, 0,
			2001, 5, 21, 14, 0,
			redEvents[0],
		)
		checkEventTimeRange(t,
			2001, 5, 22, 16, 0,
			2001, 5, 22, 17, 0,
			redEvents[1],
		)
		checkEventTimeRange(t,
			2001, 5, 28, 8, 0,
			2001, 5, 28, 14, 0,
			redEvents[2],
		)

		checkEventsAreSortedLinkedList(t, greenEvents)
		checkEventsAreSortedLinkedList(t, redEvents)
	}
}

func checkEventTimeRange(t *testing.T, expStartYear, expStartMonth, expStartDay, expStartHour, expStartMinute, expEndYear, expEndMonth, expEndDay, expEndHour, expEndMinute int, rcv *event) {
	require.Equal(t, expStartYear, rcv.Start.Year())
	require.Equal(t, expStartMonth, int(rcv.Start.Month()))
	require.Equal(t, expStartDay, rcv.Start.Day())
	require.Equal(t, expStartHour, rcv.Start.Hour())
	require.Equal(t, expStartMinute, rcv.Start.Minute())
	//
	require.Equal(t, expEndYear, rcv.End.Year())
	require.Equal(t, expEndMonth, int(rcv.End.Month()))
	require.Equal(t, expEndDay, rcv.End.Day())
	require.Equal(t, expEndHour, rcv.End.Hour())
	require.Equal(t, expEndMinute, rcv.End.Minute())
}

func checkEventsAreSortedLinkedList(t *testing.T, events []*event) {
	for i := 0; i < len(events); i++ {
		event := events[i]

		if i == 0 {
			require.Nil(t, event.Prev)
		}
		if i == len(events)-1 {
			require.Nil(t, event.Next)
		}

		if prevEvent := event.Prev; prevEvent != nil {
			require.Equal(t, event.Prev, prevEvent)
			require.Equal(t, prevEvent.Next, event)
			require.True(t, event.Start.Equal(prevEvent.End) || event.Start.After(prevEvent.End))
			require.Equal(t, event.Type, prevEvent.Type)
		}

		if nextEvent := event.Next; nextEvent != nil {
			require.Equal(t, event.Next, nextEvent)
			require.Equal(t, nextEvent.Prev, event)
			require.True(t, event.End.Equal(nextEvent.Start) || event.End.Before(nextEvent.Start))
			require.Equal(t, event.Type, nextEvent.Type)
		}

		require.True(t, event.End.After(event.Start))
	}
}
