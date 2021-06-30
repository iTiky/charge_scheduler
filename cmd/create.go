package main

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/itiky/charge_scheduler/schema"
)

const (
	FlagWeekly    = "weekly"
)

// CreateSingleEventCmd returns create schema.SingleEvent object command.
func CreateSingleEventCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [scheduleType] [eventStartDateTime] [eventEndTime]",
		Short:   "Create a schedule event (single / recurrent) of a specified type",
		Example: "create Available 2020-02-21T12:00:00Z 15:30 --weekly",
		Long: `Arguments:
  [scheduleType] - schedule type (Available / Occupied);
  [eventStartDateTime] - event start dateTime (RFC 3339);
  [eventEndTime] - event end time during the {eventStartDateTime} day (HH:MM);
`,
		Args: cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			logger, err := getLogger(cmd)
			if err != nil {
				log.Fatal(err)
			}

			// Parse inputs
			eventType := schema.SingleEventType(args[0])
			if !eventType.IsValid() {
				logger.Fatal().Str("arg", "scheduleType").Msg("invalid")
			}

			eventStart, err := time.Parse(time.RFC3339, args[1])
			if err != nil {
				logger.Fatal().Str("arg", "eventStartDateTime").Err(err).Msg("invalid")
			}

			eventEndTime, err := time.Parse("15:04", args[2])
			if err != nil {
				logger.Fatal().Str("arg", "eventEndTime").Err(err).Msg("invalid")
			}

			isWeekly, err := cmd.Flags().GetBool(FlagWeekly)
			if err != nil {
				logger.Fatal().Str("flag", FlagWeekly).Err(err).Msg("invalid")
			}

			// Init dependencies and request
			svc := getService(logger, cmd)
			if isWeekly {
				if err := svc.AddPeriodicEvent(context.TODO(), eventType, eventStart, uint(eventEndTime.Hour()), uint(eventEndTime.Minute())); err != nil {
					logger.Fatal().Err(err).Msg("svc.AddPeriodicEvent")
				}
			} else {
				if err := svc.AddSingleEvent(context.TODO(), eventType, eventStart, uint(eventEndTime.Hour()), uint(eventEndTime.Minute())); err != nil {
					logger.Fatal().Err(err).Msg("svc.AddSingleEvent")
				}
			}
		},
	}
	cmd.Flags().Bool(FlagWeekly, false, "(optional) recurrent schedule event type")

	return cmd
}

func init() {
	rootCmd.AddCommand(CreateSingleEventCmd())
}
