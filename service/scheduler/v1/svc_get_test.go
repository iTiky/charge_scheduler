package v1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func (s *ServiceTestSuite) Test_mergeGreenRedEvents() {
	t := s.T()
	targetSvc := s.r.Svc.(*Scheduler)

	// all greens removed
	{
		greens := buildEventsLinkedList(t, []*event{
			buildEvent(
				2000, 1, 5, 9, 0,
				2000, 1, 5, 12, 30,
			),
			buildEvent(
				2000, 1, 6, 12, 0,
				2000, 1, 6, 13, 30,
			),
			buildEvent(
				2000, 1, 6, 14, 0,
				2000, 1, 6, 15, 30,
			),
			buildEvent(
				2000, 1, 7, 15, 30,
				2000, 1, 7, 16, 0,
			),
		})

		reds := buildEventsLinkedList(t, []*event{
			buildEvent(
				2000, 1, 5, 9, 0,
				2000, 1, 5, 12, 30,
			),
			buildEvent(
				2000, 1, 6, 11, 0,
				2000, 1, 6, 16, 30,
			),
			buildEvent(
				2000, 1, 7, 12, 0,
				2000, 1, 7, 18, 0,
			),
		})

		received := targetSvc.mergeGreenRedEvents(greens, reds)
		require.Nil(t, received)
	}

	// all greens preserved, no intersection
	{
		greens := buildEventsLinkedList(t, []*event{
			buildEvent(
				2000, 1, 5, 9, 0,
				2000, 1, 5, 12, 00,
			),
			buildEvent(
				2000, 1, 6, 9, 0,
				2000, 1, 6, 12, 0,
			),
		})

		reds := buildEventsLinkedList(t, []*event{
			buildEvent(
				2000, 1, 5, 13, 0,
				2000, 1, 5, 14, 0,
			),
			buildEvent(
				2000, 1, 6, 7, 0,
				2000, 1, 6, 8, 30,
			),
			buildEvent(
				2000, 1, 7, 12, 0,
				2000, 1, 7, 18, 0,
			),
		})

		expected := buildEventsLinkedList(t, []*event{
			buildEvent(
				2000, 1, 5, 9, 0,
				2000, 1, 5, 12, 00,
			),
			buildEvent(
				2000, 1, 6, 9, 0,
				2000, 1, 6, 12, 0,
			),
		})

		received := targetSvc.mergeGreenRedEvents(greens, reds)
		require.NotNil(t, received)
		checkEventsLinkedListsEqual(t, expected[0], received)
	}

	// complex with splits and removals
	{
		greens := buildEventsLinkedList(t, []*event{
			// Stays "as is"
			buildEvent(
				2000, 1, 1, 9, 0,
				2000, 1, 1, 12, 0,
			),
			// Split into two (red in the middle)
			buildEvent(
				2000, 1, 2, 12, 0,
				2000, 1, 2, 15, 0,
			),
			// Split into three (two reds in the middle)
			buildEvent(
				2000, 1, 3, 9, 0,
				2000, 1, 3, 15, 0,
			),
			// Shrink (red at the start)
			buildEvent(
				2000, 1, 4, 14, 0,
				2000, 1, 4, 16, 0,
			),
			// Shrink (red at the end)
			buildEvent(
				2000, 1, 5, 14, 0,
				2000, 1, 5, 16, 0,
			),
		})

		reds := buildEventsLinkedList(t, []*event{
			buildEvent(
				2000, 1, 2, 13, 0,
				2000, 1, 2, 14, 0,
			),
			buildEvent(
				2000, 1, 3, 10, 0,
				2000, 1, 3, 11, 0,
			),
			buildEvent(
				2000, 1, 3, 13, 0,
				2000, 1, 3, 14, 0,
			),
			buildEvent(
				2000, 1, 4, 13, 0,
				2000, 1, 4, 14, 30,
			),
			buildEvent(
				2000, 1, 5, 15, 30,
				2000, 1, 5, 17, 0,
			),
			buildEvent(
				2000, 1, 6, 9, 0,
				2000, 1, 6, 18, 0,
			),
		})

		expected := buildEventsLinkedList(t, []*event{
			// "As is"
			buildEvent(
				2000, 1, 1, 9, 0,
				2000, 1, 1, 12, 0,
			),
			// "Red in the middle" 1st
			buildEvent(
				2000, 1, 2, 12, 0,
				2000, 1, 2, 13, 0,
			),
			// "Red in the middle" 2nd
			buildEvent(
				2000, 1, 2, 14, 0,
				2000, 1, 2, 15, 0,
			),
			// "Two reds in the middle" 1st
			buildEvent(
				2000, 1, 3, 9, 0,
				2000, 1, 3, 10, 0,
			),
			// "Two reds in the middle" 2nd
			buildEvent(
				2000, 1, 3, 11, 0,
				2000, 1, 3, 13, 0,
			),
			// "Two reds in the middle" 3rd
			buildEvent(
				2000, 1, 3, 14, 0,
				2000, 1, 3, 15, 0,
			),
			// "Red at the start"
			buildEvent(
				2000, 1, 4, 14, 30,
				2000, 1, 4, 16, 0,
			),
			// "Red at the emd"
			buildEvent(
				2000, 1, 5, 14, 0,
				2000, 1, 5, 15, 30,
			),
		})

		received := targetSvc.mergeGreenRedEvents(greens, reds)
		require.NotNil(t, received)
		checkEventsLinkedListsEqual(t, expected[0], received)
	}
}

func (s *ServiceTestSuite) Test_buildAgendaResults() {
	t := s.T()
	targetSvc := s.r.Svc.(*Scheduler)

	// empty
	{
		start := time.Date(2000, 1, 1, 15, 30, 0, 0, time.UTC)
		endDur := 72 * time.Hour
		res := targetSvc.buildAgendaResults(nil, start, endDur, 30*time.Minute)
		require.Len(t, res, 3)

		require.Equal(t, res[0].Date, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
		require.Empty(t, res[0].TimeSlots)
		require.Equal(t, res[1].Date, time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC))
		require.Empty(t, res[1].TimeSlots)
		require.Equal(t, res[2].Date, time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC))
		require.Empty(t, res[2].TimeSlots)
	}

	// complex without pre/post fills
	{
		start := time.Date(2000, 1, 1, 15, 30, 0, 0, time.UTC)
		endDur := 4 * dayDur
		desDur := 30 * time.Minute

		greens := buildEventsLinkedList(t, []*event{
			buildEvent(
				2000, 1, 1, 9, 0,
				2000, 1, 1, 10, 30,
			),
			buildEvent(
				2000, 1, 2, 12, 0,
				2000, 1, 2, 12, 45,
			),
			buildEvent(
				2000, 1, 4, 8, 0,
				2000, 1, 4, 8, 15,
			),
		})

		res := targetSvc.buildAgendaResults(greens[0], start, endDur, desDur)
		require.Len(t, res, 4)

		require.Len(t, res[0].TimeSlots, 3)
		require.Equal(t, res[0].Date, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
		require.Equal(t, res[0].TimeSlots[0].Start, time.Date(2000, 1, 1, 9, 0, 0, 0, time.UTC))
		require.Equal(t, res[0].TimeSlots[0].Duration, desDur)
		require.Equal(t, res[0].TimeSlots[1].Start, time.Date(2000, 1, 1, 9, 30, 0, 0, time.UTC))
		require.Equal(t, res[0].TimeSlots[1].Duration, desDur)
		require.Equal(t, res[0].TimeSlots[2].Start, time.Date(2000, 1, 1, 10, 0, 0, 0, time.UTC))
		require.Equal(t, res[0].TimeSlots[2].Duration, desDur)

		require.Len(t, res[1].TimeSlots, 1)
		require.Equal(t, res[1].Date, time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC))
		require.Equal(t, res[1].TimeSlots[0].Start, time.Date(2000, 1, 2, 12, 0, 0, 0, time.UTC))
		require.Equal(t, res[1].TimeSlots[0].Duration, desDur)

		require.Len(t, res[2].TimeSlots, 0)
		require.Equal(t, res[2].Date, time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC))

		require.Len(t, res[3].TimeSlots, 0)
		require.Equal(t, res[3].Date, time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC))
	}

	// complex with pre/post fills
	{
		start := time.Date(2000, 1, 1, 15, 30, 0, 0, time.UTC)
		endDur := 4 * dayDur
		desDur := 30 * time.Minute

		greens := buildEventsLinkedList(t, []*event{
			buildEvent(
				2000, 1, 3, 9, 0,
				2000, 1, 3, 9, 30,
			),
		})

		res := targetSvc.buildAgendaResults(greens[0], start, endDur, desDur)
		require.Len(t, res, 4)

		require.Len(t, res[0].TimeSlots, 0)
		require.Equal(t, res[0].Date, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))

		require.Len(t, res[1].TimeSlots, 0)
		require.Equal(t, res[1].Date, time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC))

		require.Len(t, res[2].TimeSlots, 1)
		require.Equal(t, res[2].TimeSlots[0].Start, time.Date(2000, 1, 3, 9, 0, 0, 0, time.UTC))
		require.Equal(t, res[2].TimeSlots[0].Duration, desDur)
		require.Equal(t, res[2].Date, time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC))

		require.Len(t, res[3].TimeSlots, 0)
		require.Equal(t, res[3].Date, time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC))
	}
}

func buildEvent(startYear, startMonth, startDay, startHour, startMinute, endYear, endMonth, endDay, endHour, endMinute int) *event {
	return &event{
		Start: time.Date(startYear, time.Month(startMonth), startDay, startHour, startMinute, 0, 0, time.UTC),
		End:   time.Date(endYear, time.Month(endMonth), endDay, endHour, endMinute, 0, 0, time.UTC),
	}
}

func buildEventsLinkedList(t *testing.T, events []*event) []*event {
	for i := 0; i < len(events); i++ {
		event := events[i]

		if i == 0 {
			event.Prev = nil
		}
		if i == len(events)-1 {
			event.Next = nil
		}

		if i > 0 {
			prevEvent := events[i-1]
			prevEvent.Next = event
			event.Prev = prevEvent
		}

		if i < len(events)-1 {
			nextEvent := events[i+1]
			nextEvent.Prev = event
			event.Next = nextEvent
		}
	}
	checkEventsAreSortedLinkedList(t, events)

	return events
}

// nolint:staticcheck
func checkEventsLinkedListsEqual(t *testing.T, headA, headB *event) {
	idx := 0
	for ; ; headA, headB = headA.Next, headB.Next {
		if headA == nil && headB == nil {
			return
		}

		if headA == nil && headB != nil {
			require.True(t, false, "index[%d]: A -> nil, B -> not nil", idx)
		}
		if headA != nil && headB == nil {
			require.True(t, false, "index[%d]: A -> not nil, B -> nil", idx)
		}

		require.True(t, headA.Start.Equal(headB.Start), "index[%d]: Start", idx)
		require.True(t, headA.End.Equal(headB.End), "index[%d]: End", idx)

		idx++
	}
}
