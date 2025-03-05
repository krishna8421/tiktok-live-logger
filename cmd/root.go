package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	commit  = "dev"
	date    = "05-03-2025"
)

var rootCmd = &cobra.Command{
	Use:   "tiktok-live-logger",
	Short: "A TikTok live stream logger with beautiful TUI",
	Long: `A command line tool to log and view TikTok live streams.
It connects to live streams, logs all events (chat, gifts, etc.) and provides
a beautiful interface to view the logs.`,
	Version: fmt.Sprintf("%s (commit: %s, date: %s)", version, commit, date),
}

func init() {
	// Add commands
	rootCmd.AddCommand(logCmd)
	// rootCmd.AddCommand(listCmd)
	// rootCmd.AddCommand(cleanCmd)
	// rootCmd.AddCommand(configCmd)

	// Add flags
	rootCmd.PersistentFlags().StringP("db", "d", "", "Path to database file")
	rootCmd.PersistentFlags().BoolP("debug", "v", false, "Enable debug mode")
}

func Execute() error {
	return rootCmd.Execute()
}

// GetDBPath returns the path to the database file
func GetDBPath() string {
	dbPath, err := rootCmd.Flags().GetString("db")
	if err == nil && dbPath != "" {
		return dbPath
	}
	return filepath.Join(os.Getenv("HOME"), ".tiktok-live-logger", "events.db")
}

// IsDebug returns whether debug mode is enabled
func IsDebug() bool {
	debug, _ := rootCmd.Flags().GetBool("debug")
	return debug
} 