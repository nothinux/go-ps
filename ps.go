//go:build linux
// +build linux

package ps

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrNoProcess = errors.New("no process with provided name found")
	statRX       = regexp.MustCompile(`(\w.+):\s+(.+)`)
)

// Process contains information about running process
type Process struct {
	// The process id
	Pid int
	// The executable filename of running process
	Comm string
	// Full command line of running process
	CmdLine string
	// state of this process
	State string
	// The parent process id from this process
	Ppid int
	//  The process group id from this process
	Pgrp int
	// The User id and group id who running process
	UID, GID int
}

// ProcessIsExists return true if given process name is exists
func ProcessIsExists(name string) (bool, error) {
	_, err := FindProcessName(name)
	if err != nil {
		return false, err
	}

	return true, nil
}

// ProcessIsExists return true if given pid is exists
func PidIsExists(pid int) (bool, error) {
	_, err := FindProcess(pid)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetPidState returns state for given process id
func GetPidState(pid int) (string, error) {
	proc, err := FindProcess(pid)
	if err != nil {
		return "", err
	}

	return StateToString(proc.State), nil
}

// FindPpid returns parent process id from given process name
func FindPpid(name string) (int, error) {
	procs, err := FindProcessName(name)
	if err != nil {
		return 0, err
	}

	return procs.Ppid, nil
}

// FindPid returns process id from given process name
func FindPid(name string) (int, error) {
	proc, err := FindProcessName(name)
	if err != nil {
		return 0, err
	}
	return proc.Pid, nil
}

// FindPGid returns process group id from given process name
func FindPGid(name string) (int, error) {
	proc, err := FindProcessName(name)
	if err != nil {
		return 0, err
	}

	return proc.Pgrp, nil
}

// FindProcess returns Process struct from given process id
func FindProcess(pid int) (Process, error) {
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

// FindProcessName returns Process based on finding,
// it will be match if given process name same with executable filename
func FindProcessName(name string) (Process, error) {
	procs, err := GetProcess()
	if err != nil {
		return Process{}, err
	}

	for _, proc := range procs {
		if proc.Comm == name {
			return proc, nil
		}
	}

	return Process{}, ErrNoProcess
}

// FindProcessNameContains returns slice of Process based on finding,
// it will returns process containing name.
// Using this is not very accurate, if you want to search process based on
// executable file name, use FindProcessName instead
func FindProcessNameContains(name string) ([]Process, error) {
	var p []Process

	procs, err := GetProcess()
	if err != nil {
		return nil, err
	}

	for _, proc := range procs {
		if strings.Contains(proc.CmdLine, name) {
			p = append(p, proc)
		}
	}

	return p, nil
}

// GetProcess returns all process information pid, comm, ppid, pgrp
func GetProcess() ([]Process, error) {
	var procs []Process

	pids, err := GetPids()
	if err != nil {
		return nil, err
	}

	for _, pid := range pids {
		stat, err := readStatusfile(pid)
		if err != nil {
			return nil, err
		}

		// read full command line for this process
		cmdline, err := readCmdLine(pid)
		if err != nil {
			return nil, err
		}

		// add cmdline to map
		stat["CmdLine"] = cmdline

		p := parseStatusFile(stat)

		procs = append(procs, p)
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

func readStatusfile(pid int) (map[string]string, error) {
	fname := fmt.Sprintf("/proc/%d/status", pid)

	f, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	field := strings.Split(string(f), "\n")
	var status = map[string]string{}
	// format
	for _, f := range field {
		if statRX.Match([]byte(f)) {
			st := statRX.FindStringSubmatch(f)
			status[st[1]] = st[2]
		}
	}

	return status, nil
}

func readCmdLine(pid int) (string, error) {
	name := fmt.Sprintf("/proc/%d/cmdline", pid)

	cmd, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}

	return string(cmd), nil
}

// parseStatus file returns all information in Process
func parseStatusFile(status map[string]string) Process {
	var p Process

	// TODO parse all information from /proc/$/status file
	for k, v := range status {
		switch {
		case k == "Name":
			p.Comm = v
		case k == "State":
			s := strings.Split(v, " ")
			p.State = s[0]
		case k == "Pid":
			p.Pid = toInt(v)
		case k == "PPid":
			p.Ppid = toInt(v)
		case k == "Tgid":
			p.Pgrp = toInt(v)
		case k == "Uid":
			p.UID = toInt(v)
		case k == "Gid":
			p.GID = toInt(v)
		case k == "CmdLine":
			p.CmdLine = v
		default:
		}
	}

	return p
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

func prettyName(s string) string {
	return strings.Trim(s, "()")
}

func toInt(s string) int {
	i, _ := strconv.Atoi(s)

	return i
}
