package v1

import (
	"time"

	"github.com/stretchr/testify/require"

	"github.com/itiky/charge_scheduler/schema"
)

func (s *ServiceTestSuite) Test_LogicTest1() {
	t := s.T()
	ctx := s.ctx
	targetSvc := s.r.Svc.(*Scheduler)
	require.NoError(t, s.r.StorageRes.Storage.DropData(ctx))

	// Init fixtures
	{
		require.NoError(t, targetSvc.AddPeriodicEvent(ctx,
			schema.SingleEventTypeAvailable,
			time.Date(2014, 8, 4, 9, 30, 0, 0, time.UTC),
			13, 30,
		))

		require.NoError(t, targetSvc.AddSingleEvent(ctx,
			schema.SingleEventTypeOccupied,
			time.Date(2014, 8, 11, 10, 30, 0, 0, time.UTC),
			11, 30,
		))
	}

	// Check
	{
		agendas, err := targetSvc.GetAvailableAgenda(ctx,
			time.Date(2014, 8, 10, 0, 0, 0, 0, time.UTC),
			10*dayDur,
			30*time.Minute,
		)
		require.NoError(t, err)

		t.Logf("Results:\n%s", agendas.String())

		require.Len(t, agendas, 10)
		//
		require.Equal(t, agendas[0].Date, time.Date(2014, 8, 10, 0, 0, 0, 0, time.UTC))
		require.Len(t, agendas[0].TimeSlots, 0)
		//
		require.Equal(t, agendas[1].Date, time.Date(2014, 8, 11, 0, 0, 0, 0, time.UTC))
		require.Len(t, agendas[1].TimeSlots, 6)
		require.Equal(t, agendas[1].TimeSlots[0].Start, agendas[1].Date.Add(9*time.Hour).Add(30*time.Minute))
		require.Equal(t, agendas[1].TimeSlots[1].Start, agendas[1].Date.Add(10*time.Hour).Add(0*time.Minute))
		require.Equal(t, agendas[1].TimeSlots[2].Start, agendas[1].Date.Add(11*time.Hour).Add(30*time.Minute))
		require.Equal(t, agendas[1].TimeSlots[3].Start, agendas[1].Date.Add(12*time.Hour).Add(0*time.Minute))
		require.Equal(t, agendas[1].TimeSlots[4].Start, agendas[1].Date.Add(12*time.Hour).Add(30*time.Minute))
		require.Equal(t, agendas[1].TimeSlots[5].Start, agendas[1].Date.Add(13*time.Hour).Add(0*time.Minute))
		//
		require.Equal(t, agendas[2].Date, time.Date(2014, 8, 12, 0, 0, 0, 0, time.UTC))
		require.Len(t, agendas[2].TimeSlots, 0)
		//
		require.Equal(t, agendas[3].Date, time.Date(2014, 8, 13, 0, 0, 0, 0, time.UTC))
		require.Len(t, agendas[3].TimeSlots, 0)
		//
		require.Equal(t, agendas[4].Date, time.Date(2014, 8, 14, 0, 0, 0, 0, time.UTC))
		require.Len(t, agendas[4].TimeSlots, 0)
		//
		require.Equal(t, agendas[5].Date, time.Date(2014, 8, 15, 0, 0, 0, 0, time.UTC))
		require.Len(t, agendas[5].TimeSlots, 0)
		//
		require.Equal(t, agendas[6].Date, time.Date(2014, 8, 16, 0, 0, 0, 0, time.UTC))
		require.Len(t, agendas[6].TimeSlots, 0)
		//
		require.Equal(t, agendas[7].Date, time.Date(2014, 8, 17, 0, 0, 0, 0, time.UTC))
		require.Len(t, agendas[7].TimeSlots, 0)
		//
		require.Equal(t, agendas[8].Date, time.Date(2014, 8, 18, 0, 0, 0, 0, time.UTC))
		require.Len(t, agendas[8].TimeSlots, 8)
		require.Equal(t, agendas[8].TimeSlots[0].Start, agendas[8].Date.Add(9*time.Hour).Add(30*time.Minute))
		require.Equal(t, agendas[8].TimeSlots[1].Start, agendas[8].Date.Add(10*time.Hour).Add(0*time.Minute))
		require.Equal(t, agendas[8].TimeSlots[2].Start, agendas[8].Date.Add(10*time.Hour).Add(30*time.Minute))
		require.Equal(t, agendas[8].TimeSlots[3].Start, agendas[8].Date.Add(11*time.Hour).Add(0*time.Minute))
		require.Equal(t, agendas[8].TimeSlots[4].Start, agendas[8].Date.Add(11*time.Hour).Add(30*time.Minute))
		require.Equal(t, agendas[8].TimeSlots[5].Start, agendas[8].Date.Add(12*time.Hour).Add(0*time.Minute))
		require.Equal(t, agendas[8].TimeSlots[6].Start, agendas[8].Date.Add(12*time.Hour).Add(30*time.Minute))
		require.Equal(t, agendas[8].TimeSlots[7].Start, agendas[8].Date.Add(13*time.Hour).Add(0*time.Minute))
		//
		require.Equal(t, agendas[9].Date, time.Date(2014, 8, 19, 0, 0, 0, 0, time.UTC))
		require.Len(t, agendas[9].TimeSlots, 0)
	}
}
