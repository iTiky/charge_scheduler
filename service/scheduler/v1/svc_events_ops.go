package v1

func (svc Scheduler) checkEventsIntersect(e1, e2 *event) bool {
	earlierEvent, laterEvent := e1, e2
	if laterEvent.Start.Before(earlierEvent.Start) {
		earlierEvent, laterEvent = laterEvent, earlierEvent
	}

	if laterEvent.Start.Equal(earlierEvent.Start) || laterEvent.Start.Equal(earlierEvent.End) {
		return true
	}
	if laterEvent.Start.After(earlierEvent.Start) && laterEvent.Start.Before(earlierEvent.End) {
		return true
	}

	return false
}
