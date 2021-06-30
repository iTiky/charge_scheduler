package sqlite

import (
	"time"

	"github.com/stretchr/testify/require"
	"github.com/teambition/rrule-go"

	"github.com/itiky/charge_scheduler/schema"
)

func (s *StorageTestSuite) Test_SingleEvent() {
	t := s.T()
	ctx := s.ctx
	targetSt := s.r.Storage

	// Init fixtures
	now := time.Now().UTC()
	events := []schema.SingleEvent{
		{
			Id:            1,
			Type:          schema.SingleEventTypeAvailable,
			StartDateTime: now,
			EndHours:      12,
			EndMinutes:    30,
			CreatedAt:     now,
		},
		{
			Id:            2,
			Type:          schema.SingleEventTypeOccupied,
			StartDateTime: now.Add(1 * time.Minute),
			EndHours:      10,
			EndMinutes:    0,
			CreatedAt:     now,
		},
	}

	// ok: GetSingleEvent: non-existing
	{
		res, err := targetSt.GetSingleEvent(ctx, 1)
		require.NoError(t, err)
		require.Nil(t, res)
	}

	// ok: CreateSingleEvent / GetSingleEvent
	{
		for _, event := range events {
			id, err := targetSt.CreateSingleEvent(ctx, event)
			require.NoError(t, err)
			require.NotEmpty(t, id)

			res, err := targetSt.GetSingleEvent(ctx, id)
			require.NoError(t, err)
			require.NotNil(t, res)
			require.Equal(t, event, *res)
		}
	}

	// ok: GetSingleEventsWithinRange: empty
	{
		res, err := targetSt.GetSingleEventsWithinRange(ctx, now.Add(5*time.Minute), now.Add(10*time.Minute))
		require.NoError(t, err)
		require.Empty(t, res)
	}

	// ok: GetSingleEventsWithinRange: filtered
	{
		res, err := targetSt.GetSingleEventsWithinRange(ctx, now, now.Add(30*time.Second))
		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, events[0:1], res)
	}
}

func (s *StorageTestSuite) Test_PeriodicEvent() {
	t := s.T()
	ctx := s.ctx
	targetSt := s.r.Storage

	// Init fixtures
	now := time.Now().UTC()

	rule1, err := rrule.NewRRule(rrule.ROption{
		Freq:    rrule.DAILY,
		Count:   10,
		Dtstart: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	rule2, err := rrule.NewRRule(rrule.ROption{
		Freq:    rrule.WEEKLY,
		Count:   10,
		Dtstart: time.Date(2020, 1, 5, 12, 30, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	events := []schema.PeriodicEvent{
		{
			Id:         1,
			Type:       schema.SingleEventTypeAvailable,
			Rrule:      *rule1,
			EndHours:   15,
			EndMinutes: 30,
			CreatedAt:  now,
		},
		{
			Id:         2,
			Type:       schema.SingleEventTypeOccupied,
			Rrule:      *rule2,
			EndHours:   0,
			EndMinutes: 0,
			CreatedAt:  now,
		},
	}

	// ok: GetPeriodicEvent: non-existing
	{
		res, err := targetSt.GetPeriodicEvent(ctx, 1)
		require.NoError(t, err)
		require.Nil(t, res)
	}

	// ok: CreatePeriodicEvent / GetPeriodicEvent
	{
		for _, event := range events {
			id, err := targetSt.CreatePeriodicEvent(ctx, event)
			require.NoError(t, err)
			require.NotEmpty(t, id)

			res, err := targetSt.GetPeriodicEvent(ctx, id)
			require.NoError(t, err)
			require.NotNil(t, res)
			require.Equal(t, event, *res)
		}
	}

	// ok: GetAllPeriodicEvents
	{
		res, err := targetSt.GetAllPeriodicEvents(ctx)
		require.NoError(t, err)
		require.Len(t, res, 2)
		require.ElementsMatch(t, events, res)
	}
}
