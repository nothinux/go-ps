package ps

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Process struct {
	Pid  int
	Comm string
	Ppid int
	Pgrp int
}

// FindPpid return parent process id of provided process name
func FindPpid(name string) (int, error) {
	procs, err := GetProcess()
	if err != nil {
		return 0, err
	}

	for _, proc := range procs {
		if proc.Comm == fmt.Sprintf("(%s)", name) {
			// retuen parrent process id
			return proc.Ppid, nil
		}
	}

	return 0, fmt.Errorf("no process with provided name found")
}

// Findpid return process id of provided process name
func Findpid(name string) (int, error) {
	procs, err := GetProcess()
	if err != nil {
		return 0, err
	}

	for _, proc := range procs {
		if proc.Comm == fmt.Sprintf("(%s)", name) {
			// retuen parrent process id
			return proc.Pid, nil
		}
	}

	return 0, fmt.Errorf("no process with provided name found")
}

// GetProcess get all process information pid, comm, ppid, pgrp
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
			Pid:  toInt(stat[0]),
			Comm: stat[1],
			Ppid: toInt(stat[3]),
			Pgrp: toInt(stat[4]),
		})
	}

	return procs, nil
}

// GetPids get all pid from /proc file
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

func toInt(s string) int {
	i, _ := strconv.Atoi(s)

	return i
}
