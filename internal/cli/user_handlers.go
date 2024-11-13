package cli

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/siaal/gator/internal/database"
	"github.com/siaal/gator/internal/state"
)

func handlerUsers(s *state.State, _ Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("db err: %w", err)
	}
	for _, u := range users {
		if u.Name == s.Config.CurrentUsername {
			fmt.Printf("* %s (current)\n", u.Name)
		} else {
			fmt.Printf("* %s\n", u.Name)
		}
	}
	return nil
}
func handlerReset(s *state.State, cmd Command) error {
	err := s.DB.ClearUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Cleared users")
	return nil
}

func handlerLogin(s *state.State, cmd Command) error {
	username := cmd.Args[0]
	user, err := s.DB.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}
	slog.Debug("Logged in as user", "user", user)
	if err = s.Config.SetUser(user.Name); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	fmt.Printf("User changed to: %s\n", username)
	return nil
}

func handlerRegister(s *state.State, cmd Command) error {
	username := cmd.Args[0]
	now := time.Now().UTC()
	userParams := database.CreateUserParams{

		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      username,
	}
	user, err := s.DB.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	if err = s.Config.SetUser(user.Name); err != nil {
		return fmt.Errorf("user created successful, however switch user failed: %w", err)

	}
	fmt.Println("User created: " + user.Name)
	slog.Debug("Created user", "user", user)
	return nil
}
