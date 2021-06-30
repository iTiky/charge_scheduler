package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/itiky/charge_scheduler/schema"
)

func (s EventsStorage) GetSingleEvent(ctx context.Context, id int64) (retObj *schema.SingleEvent, retErr error) {
	dbObj := singleEvent{}
	err := s.Db.GetContext(ctx, &dbObj, "SELECT rowid, type, start_date_time, end_hours, end_minutes, created_at FROM single_events WHERE rowid=?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
		retErr = fmt.Errorf("s.Db.GetContext: %w", err)
		return
	}

	obj, err := dbObj.ToSchema()
	if err != nil {
		retErr = fmt.Errorf("obj unmarshal: %w", err)
		return
	}
	retObj = &obj

	return
}

func (s EventsStorage) GetSingleEventsWithinRange(ctx context.Context, rangeStart, rangeEnd time.Time) (retObjs []schema.SingleEvent, retErr error) {
	var dbObjs []singleEvent
	err := s.Db.SelectContext(ctx, &dbObjs, "SELECT rowid, type, start_date_time, end_hours, end_minutes, created_at FROM single_events WHERE start_date_time >= ? AND start_date_time <= ?", rangeStart, rangeEnd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
		retErr = fmt.Errorf("s.Db.SelectContext: %w", err)
		return
	}

	objs, err := s.unmarshalSingleEvents(dbObjs)
	if err != nil {
		retErr = err
		return
	}

	return objs, nil
}

func (s EventsStorage) GetPeriodicEvent(ctx context.Context, id int64) (retObj *schema.PeriodicEvent, retErr error) {
	dbObj := periodicEvent{}
	err := s.Db.GetContext(ctx, &dbObj, "SELECT rowid, type, rrule, end_hours, end_minutes, created_at FROM periodic_events WHERE rowid=?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
		retErr = fmt.Errorf("s.Db.GetContext: %w", err)
		return
	}

	obj, err := dbObj.ToSchema()
	if err != nil {
		retErr = fmt.Errorf("obj unmarshal: %w", err)
		return
	}
	retObj = &obj

	return
}

func (s EventsStorage) GetAllPeriodicEvents(ctx context.Context) (retObjs []schema.PeriodicEvent, retErr error) {
	var dbObjs []periodicEvent
	err := s.Db.SelectContext(ctx, &dbObjs, "SELECT rowid, type, rrule, end_hours, end_minutes, created_at FROM periodic_events")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
		retErr = fmt.Errorf("s.Db.SelectContext: %w", err)
		return
	}

	objs, err := s.unmarshalPeriodicEvents(dbObjs)
	if err != nil {
		retErr = err
		return
	}

	return objs, nil
}

func (s EventsStorage) unmarshalSingleEvents(dbObjs []singleEvent) (retObjs []schema.SingleEvent, retErr error) {
	retObjs = make([]schema.SingleEvent, 0, len(dbObjs))
	for i, dbObj := range dbObjs {
		obj, err := dbObj.ToSchema()
		if err != nil {
			retErr = fmt.Errorf("dbObj[%d] unmarshal: %w", i, err)
			return
		}
		retObjs = append(retObjs, obj)
	}

	return
}

func (s EventsStorage) unmarshalPeriodicEvents(dbObjs []periodicEvent) (retObjs []schema.PeriodicEvent, retErr error) {
	retObjs = make([]schema.PeriodicEvent, 0, len(dbObjs))
	for i, dbObj := range dbObjs {
		obj, err := dbObj.ToSchema()
		if err != nil {
			retErr = fmt.Errorf("dbObj[%d] unmarshal: %w", i, err)
			return
		}
		retObjs = append(retObjs, obj)
	}

	return
}
