package cli

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/siaal/gator/internal/database"
)

type Command struct {
	Name string
	Args []string
}

type CmdHandler func(*State, Command) error

func handlerUsers(s *State, _ Command) error {
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
func handlerReset(s *State, cmd Command) error {
	err := s.DB.ClearUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Cleared users")
	return nil
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("login takes 1 argument {username}, got: %d", len(cmd.Args))
	}
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

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("register takes 1 argument {username}, got: %d", len(cmd.Args))
	}
	username := cmd.Args[0]
	now := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
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

type Commands struct {
	Cmds map[string]CmdHandler
}

func (c *Commands) Register(name string, f CmdHandler) error {
	c.Cmds[name] = f
	return nil
}

func (c *Commands) Run(s *State, cmd Command) error {
	callback, ok := c.Cmds[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown cmd: %s", cmd.Name)
	}
	return callback(s, cmd)
}

func DefaultCommands() Commands {
	cmds := Commands{Cmds: map[string]CmdHandler{}}
	cmds.Register("login", handlerLogin)
	cmds.Register("register", handlerRegister)
	cmds.Register("reset", handlerReset)
	cmds.Register("users", handlerUsers)
	return cmds
}
