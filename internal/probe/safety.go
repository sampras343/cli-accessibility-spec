// internal/probe/safety.go
package probe

import "strings"

var blockedSubcommands = []string{
	"delete", "rm", "remove", "destroy", "purge", "drop",
	"format", "init", "reset", "clean", "wipe", "nuke",
	"truncate", "uninstall", "erase",
}

func IsBlockedSubcommand(name string) bool {
	lower := strings.ToLower(name)
	for _, blocked := range blockedSubcommands {
		if lower == blocked {
			return true
		}
	}
	return false
}
