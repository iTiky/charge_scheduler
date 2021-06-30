package schema

import (
	"fmt"
	"strings"
	"time"

	"github.com/teambition/rrule-go"

	"github.com/itiky/charge_scheduler/common"
)

type (
	SingleEvent struct {
		Id            int64           `json:"id"`
		Type          SingleEventType `json:"type"`
		StartDateTime time.Time       `json:"start_date_time"`
		EndHours      uint            `json:"end_hours"`
		EndMinutes    uint            `json:"end_minutes"`
		CreatedAt     time.Time       `json:"created_at"`
	}

	SingleEventType string
)

const (
	SingleEventTypeAvailable SingleEventType = "Available"
	SingleEventTypeOccupied  SingleEventType = "Occupied"
)

func (t SingleEventType) IsValid() bool {
	switch t {
	case SingleEventTypeAvailable, SingleEventTypeOccupied:
		return true
	default:
		return false
	}
}

func (t SingleEventType) String() string {
	return string(t)
}

func (e SingleEvent) String() string {
	str := strings.Builder{}
	str.WriteString("SingleEvent:\n")
	str.WriteString(fmt.Sprintf("  Id: %d\n", e.Id))
	str.WriteString(fmt.Sprintf("  Type: %s\n", e.Type.String()))
	str.WriteString(fmt.Sprintf("  Start: %s\n", e.StartDateTime.Format(common.TimeFmt)))
	str.WriteString(fmt.Sprintf("  End: %02d:%02d\n", e.EndHours, e.EndMinutes))
	str.WriteString(fmt.Sprintf("  CreatedAt: %s\n", e.CreatedAt.Format(common.TimeFmt)))

	return str.String()
}

type PeriodicEvent struct {
	Id         int64           `json:"id"`
	Type       SingleEventType `json:"type"`
	Rrule      rrule.RRule     `json:"rrule"`
	EndHours   uint            `json:"end_hours"`
	EndMinutes uint            `json:"end_minutes"`
	CreatedAt  time.Time       `json:"created_at"`
}

func (e PeriodicEvent) String() string {
	str := strings.Builder{}
	str.WriteString("PeriodicEvent:\n")
	str.WriteString(fmt.Sprintf("  Id: %d\n", e.Id))
	str.WriteString(fmt.Sprintf("  Type: %s\n", e.Type.String()))
	str.WriteString(fmt.Sprintf("  RRule: %s\n", e.Rrule.String()))
	str.WriteString(fmt.Sprintf("  End: %02d:%02d\n", e.EndHours, e.EndMinutes))
	str.WriteString(fmt.Sprintf("  CreatedAt: %s\n", e.CreatedAt.Format(common.TimeFmt)))

	return str.String()
}
