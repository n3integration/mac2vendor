package actions

import (
	"sync"

	"gopkg.in/urfave/cli.v1"
)

var registry = struct {
	actions []cli.Command
	mu      sync.Mutex
}{
	actions: make([]cli.Command, 0),
}

// GetCommands returns the supported application commands
func GetCommands() []cli.Command {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	return registry.actions
}

func register(cmd cli.Command) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.actions = append(registry.actions, cmd)
}
