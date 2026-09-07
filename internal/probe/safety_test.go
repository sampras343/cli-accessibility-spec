// internal/probe/safety_test.go
package probe

import "testing"

func TestIsBlockedSubcommand(t *testing.T) {
	blocked := []string{"delete", "rm", "remove", "destroy", "purge", "DROP", "Format", "init", "reset", "clean", "wipe", "nuke", "uninstall"}
	for _, s := range blocked {
		if !IsBlockedSubcommand(s) {
			t.Errorf("%q should be blocked", s)
		}
	}

	allowed := []string{"list", "get", "show", "status", "help", "version", "config", "describe"}
	for _, s := range allowed {
		if IsBlockedSubcommand(s) {
			t.Errorf("%q should be allowed", s)
		}
	}
}
