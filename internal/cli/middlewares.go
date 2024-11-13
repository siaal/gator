package cli

import "fmt"

func requireLoggedIn(f CmdHandler) CmdHandler {
	return func(s *State, c Command) error {
		if s.Config.CurrentUsername == "" {
			return fmt.Errorf("%s requires you to be logged in", c.Name)
		}
		return f(s, c)
	}
}
