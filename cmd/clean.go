package cmd

import (
	"fmt"
	"strconv"

	"tiktok-live-logger/pkg/database"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean [days]",
	Short: "Clean old logs",
	Long: `Remove logs older than the specified number of days.
If no days are specified, defaults to 30 days.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		days := 30
		if len(args) > 0 {
			var err error
			days, err = strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid number of days: %w", err)
			}
		}

		// Initialize database
		db, err := database.NewDB(GetDBPath())
		if err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer db.Close()

		// Delete old events
		rows, err := db.DeleteOldEvents(days)
		if err != nil {
			return fmt.Errorf("failed to delete old events: %w", err)
		}

		fmt.Printf("Deleted %d events older than %d days\n", rows, days)
		return nil
	},
} 