package v1

import (
	"time"

	"github.com/stretchr/testify/require"
)

func (s *ServiceTestSuite) Test_checkEventsIntersect() {
	t := s.T()
	targetSvc := s.r.Svc.(*Scheduler)

	now := time.Now()

	// no intersect
	{
		// e1 < e2
		{
			e1 := &event{Start: now.Add(-1 * time.Second), End: now}
			e2 := &event{Start: now.Add(1 * time.Second), End: now.Add(2 * time.Second)}
			require.False(t, targetSvc.checkEventsIntersect(e1, e2))
		}
		// e1 > e2
		{
			e1 := &event{Start: now.Add(1 * time.Second), End: now.Add(2 * time.Second)}
			e2 := &event{Start: now.Add(-1 * time.Second), End: now}
			require.False(t, targetSvc.checkEventsIntersect(e1, e2))
		}
	}

	// intersect: one point
	{
		// e1.End == e2.Start
		{
			e1 := &event{Start: now.Add(-1 * time.Second), End: now}
			e2 := &event{Start: now, End: now.Add(1 * time.Second)}
			require.True(t, targetSvc.checkEventsIntersect(e1, e2))
		}
		// e1.Start == e2.End
		{
			e1 := &event{Start: now, End: now.Add(1 * time.Second)}
			e2 := &event{Start: now.Add(-1 * time.Second), End: now}
			require.True(t, targetSvc.checkEventsIntersect(e1, e2))
		}
	}

	// intersect: intersect
	{
		// e1 after e2
		{
			e1 := &event{Start: now, End: now.Add(2 * time.Second)}
			e2 := &event{Start: now.Add(-1 * time.Second), End: now.Add(1 * time.Second)}
			require.True(t, targetSvc.checkEventsIntersect(e1, e2))
		}
		// e1 before e2
		{
			e1 := &event{Start: now.Add(-2 * time.Second), End: now}
			e2 := &event{Start: now.Add(-1 * time.Second), End: now.Add(1 * time.Second)}
			require.True(t, targetSvc.checkEventsIntersect(e1, e2))
		}
	}

	// intersect: inner
	{
		// e1 in e2
		{
			e1 := &event{Start: now.Add(-1 * time.Second), End: now.Add(1 * time.Second)}
			e2 := &event{Start: now.Add(-2 * time.Second), End: now.Add(2 * time.Second)}
			require.True(t, targetSvc.checkEventsIntersect(e1, e2))
		}
		// e2 in e1
		{
			e1 := &event{Start: now.Add(-2 * time.Second), End: now.Add(2 * time.Second)}
			e2 := &event{Start: now.Add(-1 * time.Second), End: now.Add(1 * time.Second)}
			require.True(t, targetSvc.checkEventsIntersect(e1, e2))
		}
	}
}
