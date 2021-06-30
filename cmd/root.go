package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/itiky/charge_scheduler/service/scheduler"
	v1 "github.com/itiky/charge_scheduler/service/scheduler/v1"
	"github.com/itiky/charge_scheduler/storage/events/sqlite"
	"github.com/itiky/charge_scheduler/storage/sqlite_base"
)

const (
	FlagLogLevel = "log-level"
	FlagDbPath   = "db-path"
)

// rootCmd is a base command.
var rootCmd = &cobra.Command{
	Use:   "charge-scheduler",
	Short: "Charge scheduler app",
}

func getLogger(cmd *cobra.Command) (zerolog.Logger, error) {
	logLevelRaw, err := cmd.Flags().GetString(FlagLogLevel)
	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("reading %s flag: %w", FlagLogLevel, err)
	}

	logLevel, err := zerolog.ParseLevel(logLevelRaw)
	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("parsing %s flag: %w", FlagLogLevel, err)
	}

	return zerolog.New(os.Stderr).
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		Level(logLevel).
		With().
		Timestamp().
		Logger(), nil
}

func getService(logger zerolog.Logger, cmd *cobra.Command) scheduler.Scheduler {
	dbPath, err := cmd.Flags().GetString(FlagDbPath)
	if err != nil {
		logger.Fatal().Str("flag", FlagDbPath).Err(err).Msg("reading")
	}

	baseSt, err := sqlite_base.NewSQLiteBase(logger, dbPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("baseStorage init")
	}

	if err := baseSt.Migrate(); err != nil {
		logger.Fatal().Err(err).Msg("baseStorage migration")
	}

	eventsSt, err := sqlite.NewEventsStorage(baseSt)
	if err != nil {
		logger.Fatal().Err(err).Msg("eventsStorage init")
	}

	svc, err := v1.NewScheduler(logger, eventsSt)
	if err != nil {
		logger.Fatal().Err(err).Msg("schedulerService init")
	}

	return svc
}

func main() {
	rootCmd.PersistentFlags().String(FlagLogLevel, "debug", "Logging level")
	rootCmd.PersistentFlags().String(FlagDbPath, "./sqlite.db", "Path to SQLite3 database")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("rootCmd.Execute: %v", err)
	}
}
