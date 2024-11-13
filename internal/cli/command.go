package cli

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/siaal/gator/internal/database"
	"github.com/siaal/gator/rss"
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

func handlerFeeds(s *State, cmd Command) error {
	feeds, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("fetch err: %w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("FEED\n* Name: %s\n* URL : %s\n* User: %s\n\n", feed.Name, feed.Url, feed.UserName)
	}
	return nil
}
func handlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("addfeed requires 2 arguments: {name} {url}. Got: %+v", cmd.Args)
	}
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

func handlerAggregate(s *State, cmd Command) error {
	url := "https://www.wagslane.dev/index.xml"
	ctx := context.Background()
	feed, err := rss.FetchFeed(ctx, url)
	if err != nil {
		return fmt.Errorf("feed fetch err: %w", err)
	}
	fmt.Println(feed)
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
func DefaultCommands() Commands {
	cmds := Commands{Cmds: map[string]CmdHandler{}}
	cmds.Register("login", handlerLogin)
	cmds.Register("register", handlerRegister)
	cmds.Register("reset", handlerReset)
	cmds.Register("users", handlerUsers)
	cmds.Register("agg", handlerAggregate)
	cmds.Register("addfeed", handlerAddFeed)
	cmds.Register("feeds", handlerFeeds)
	cmds.Register("follow", handlerFollow)
	cmds.Register("following", handlerFollowing)
	return cmds
}
