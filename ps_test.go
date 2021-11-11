package ps

import (
	"strings"
	"testing"
)

func TestGetProcess(t *testing.T) {
	_, err := GetProcess()
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindPpid(t *testing.T) {
	t.Run("Get parent pid of available process", func(t *testing.T) {
		pid, err := FindPpid("systemd")
		if err != nil {
			t.Fatal(err)
		}

		if pid != 0 {
			t.Fatalf("expected pid is 0, got %v", pid)
		}
	})

	t.Run("Get parent pid of unavailable process", func(t *testing.T) {
		_, err := FindPpid("anythingxxx")
		if err != nil {
			if !strings.Contains(err.Error(), "no process with provided name found") {
				t.Fatal(err)
			}
		}
	})
}
