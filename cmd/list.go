package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"tiktok-live-logger/pkg/database"
	"tiktok-live-logger/pkg/ui"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all logged live streams",
	Long: `View a list of all logged live streams and their events.
Select a stream to view its detailed logs.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize database
		dbPath := filepath.Join(os.Getenv("HOME"), ".tiktok-live-logger", "events.db")
		db, err := database.NewDB(dbPath)
		if err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer db.Close()

		// Get all usernames
		usernames, err := db.GetAllUsernames()
		if err != nil {
			return fmt.Errorf("failed to get usernames: %w", err)
		}

		// Create list items
		var items []list.Item
		for _, username := range usernames {
			events, err := db.GetEventsByUsername(username)
			if err != nil {
				continue
			}

			if len(events) > 0 {
				// Group events by date
				dates := make(map[string]time.Time)
				for _, event := range events {
					date := event.Timestamp.Format("2006-01-02")
					if _, exists := dates[date]; !exists {
						dates[date] = event.Timestamp
					}
				}

				// Add username and dates to list
				items = append(items, listItem{
					title:    username,
					desc:     fmt.Sprintf("%d sessions", len(dates)),
					username: username,
					dates:    dates,
				})
			}
		}

		// Initialize UI
		model := ui.NewModel("")
		model.SetListItems(items)
		model.SetViewMode(true)

		// Start the program
		p := tea.NewProgram(model, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("failed to run UI: %w", err)
		}

		return nil
	},
}

type listItem struct {
	title    string
	desc     string
	username string
	dates    map[string]time.Time
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.desc }
func (i listItem) FilterValue() string { return i.title } 