//go:build linux
// +build linux

package ps

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var (
	ErrNoProcess = errors.New("no process with provided name found")
)

// Process contains information about running process
type Process struct {
	// The process id
	Pid int
	// The executable filename of running process
	Comm string
	// state of this process
	State string
	// The parent process id from this process
	Ppid int
	//  The process group id from this process
	Pgrp int
}

// ProcessIsExists return true if given process name is exists
func ProcessIsExists(name string) (bool, error) {
	_, err := FindProcess(name)
	if err != nil {
		return false, err
	}

	return true, nil
}

// ProcessIsExists return true if given pid is exists
func PidIsExists(pid int) (bool, error) {
	_, err := FindProcessFromPid(pid)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetPidState returns state for given process id
func GetPidState(pid int) (string, error) {
	proc, err := FindProcessFromPid(pid)
	if err != nil {
		return "", err
	}

	return StateToString(proc.State), nil
}

// FindPpid returns parent process id from given process name
func FindPpid(name string) (int, error) {
	procs, err := FindProcess(name)
	if err != nil {
		return 0, err
	}

	return procs.Ppid, nil
}

// FindPid returns process id from given process name
func FindPid(name string) (int, error) {
	proc, err := FindProcess(name)
	if err != nil {
		return 0, err
	}
	return proc.Pid, nil
}

// FindPGid returns process group id from given process name
func FindPGid(name string) (int, error) {
	proc, err := FindProcess(name)
	if err != nil {
		return 0, err
	}

	return proc.Pgrp, nil
}

// FindProcessFromPid return Process struct from given process id
func FindProcessFromPid(pid int) (Process, error) {
	procs, err := GetProcess()
	if err != nil {
		return Process{}, err
	}

	for _, proc := range procs {
		if proc.Pid == pid {
			return proc, nil
		}
	}

	return Process{}, ErrNoProcess
}

// FindProcess return Process struct from given process name,
// it will be match if given process name same with executable filename
func FindProcess(name string) (Process, error) {
	procs, err := GetProcess()
	if err != nil {
		return Process{}, err
	}

	for _, proc := range procs {
		if proc.Comm == fmt.Sprintf("(%s)", name) {
			return proc, nil
		}
	}

	return Process{}, ErrNoProcess
}

// GetProcess returns all process information pid, comm, ppid, pgrp
func GetProcess() ([]Process, error) {
	var procs []Process

	pids, err := GetPids()
	if err != nil {
		return nil, err
	}

	for _, pid := range pids {
		stat, err := readStatsfile(pid)
		if err != nil {
			return nil, err
		}

		procs = append(procs, Process{
			Pid:   toInt(stat[0]),
			Comm:  stat[1],
			State: stat[2],
			Ppid:  toInt(stat[3]),
			Pgrp:  toInt(stat[4]),
		})
	}

	return procs, nil
}

// GetPids returns a slice of process ID
func GetPids() ([]int, error) {
	var pids []int

	dirs, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, dir := range dirs {
		pid, err := strconv.Atoi(dir.Name())
		if err != nil {
			continue
		}

		pids = append(pids, pid)
	}

	return pids, nil
}

func readStatsfile(pid int) ([]string, error) {
	fname := fmt.Sprintf("/proc/%d/stat", pid)

	f, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	field := strings.Split(string(f), " ")

	return field, nil
}

// StateToString returns state representation for given state
func StateToString(state string) string {
	// https://man7.org/linux/man-pages/man5/proc.5.html
	states := map[string]string{
		"R": "Running",
		"S": "Sleeping",
		"D": "Waiting",
		"Z": "Zombie",
		"T": "Stopped",
		"t": "Tracing stop",
		"X": "Dead",
		"x": "Dead",
		"K": "Wakekill",
		"W": "Waking",
		"P": "Parked",
	}

	return states[state]
}

func toInt(s string) int {
	i, _ := strconv.Atoi(s)

	return i
}
