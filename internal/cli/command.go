package cli

import (
	"fmt"
)

type Command struct {
	Name string
	Args []string
}

type CmdHandler func(*State, Command) error

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

	// Users
	cmds.Register("login", handlerLogin)
	cmds.Register("register", handlerRegister)
	cmds.Register("reset", handlerReset)
	cmds.Register("users", handlerUsers)

	// Feeds
	cmds.Register("agg", requireLoggedIn(handlerAggregate))
	cmds.Register("addfeed", requireLoggedIn(handlerAddFeed))
	cmds.Register("feeds", requireLoggedIn(handlerFeeds))

	// Follows
	cmds.Register("follow", requireLoggedIn(handlerFollow))
	cmds.Register("following", requireLoggedIn(handlerFollowing))
	return cmds
}
