package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Config struct {
	DefaultDaysToKeep int    `json:"default_days_to_keep"`
	DatabasePath      string `json:"database_path"`
	DebugMode         bool   `json:"debug_mode"`
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long: `View and modify the application configuration.
The configuration file is stored in ~/.tiktok-live-logger/config.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir := filepath.Join(os.Getenv("HOME"), ".tiktok-live-logger")
		configPath := filepath.Join(configDir, "config.json")

		// Create config directory if it doesn't exist
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		// Load existing config or create default
		config := &Config{
			DefaultDaysToKeep: 30,
			DatabasePath:      filepath.Join(configDir, "events.db"),
			DebugMode:         false,
		}

		if data, err := os.ReadFile(configPath); err == nil {
			if err := json.Unmarshal(data, config); err != nil {
				return fmt.Errorf("failed to parse config file: %w", err)
			}
		}

		// If no subcommand is provided, show current config
		if len(args) == 0 {
			data, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal config: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		// Handle subcommands
		switch args[0] {
		case "set":
			if len(args) < 3 {
				return fmt.Errorf("usage: config set <key> <value>")
			}
			key := args[1]
			value := args[2]

			switch key {
			case "default_days_to_keep":
				var days int
				if _, err := fmt.Sscanf(value, "%d", &days); err != nil {
					return fmt.Errorf("invalid number of days: %w", err)
				}
				config.DefaultDaysToKeep = days
			case "database_path":
				config.DatabasePath = value
			case "debug_mode":
				config.DebugMode = value == "true"
			default:
				return fmt.Errorf("unknown config key: %s", key)
			}

			// Save updated config
			data, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal config: %w", err)
			}
			if err := os.WriteFile(configPath, data, 0644); err != nil {
				return fmt.Errorf("failed to write config file: %w", err)
			}
			fmt.Println("Configuration updated successfully")
		default:
			return fmt.Errorf("unknown subcommand: %s", args[0])
		}

		return nil
	},
} 