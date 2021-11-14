package ps

import (
	"fmt"
	"testing"
)

// All tests in here only work on Environment with systemd as init

func TestShowAllProcess(t *testing.T) {
	procs, err := GetProcess()
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range procs {
		t.Logf("procs name: %v \tprocs cmd: %v\n", p.Comm, p.CmdLine)
	}
}

func TestProcessIsExists(t *testing.T) {
	exists, err := ProcessIsExists("systemd")
	if err != nil {
		t.Fatal(err)
	}

	// systemd must exists
	if !exists {
		t.Fatalf("expected systemd is exists, but got not exists")
	}
}

func TestPidIsExists(t *testing.T) {
	exists, err := PidIsExists(1)
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatalf("expected pid 1 is exists, but got not exists")
	}
}

func TestGetPidState(t *testing.T) {
	state, err := GetPidState(1)
	if err != nil {
		t.Fatal(err)
	}

	if state != "Sleeping" {
		t.Fatalf("expected state Sleeping, got %v", state)
	}
}

func TestFindPid(t *testing.T) {
	t.Run("Get pid of available process", func(t *testing.T) {
		pid, err := FindPid("systemd")
		if err != nil {
			t.Fatal(err)
		}

		if pid != 1 {
			t.Fatalf("expected pid is 1, got %v", pid)
		}
	})

	t.Run("Get pid of unavailable process", func(t *testing.T) {
		_, err := FindPid("anythingxxx")
		if err != nil {
			if err != ErrNoProcess {
				t.Fatal(err)
			}
		}
	})
}

func TestFindPpid(t *testing.T) {
	t.Run("Get parent pid of available process", func(t *testing.T) {
		ppid, err := FindPpid("systemd")
		if err != nil {
			t.Fatal(err)
		}

		if ppid != 0 {
			t.Fatalf("expected parent pid is 0, got %v", ppid)
		}
	})

	t.Run("Get parent pid of unavailable process", func(t *testing.T) {
		_, err := FindPpid("anythingxxx")
		if err != nil {
			if err != ErrNoProcess {
				t.Fatal(err)
			}
		}
	})
}

func TestFindPGid(t *testing.T) {
	t.Run("Get process group id of available process", func(t *testing.T) {
		pgid, err := FindPGid("systemd")
		if err != nil {
			t.Fatal(err)
		}

		if pgid != 1 {
			t.Fatalf("expected pid is 1, got %v", pgid)
		}
	})

	t.Run("Get pid of unavailable process", func(t *testing.T) {
		_, err := FindPGid("anythingxxx")
		if err != nil {
			if err != ErrNoProcess {
				t.Fatal(err)
			}
		}
	})
}

func TestFindProcess(t *testing.T) {
	p, err := FindProcess(1)
	if err != nil {
		t.Fatal(err)
	}

	if p.Comm != "systemd" {
		t.Fatalf("expected comm systemd, got %v", p.Comm)
	}

	if p.Pid != 1 {
		t.Fatalf("expected pid 1, got %v", p.Pid)
	}

	if p.State != "S" {
		t.Fatalf("expected state S, got %v", p.State)
	}
}

func TestFindProcessName(t *testing.T) {
	p, err := FindProcessName("systemd")
	if err != nil {
		t.Fatal(err)
	}

	if p.Comm != "systemd" {
		t.Fatalf("expected comm systemd, got %v", p.Comm)
	}

	if p.Pid != 1 {
		t.Fatalf("expected pid 1, got %v", p.Pid)
	}

	if p.State != "S" {
		t.Fatalf("expected state S, got %v", p.State)
	}
}

func TestProcessNameTree(t *testing.T) {
	ps, err := FindProcessNameTree("systemd")
	if err != nil {
		t.Fatal(err)
	}

	// make sure last element have expected parent id by compare it with pid
	// from parent process
	if ps[len(ps)-1].Ppid != ps[0].Pid {
		t.Fatalf("child parent id not same with parent process id")
	}
}

func TestGetProcess(t *testing.T) {
	procs, err := GetProcess()
	if err != nil {
		t.Fatal(err)
	}

	// assume index 0 is pid 1
	if procs[0].Comm != "systemd" {
		t.Fatalf("expected comm systemd, got %v", procs[0].Comm)
	}
}

func TestGetPids(t *testing.T) {
	pids, err := GetPids()
	if err != nil {
		t.Fatal(err)
	}

	// check pid 1 is exists
	if pids[0] != 1 {
		t.Fatalf("no pid 1 found")
	}
}

func TestStateToString(t *testing.T) {
	tests := []struct {
		state  string
		rState string
	}{
		{
			state:  "S",
			rState: "Sleeping",
		},
		{
			state:  "R",
			rState: "Running",
		},
		{
			state:  "D",
			rState: "Waiting",
		},
		{
			state:  "Z",
			rState: "Zombie",
		},
		{
			state:  "T",
			rState: "Stopped",
		},
		{
			state:  "t",
			rState: "Tracing stop",
		},
		{
			state:  "X",
			rState: "Dead",
		},
		{
			state:  "x",
			rState: "Dead",
		},
		{
			state:  "K",
			rState: "Wakekill",
		},
		{
			state:  "W",
			rState: "Waking",
		},
		{
			state:  "P",
			rState: "Parked",
		},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("Test state %s representation", tt.state)

		t.Run(name, func(t *testing.T) {
			state := StateToString(tt.state)

			if tt.rState != state {
				t.Fatalf("expected %v, got %v", tt.rState, state)
			}
		})
	}
}
