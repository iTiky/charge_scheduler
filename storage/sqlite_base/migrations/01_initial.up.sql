CREATE TABLE single_events
(
    type            TEXT      NOT NULL,
    start_date_time TIMESTAMP NOT NULL,
    end_hours       INTEGER   NOT NULL,
    end_minutes     INTEGER   NOT NULL,
    created_at      TIMESTAMP NOT NULL
);

CREATE TABLE periodic_events
(
    type        TEXT      NOT NULL,
    rrule       TEXT      NOT NULL,
    end_hours   INTEGER   NOT NULL,
    end_minutes INTEGER   NOT NULL,
    created_at  TIMESTAMP NOT NULL
);