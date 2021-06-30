package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

const (
	FlagChargeDur = "charge-duration"
)

// GetAgendaCmd returns get agenda command.
func GetAgendaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "agenda [periodStartDateTime] [periodDur]",
		Short:   "Get available charging slots for a specified period and charging time",
		Example: "agenda 2020-02-21T12:00:00Z 240h 30m",
		Long: `Arguments:
  [periodStartDateTime] - period start dateTime (RFC 3339);
  [periodDur] - requested period duration;
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

			periodDur, err := time.ParseDuration(args[1])
			if err != nil {
				logger.Fatal().Str("arg", "periodDur").Err(err).Msg("invalid")
			}

			chargingDur, err := cmd.Flags().GetDuration(FlagChargeDur)
			if err != nil {
				logger.Fatal().Str("flag", FlagChargeDur).Err(err).Msg("invalid")
			}

			// Init dependencies and request
			svc := getService(logger, cmd)
			agenda, err := svc.GetAvailableAgenda(context.TODO(), periodStart, periodDur, chargingDur)
			if err != nil {
				logger.Fatal().Err(err).Msg("svc.GetAvailableAgenda")
			}

			// Print response
			if len(agenda) == 0 {
				logger.Fatal().Msg("agenda is empty")
			}
			fmt.Print(agenda.String())
		},
	}
	cmd.Flags().Duration(FlagChargeDur, 30*time.Minute, "(optional) desired charging duration")

	return cmd
}

func init() {
	rootCmd.AddCommand(GetAgendaCmd())
}
