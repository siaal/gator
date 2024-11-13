package cli

import (
	"fmt"

	"github.com/siaal/gator/internal/state"
)

type Command struct {
	Name string
	Args []string
}

type CmdHandler func(*state.State, Command) error

type Commands struct {
	Cmds map[string]CmdHandler
}

func (c *Commands) Register(name string, f CmdHandler) error {
	c.Cmds[name] = f
	return nil
}

func (c *Commands) Run(s *state.State, cmd Command) error {
	callback, ok := c.Cmds[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown cmd: %s", cmd.Name)
	}
	return callback(s, cmd)
}

func DefaultCommands() Commands {
	cmds := Commands{Cmds: map[string]CmdHandler{}}

	// Users
	cmds.Register("login", requireArgsNum(1, handlerLogin))
	cmds.Register("register", requireArgsNum(1, handlerRegister))
	cmds.Register("reset", handlerReset)
	cmds.Register("users", handlerUsers)

	// Feeds
	cmds.Register("agg", requireLoggedIn(handlerAggregate))
	cmds.Register("addfeed", requireLoggedIn(requireArgsNum(2, handlerAddFeed)))
	cmds.Register("feeds", requireLoggedIn(handlerFeeds))

	// Follows
	cmds.Register("follow", requireLoggedIn(requireArgsNum(1, handlerFollow)))
	cmds.Register("following", requireLoggedIn(handlerFollowing))
	cmds.Register("unfollow", requireLoggedIn(requireArgsNum(1, handlerUnfollow)))
	return cmds
}
