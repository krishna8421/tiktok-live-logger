package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"tiktok-live-logger/pkg/database"
	"tiktok-live-logger/pkg/tiktok"
	"tiktok-live-logger/pkg/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log [username]",
	Short: "Log a TikTok live stream",
	Long: `Connect to a TikTok live stream and log all events (chat, gifts, etc.)
to a SQLite database while displaying them in a beautiful TUI interface.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]

		// Initialize database
		dbPath := GetDBPath()
		if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
			return fmt.Errorf("failed to create database directory: %w", err)
		}

		db, err := database.NewDB(dbPath)
		if err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer db.Close()

		// Initialize TikTok client
		client, err := tiktok.NewClient(IsDebug())
		if err != nil {
			return fmt.Errorf("failed to initialize TikTok client: %w", err)
		}
		defer client.Close()

		// Initialize UI with username
		model := ui.NewModel(username)

		// Set up event handler
		onEvent := func(event tiktok.Event) {
			// Update UI with new event
			model.AddEvent(event.Content)
			
			// Save event to database
			if err := db.SaveEvent(event.Type, event.Content, event.Timestamp, username); err != nil {
				model.SetError(fmt.Errorf("failed to save event: %w", err))
			}
		}

		// Start tracking user
		if err := client.TrackUser(username, onEvent); err != nil {
			return fmt.Errorf("failed to track user: %w", err)
		}

		// Start the program
		p := tea.NewProgram(model, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("failed to run UI: %w", err)
		}

		return nil
	},
} 