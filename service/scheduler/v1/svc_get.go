package v1

import (
	"context"
	"fmt"
	"time"

	"github.com/itiky/charge_scheduler/common"
	"github.com/itiky/charge_scheduler/schema"
)

const dayDur = 24 * time.Hour

func (svc Scheduler) GetAvailableAgenda(ctx context.Context, periodStart time.Time, periodDur, desiredDur time.Duration) (retAgendas schema.AgendaResults, retErr error) {
	// Input checks
	if periodStart.IsZero() {
		retErr = fmt.Errorf("%s: zero: %w", "periodStart", common.ErrInvalidInput)
		return
	}
	if periodDur <= 0 {
		retErr = fmt.Errorf("%s: must be GT 0: %w", "periodDur", common.ErrInvalidInput)
		return
	}
	if desiredDur <= 0 {
		retErr = fmt.Errorf("%s: must be GT 0: %w", "desiredDur", common.ErrInvalidInput)
		return
	}

	// Get existing events [-1 day : +periodDur +1 day]
	rangeStart, rangeEnd := periodStart.Add(-dayDur), periodStart.Add(periodDur).Add(dayDur)
	greenEvents, redEvents, err := svc.getGreenRedEvents(ctx, rangeStart, rangeEnd)
	if err != nil {
		retErr = fmt.Errorf("svc.getGreenRedEvents: %w", err)
		return
	}

	// Remove reds from greens and build the result
	greenHead := svc.mergeGreenRedEvents(greenEvents, redEvents)
	retAgendas = svc.buildAgendaResults(greenHead, periodStart, periodDur, desiredDur)

	return
}

// mergeGreenRedEvents alters greens by splitting / removing its elements (subtracting red elements).
func (svc Scheduler) mergeGreenRedEvents(greenEvents, redEvents []*event) *event {
	if len(greenEvents) == 0 {
		return nil
	}
	if len(redEvents) == 0 {
		return greenEvents[0]
	}

	greenHead := greenEvents[0]
	for _, redCur := range redEvents {
		for greenCur := greenHead; greenCur != nil; greenCur = greenCur.Next {
			// Optimization
			if greenCur.Start.After(redCur.End) {
				break
			}

			// Check if green and red intersects
			if svc.checkEventsIntersect(redCur, greenCur) {
				// Remove green element / exchange with splits
				greenCurSplits := svc.splitGreenEventWithRedEvent(greenCur, redCur)
				if len(greenCurSplits) == 0 {
					// Remove
					if greenCur.Prev != nil {
						greenCur.Prev.Next = greenCur.Next
						if greenCur.Next != nil {
							greenCur.Next.Prev = greenCur.Prev
						}
					} else {
						greenHead = greenCur.Next
						if greenHead != nil {
							greenHead.Prev = nil
						}
					}
				} else {
					// Exchange
					if greenCur.Prev != nil {
						greenCur.Prev.Next = greenCurSplits[0]
						greenCurSplits[0].Prev = greenCur.Prev
					} else {
						greenHead = greenCurSplits[0]
					}

					if greenCur.Next != nil {
						greenCurSplits[len(greenCurSplits)-1].Next = greenCur.Next
						greenCur.Next.Prev = greenCurSplits[len(greenCurSplits)-1]
					}

					greenCur = greenCurSplits[len(greenCurSplits)-1]
				}
			}
		}
	}

	return greenHead
}

// splitGreenEventWithRedEvent builds new green parts (if any) removing red intersection with green.
func (svc Scheduler) splitGreenEventWithRedEvent(green, red *event) (retGreenParts []*event) {
	if green.Start.Before(red.Start) {
		newGreenPart := &event{
			Id:    green.Id,
			Type:  green.Type,
			Start: green.Start,
			End:   red.Start,
			Prev:  nil,
			Next:  nil,
		}

		retGreenParts = append(retGreenParts, newGreenPart)
	}

	if green.End.After(red.End) {
		newGreenPart := &event{
			Id:    green.Id,
			Type:  green.Type,
			Start: red.End,
			End:   green.End,
			Prev:  nil,
			Next:  nil,
		}

		if len(retGreenParts) == 1 {
			retGreenParts[0].Next = newGreenPart
			newGreenPart.Prev = retGreenParts[0]
		}
		retGreenParts = append(retGreenParts, newGreenPart)
	}

	return
}

// buildAgendaResults builds schema.AgendaResult list searching for available time slots within desired duration.
func (svc Scheduler) buildAgendaResults(greenHead *event, periodStart time.Time, periodDur, desiredDur time.Duration) (retAgendas schema.AgendaResults) {
	periodEnd := periodStart.Add(periodDur)

	// removeTime return time.Time containing only date
	removeTime := func(ts time.Time) time.Time {
		return time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location())
	}

	// addEmptyAgendas adds to the results empty agendas within range [start, end)
	addEmptyAgendas := func(start, end time.Time, startInclusive, endInclusive bool) {
		start, end = removeTime(start), removeTime(end)

		if !startInclusive {
			start = start.Add(dayDur)
		}
		if endInclusive {
			end = end.Add(dayDur)
		}

		if start.Equal(end) || start.After(end) {
			return
		}

		for ; !start.Equal(end); {
			retAgendas = append(retAgendas, schema.AgendaResult{
				Date:      start,
				TimeSlots: nil,
			})
			start = start.Add(dayDur)
		}
	}

	// Edge case
	if greenHead == nil {
		addEmptyAgendas(periodStart, periodEnd, true, false)
		return
	}

	// Prefill
	addEmptyAgendas(periodStart, greenHead.Start, true, false)

	greenLast := greenHead
	for greenCur := greenHead; greenCur != nil; greenCur = greenCur.Next {
		// Middle fill
		addEmptyAgendas(greenLast.Start, greenCur.Start, false, false)
		greenLast = greenCur

		// Check if prev agenda should be used (instead of the new one)
		var agenda schema.AgendaResult
		createAgenda := true
		if len(retAgendas) > 0 {
			agenda = retAgendas[len(retAgendas)-1]
			if removeTime(greenCur.Start).Equal(agenda.Date) {
				createAgenda = false
				retAgendas = retAgendas[:len(retAgendas)-1]
			}
		}
		if createAgenda {
			agenda = schema.AgendaResult{
				Date: removeTime(greenCur.Start),
			}
		}

		// Check if time chunk is big enough
		eventDur := greenCur.End.Sub(greenCur.Start)
		if eventDur < desiredDur {
			retAgendas = append(retAgendas, agenda)
			continue
		}

		// Get time slots
		curTs := greenCur.Start
		for i := 0; i < int(eventDur/desiredDur); i++ {
			agenda.TimeSlots = append(agenda.TimeSlots, schema.TimeSlot{
				Start:    curTs,
				Duration: desiredDur,
			})
			curTs = curTs.Add(desiredDur)
		}
		retAgendas = append(retAgendas, agenda)
	}

	// Postfill
	addEmptyAgendas(greenLast.Start, periodEnd, false, false)

	return
}

func printList(greenHead *event) {
	fmt.Println("\nGreen list:")
	for greenCur := greenHead; greenCur != nil; greenCur = greenCur.Next {
		fmt.Printf("%s -> %s\n", greenCur.Start.Format(common.TimeFmt), greenCur.End.Format(common.TimeFmt))
	}
}
