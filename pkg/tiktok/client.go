package tiktok

import (
	"fmt"
	"time"

	"tiktok-live-logger/pkg/logger"

	"github.com/Davincible/gotiktoklive"
	"github.com/pkg/errors"
)

type Client struct {
	tiktok *gotiktoklive.TikTok
	logger *logger.Logger
}

func NewClient(debugMode bool) (*Client, error) {
	logger, err := logger.NewLogger(debugMode)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return &Client{
		tiktok: gotiktoklive.NewTikTok(),
		logger: logger,
	}, nil
}

type LiveStats struct {
	ViewerCount  int64
	LikeCount    int64
	ShareCount   int64
	CommentCount int64
}

type Event struct {
	Type      string
	Content   string
	Timestamp time.Time
}

type EventHandler func(Event)

func (c *Client) TrackUser(username string, onEvent EventHandler) error {
	c.logger.Info("Starting to track user: %s", username)

	live, err := c.tiktok.TrackUser(username)
	if err != nil {
		c.logger.ErrorWithStack(err, "Failed to track user: %s", username)
		return errors.Wrap(err, "failed to track user")
	}

	// Handle events from the channel
	go func() {
		for event := range live.Events {
			switch e := event.(type) {
			case gotiktoklive.ChatEvent:
				event := Event{
					Type:      "chat",
					Content:   fmt.Sprintf("%s: %s", e.User.Nickname, e.Comment),
					Timestamp: time.Unix(e.Timestamp, 0),
				}
				c.logger.Debug("Chat message: %s", event.Content)
				onEvent(event)

			case gotiktoklive.GiftEvent:
				event := Event{
					Type:      "gift",
					Content:   fmt.Sprintf("%s sent %s (x%d)", e.User.Nickname, e.Name, e.RepeatCount),
					Timestamp: time.Unix(e.Timestamp, 0),
				}
				c.logger.Info("Gift received: %s", event.Content)
				onEvent(event)

			case gotiktoklive.LikeEvent:
				event := Event{
					Type:      "like",
					Content:   fmt.Sprintf("%s sent %d likes", e.User.Nickname, e.Likes),
					Timestamp: time.Now(),
				}
				c.logger.Debug("Likes received: %s", event.Content)
				onEvent(event)

			case gotiktoklive.UserEvent:
				switch e.Event {
				case gotiktoklive.USER_FOLLOW:
					event := Event{
						Type:      "follow",
						Content:   fmt.Sprintf("%s followed the streamer", e.User.Nickname),
						Timestamp: time.Now(),
					}
					c.logger.Info("New follower: %s", event.Content)
					onEvent(event)

				case gotiktoklive.USER_SHARE:
					event := Event{
						Type:      "share",
						Content:   fmt.Sprintf("%s shared the stream", e.User.Nickname),
						Timestamp: time.Now(),
					}
					c.logger.Info("Stream shared: %s", event.Content)
					onEvent(event)
				}

			case gotiktoklive.ViewersEvent:
				event := Event{
					Type:      "stats",
					Content:   fmt.Sprintf("Viewer count: %d", e.Viewers),
					Timestamp: time.Now(),
				}
				c.logger.Debug("Room stats updated: %s", event.Content)
				onEvent(event)
			}
		}
	}()

	c.logger.Info("Successfully started tracking user: %s", username)
	return nil
}

func (c *Client) GetLiveStats(username string) (*LiveStats, error) {
	c.logger.Info("Getting live stats for user: %s", username)

	// Get room info directly from TikTok instance
	roomInfo, err := c.tiktok.GetRoomInfo(username)
	if err != nil {
		c.logger.ErrorWithStack(err, "Failed to get room info for user: %s", username)
		return nil, errors.Wrap(err, "failed to get room info")
	}

	stats := &LiveStats{
		ViewerCount:  int64(roomInfo.Stats.TotalUser),
		LikeCount:    int64(roomInfo.Stats.LikeCount),
		ShareCount:   int64(roomInfo.Stats.ShareCount),
		CommentCount: int64(roomInfo.Stats.DiggCount),
	}

	c.logger.Debug("Live stats for %s: viewers=%d, likes=%d, shares=%d, comments=%d",
		username, stats.ViewerCount, stats.LikeCount, stats.ShareCount, stats.CommentCount)

	return stats, nil
}

func (c *Client) Close() error {
	if c.logger != nil {
		return c.logger.Close()
	}
	return nil
} 