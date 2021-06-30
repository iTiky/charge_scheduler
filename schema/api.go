package schema

import (
	"fmt"
	"strings"
	"time"

	"github.com/itiky/charge_scheduler/common"
)

type (
	AgendaResult struct {
		Date      time.Time
		TimeSlots []TimeSlot
	}

	TimeSlot struct {
		Start    time.Time
		Duration time.Duration
	}

	AgendaResults []AgendaResult
)

func (r AgendaResults) String() string {
	str := strings.Builder{}
	for _, item := range r {
		str.WriteString(item.String())
	}

	return str.String()
}

func (r AgendaResult) String() string {
	str := strings.Builder{}
	str.WriteString("Agenda:\n")
	str.WriteString(fmt.Sprintf("  Date: %s\n", r.Date.Format("02.01.2006")))
	if len(r.TimeSlots) == 0 {
		str.WriteString("  Slots: none\n")
	} else {
		str.WriteString("  Slots:\n")
		for _, slot := range r.TimeSlots {
			str.WriteString(fmt.Sprintf("  - %s -> %s\n", slot.Start.Format(common.TimeFmt), slot.Duration))
		}
	}

	return str.String()
}
