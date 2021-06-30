package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

// ListEventsCmd returns list events command.
func ListEventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [periodStartDateTime] [periodEndDateTime]",
		Short:   "Print registered events within specified time range",
		Example: "agenda 2020-02-21T12:00:00Z 72h 30m",
		Long: `Arguments:
  [periodStartDateTime] - period start dateTime (RFC 3339);
  [periodEndDateTime] - period end dateTime (RFC 3339);
`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			logger, err := getLogger(cmd)
			if err != nil {
				log.Fatal(err)
			}

			// Parse inputs
			periodStart, err := time.Parse(time.RFC3339, args[0])
			if err != nil {
				logger.Fatal().Str("arg", "periodStartDateTime").Err(err).Msg("invalid")
			}

			periodEnd, err := time.Parse(time.RFC3339, args[1])
			if err != nil {
				logger.Fatal().Str("arg", "periodEndDateTime").Err(err).Msg("invalid")
			}

			// Init dependencies and request
			svc := getService(logger, cmd)
			sEvents, pEvents, err := svc.GetEvents(context.TODO(), periodStart, periodEnd)
			if err != nil {
				logger.Fatal().Err(err).Msg("svc.GetEvents")
			}

			// Print response
			for _, event := range sEvents {
				fmt.Print(event.String())
			}
			for _, event := range pEvents {
				fmt.Print(event.String())
			}
		},
	}

	return cmd
}

func init() {
	rootCmd.AddCommand(ListEventsCmd())
}
