package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/siaal/gator/internal/database"
)

func handlerFollow(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("follow requires 1 argument {feed name}, got %+v", cmd.Args)
	}
	feedURL := cmd.Args[0]
	userName := s.Config.CurrentUsername
	now := time.Now().UTC()
	args := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedUrl:   feedURL,
		Username:  userName,
		UpdatedAt: now,
		CreatedAt: now,
	}
	ctx := context.Background()
	feed, err := s.DB.CreateFeedFollow(ctx, args)
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}
	fmt.Println(feed)
	return nil
}

func handlerFollowing(s *State, cmd Command) error {
	ctx := context.Background()
	following, err := s.DB.GetFollowing(ctx, s.Config.CurrentUsername)
	if err != nil {
		return fmt.Errorf("fetch err: %w", err)
	}
	for _, followed := range following {
		fmt.Printf("* %s\n", followed.Name)
	}
	return nil
}
