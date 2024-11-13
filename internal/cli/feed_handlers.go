package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/siaal/gator/internal/database"
	"github.com/siaal/gator/internal/state"
	"github.com/siaal/gator/rss"
)

func handlerFeeds(s *state.State, cmd Command) error {
	feeds, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("fetch err: %w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("FEED\n* Name: %s\n* URL : %s\n* User: %s\n\n", feed.Name, feed.Url, feed.UserName)
	}
	return nil
}
func handlerAddFeed(s *state.State, cmd Command) error {
	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	ctx := context.Background()
	now := time.Now().UTC()
	addFeedParam := database.AddFeedParams{
		ID:        uuid.New(),
		Name:      feedName,
		Url:       feedURL,
		CreatedAt: now,
		UpdatedAt: now,
		Username:  s.Config.CurrentUsername,
	}
	feed, err := s.DB.AddFeed(ctx, addFeedParam)
	if err != nil {
		return fmt.Errorf("failed to add feed: %w", err)
	}
	fmt.Println(feed)
	ctx = context.Background()
	ffArgs := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UpdatedAt: now,
		CreatedAt: now,
		FeedUrl:   feedURL,
		Username:  s.Config.CurrentUsername,
	}
	ff, err := s.DB.CreateFeedFollow(ctx, ffArgs)
	if err != nil {
		return fmt.Errorf("Successfully added feed, but failed to follow: %w", err)
	}
	fmt.Println(ff)
	return nil
}

func handlerAggregate(s *state.State, cmd Command) error {
	interval := 1 * time.Hour
	if len(cmd.Args) == 1 {
		int, err := time.ParseDuration(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("Failed to parse arg as time interval: %s %w", cmd.Args[0], err)
		}
		interval = int
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		if err := rss.ScrapeFeeds(s); err != nil {
			return err
		}
	}
}
