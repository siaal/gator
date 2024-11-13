package cli

import (
	"fmt"

	"github.com/siaal/gator/internal/state"
)

func requireLoggedIn(f CmdHandler) CmdHandler {
	return func(s *state.State, c Command) error {
		if s.Config.CurrentUsername == "" {
			return fmt.Errorf("%s requires you to be logged in", c.Name)
		}
		return f(s, c)
	}
}

func requireArgsNum(nArgs int, f CmdHandler) CmdHandler {
	return func(s *state.State, c Command) error {
		if len(c.Args) != nArgs {
			return fmt.Errorf("%s requires %d args. Got %d with %+v", c.Name, nArgs, len(c.Args), c.Args)
		}
		return f(s, c)
	}
}
