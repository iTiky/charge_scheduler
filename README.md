# Charge scheduler

## Task

The goal is to write an algorithm that checks the availabilities of an agenda depending of the events attached to it.
The main method has a start date for input and is looking for the availabilities of the next `10` days.

They are two kinds of events:
- `available`, are the availability of the Charge Point for a specific day and they can be recurring week by week.
- `occupied`, times when the Charge Point is already booked.

Requirements:
- Tests with edge case covering.
- Be pragmatic about performance.
- Code must be SQLite compatible.

Basic unit test vision:
- Create data base records with fields for Event
    - kind: `available` , starts_at: `2014-08-04 09:30`, ends_at: `2014-08-04 13:30`, weekly_recurring: `true`
    - kind: `occupied` , starts_at: `2014-08-11 10:30` , ends_at: `2014-08-11 11:30`
- Get Event availabilities in data base started from date `2014-08-10`
- Compare Event availabilities with:
    - `availabilities[0]["date"]` should be equal to should be equal to `2014-08-10`
    - `availabilities[0]["slots"]` should be equal to `[]`
    - `availabilities[1]["date"]` should be equal to `2014-08-11`
    - `availabilities[1]["slots"]` should be equal to `["9:30", "10:00", "11:30","12:00", "12:30", "13:00"]`
    - `availabilities[2]["slots"]` should be equal to `[]`
    - `availabilities[6]["date"]` should be equal to `2014-08-16`
    - Event availabilities length should be equal to `10`

## Build and test

Requirements:
* At least `Golang` 1.13 compiler installed;
* `gcc` compiler (required for `SQLite3 driver`);

**Running tests**
```Bash
make test
```
Command runs unit and integration tests (including the one described in the task).

**Build**
```Bash
make install
```
Default binary file path is `./build/charge-scheduler`.

## CLI

Application has a CLI interface with build-in help and examples:
```Bash
./charge-scheduler -h
./charge-scheduler create -h
./charge-scheduler list -h
./charge-scheduler agenda -h
```

**Example**
```Bash
# Create the "Available" recurring calendar event
./charge-scheduler create Available 2014-08-04T09:30:00Z 13:30 --weekly
# Create the "Occupied" single calendar event
./charge-scheduler create Occupied 2014-08-11T10:30:00Z 11:30

# Print all the registered events so far
./charge-scheduler list 2014-08-04T00:00:00Z 2014-08-15T23:59:00Z

# Request available charging slots within 10days and 30min charging duration
./charge-scheduler agenda 2014-08-10T00:00:00Z 240h
```

Output:
```Bash
Agenda:
  Date: 10.08.2014
  Slots: none
Agenda:
  Date: 11.08.2014
  Slots:
  - 11.08.2014 09:30:00 UTC -> 30m0s
  - 11.08.2014 10:00:00 UTC -> 30m0s
  - 11.08.2014 11:30:00 UTC -> 30m0s
  - 11.08.2014 12:00:00 UTC -> 30m0s
  - 11.08.2014 12:30:00 UTC -> 30m0s
  - 11.08.2014 13:00:00 UTC -> 30m0s
Agenda:
  Date: 12.08.2014
  Slots: none
Agenda:
  Date: 13.08.2014
  Slots: none
Agenda:
  Date: 14.08.2014
  Slots: none
Agenda:
  Date: 15.08.2014
  Slots: none
Agenda:
  Date: 16.08.2014
  Slots: none
Agenda:
  Date: 17.08.2014
  Slots: none
Agenda:
  Date: 18.08.2014
  Slots:
  - 18.08.2014 09:30:00 UTC -> 30m0s
  - 18.08.2014 10:00:00 UTC -> 30m0s
  - 18.08.2014 10:30:00 UTC -> 30m0s
  - 18.08.2014 11:00:00 UTC -> 30m0s
  - 18.08.2014 11:30:00 UTC -> 30m0s
  - 18.08.2014 12:00:00 UTC -> 30m0s
  - 18.08.2014 12:30:00 UTC -> 30m0s
  - 18.08.2014 13:00:00 UTC -> 30m0s
Agenda:
  Date: 19.08.2014
  Slots: none
```

## Design

**Storage**
`SQLite 3` in-memory database is used to persist calendar events.

*Single* calendar events are stored within `single_events` tables with the following schema:
```SQL
CREATE TABLE single_events
(
    type            TEXT      NOT NULL,
    start_date_time TIMESTAMP NOT NULL,
    end_hours       INTEGER   NOT NULL,
    end_minutes     INTEGER   NOT NULL,
    created_at      TIMESTAMP NOT NULL
);
```

Event is defined with start timestamp and end hours and minutes. `HH:MM` approach limits the event duration to a single day (`00:00 - 23:59`).

*Recurring* (periodic) calendar events are stored within `periodic_events` table with the following schema:
```SQL
CREATE TABLE periodic_events
(
    type        TEXT      NOT NULL,
    rrule       TEXT      NOT NULL,
    end_hours   INTEGER   NOT NULL,
    end_minutes INTEGER   NOT NULL,
    created_at  TIMESTAMP NOT NULL
);
```

Event repeat pattern is serialized using Apple iCalendar RRule (RFC 5545). Few points regarding this decision:
* We do not reinvent formats;
* RRule allows usage of more complex (comparing to *weekly*) patterns;
* Avoid coding dateTime algos;
* Ability to add *exception* cases like:
    * remove a single event keeping the pattern generated ones;
    * alter a single event keeping the rest intact;
    

Database migrations are embedded to the application binary.


Storage layer has its own data models (schema) with serialize / deserialize functions for bi-directional service - storage models conversion.
This approach makes service layer models a bit more "usable" by services and API handlers.

**Algorithm flow**

Here is a brief algo description for agenda requests.
1. Get all registered periodic events.
2. Get all registered single events filtered by input defined time range.
3. Using the RRule engine generate periodic events as single events withing input defined time range.
4. Sort and split all the events into two groups:
    * Green events: *Available*;
    * Red events: *Occupied*;
    * Groups are structured as a sorted double linked list;
5. Iterate over Reds looking for a Green events it intersects to and alter the Greens list.
    * Green might be removed from Greens if Red if "bigger";
    * Green might shrink (partial Red-Green intersection);
    * Green might be splitted into few smaller events and inserted to Greens (Red was "in the middle" of Green);
6. An altered (or completely erased) Greens lists is aggregated to Days.
7. Time slots are searched within a day considering the desired charging duration (30 mins by default);

## Errors

* Input checks are performed along the way (from API to Storage) to avoid wrong input failures;
* User can't create an event which has intersections with already existing events (for the same event type: Green / Red);

## Implementation limitations and points of improvement

1. Unsafe for concurrent event creation requests
    * Parallel events overlapping checks do not know about each other;
    * POI: lock the DB on *create* requests;
    * POI: add requests queue for "single create at a time" approach;
2. Reread and reprocessing of all periodic events for each *agenda* request
    * POI: add a cache layer which stores "unrolled" RRule events for the current and upcoming months;
3. All-in-one app
    * No server running making each request start the DB;
    * POI: add an API level (RPC, gRPC, REST,...);
    
## Dependencies

Project has a `/vendor` folder which makes it "GitHub repo removal" proof.
* `github.com/spf13/cobra`
    * industry standard library for building a CLI interface;
* `github.com/stretchr/testify`
    * should be part the stdlib;
    * adds TestCase feature and makes testing fun;
* `github.com/rs/zerolog`
    * structured logging library;
* `github.com/mattn/go-sqlite3`
    * SQLite DB driver;
* `github.com/jmoiron/sqlx`
    * extends the standard `database/sql` interface with some helper functions;
* `github.com/golang-migrate/migrate/v4`
    * DB migration tool;
* `github.com/teambition/rrule-go`
    * `python-dateutil` library port;
    * adds RRule build and validation options;
    * adds RRule processing;