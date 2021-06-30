package sqlite

import (
	"context"
	"fmt"

	"github.com/itiky/charge_scheduler/common"
	"github.com/itiky/charge_scheduler/schema"
)

func (s EventsStorage) CreateSingleEvent(ctx context.Context, obj schema.SingleEvent) (retId int64, retErr error) {
	dbObj, err := newSingleEvent(obj)
	if err != nil {
		retErr = fmt.Errorf("obj marshal: %v: %w", err, common.ErrInvalidInput)
		return
	}

	res, err := s.Db.NamedExecContext(ctx, "INSERT INTO single_events (type, start_date_time, end_hours, end_minutes, created_at) VALUES (:type, :start_date_time, :end_hours, :end_minutes, :created_at)", dbObj)
	if err != nil {
		retErr = fmt.Errorf("s.Db.NamedExecContext: %w", err)
		return
	}

	resId, err := res.LastInsertId()
	if err != nil {
		retErr = fmt.Errorf("res.LastInsertId(): %w", err)
		return
	}
	retId = resId

	return
}

func (s EventsStorage) CreatePeriodicEvent(ctx context.Context, obj schema.PeriodicEvent) (retId int64, retErr error) {
	dbObj, err := newPeriodicEvent(obj)
	if err != nil {
		retErr = fmt.Errorf("obj marshal: %v: %w", err, common.ErrInvalidInput)
		return
	}

	res, err := s.Db.NamedExecContext(ctx, "INSERT INTO periodic_events (type, rrule, end_hours, end_minutes, created_at) VALUES (:type, :rrule, :end_hours, :end_minutes, :created_at)", dbObj)
	if err != nil {
		retErr = fmt.Errorf("s.Db.NamedExecContext: %w", err)
		return
	}

	resId, err := res.LastInsertId()
	if err != nil {
		retErr = fmt.Errorf("res.LastInsertId(): %w", err)
		return
	}
	retId = resId

	return
}
