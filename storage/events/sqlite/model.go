package sqlite

import (
	"fmt"
	"time"

	"github.com/teambition/rrule-go"

	"github.com/itiky/charge_scheduler/schema"
)

type singleEvent struct {
	Id            int64     `db:"rowid"`
	Type          string    `db:"type"`
	StartDateTime time.Time `db:"start_date_time"`
	EndHours      uint      `db:"end_hours"`
	EndMinutes    uint      `db:"end_minutes"`
	CreatedAt     time.Time `db:"created_at"`
}

func (e singleEvent) ToSchema() (schema.SingleEvent, error) {
	eType := schema.SingleEventType(e.Type)
	if !eType.IsValid() {
		return schema.SingleEvent{}, fmt.Errorf("%s: invalid", "type")
	}

	return schema.SingleEvent{
		Id:            e.Id,
		Type:          eType,
		StartDateTime: e.StartDateTime,
		EndHours:      e.EndHours,
		EndMinutes:    e.EndMinutes,
		CreatedAt:     e.CreatedAt,
	}, nil
}

func newSingleEvent(obj schema.SingleEvent) (singleEvent, error) {
	return singleEvent{
		Type:          obj.Type.String(),
		StartDateTime: obj.StartDateTime,
		EndHours:      obj.EndHours,
		EndMinutes:    obj.EndMinutes,
		CreatedAt:     obj.CreatedAt,
	}, nil
}

type periodicEvent struct {
	Id         int64     `db:"rowid"`
	Type       string    `db:"type"`
	Rrule      string    `db:"rrule"`
	EndHours   uint      `db:"end_hours"`
	EndMinutes uint      `db:"end_minutes"`
	CreatedAt  time.Time `db:"created_at"`
}

func (e periodicEvent) ToSchema() (schema.PeriodicEvent, error) {
	eType := schema.SingleEventType(e.Type)
	if !eType.IsValid() {
		return schema.PeriodicEvent{}, fmt.Errorf("%s: invalid", "type")
	}

	obj := schema.PeriodicEvent{
		Id:         e.Id,
		Type:       eType,
		EndHours:   e.EndHours,
		EndMinutes: e.EndMinutes,
		CreatedAt:  e.CreatedAt,
	}

	r, err := rrule.StrToRRule(e.Rrule)
	if err != nil {
		return schema.PeriodicEvent{}, fmt.Errorf("creating rrule (%s): %w", e.Rrule, err)
	}
	obj.Rrule = *r

	return obj, nil
}

func newPeriodicEvent(obj schema.PeriodicEvent) (periodicEvent, error) {
	return periodicEvent{
		Type:       obj.Type.String(),
		Rrule:      obj.Rrule.String(),
		EndHours:   obj.EndHours,
		EndMinutes: obj.EndMinutes,
		CreatedAt:  obj.CreatedAt,
	}, nil
}
