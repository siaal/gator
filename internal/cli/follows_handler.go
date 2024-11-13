package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/siaal/gator/internal/database"
	"github.com/siaal/gator/internal/state"
)

func handlerFollow(s *state.State, cmd Command) error {
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

func handlerFollowing(s *state.State, cmd Command) error {
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

func handlerUnfollow(s *state.State, cmd Command) error {
	ctx := context.Background()
	if err := s.DB.Unfollow(ctx, database.UnfollowParams{Username: s.Config.CurrentUsername, FeedUrl: cmd.Args[0]}); err != nil {
		return fmt.Errorf("could not unfollow %s, %w", cmd.Args[0], err)
	}
	return nil
}
