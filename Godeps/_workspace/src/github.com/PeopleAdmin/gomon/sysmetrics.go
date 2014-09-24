package gomon

import (
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Functions to return general system metrics like disk space and free memory.

const availableSpaceCol = 3
const linuxMemfile = "/proc/meminfo"

var spaceSeperator = regexp.MustCompile("[\t ]+")

// RootPartitionFree returns the space in MB available on '/'.
func RootPartitionFree() float64 {
	return float64(MbFreeAtMount("/"))
}

// MemoryFree returns the number of MB reported as 'available' by the unix
// `free` command.
func MemoryFree() float64 {
	return float64(memStat("MemFree"))
}

// MbFreeAtMount parses `df -m` for the given mount and returns the MB free in
// that directory.
func MbFreeAtMount(mountpoint string) int {
	out, _ := exec.Command("df", "-m", mountpoint).Output()
	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return 0
	}
	cols := spaceSeperator.Split(lines[1], -1)
	if len(cols) < availableSpaceCol+1 {
		return 0
	}

	mbs, err := strconv.Atoi(cols[availableSpaceCol])
	if err != nil {
		return 0
	}
	return mbs
}

// Parses /proc/meminfo and returns the requested line as MBs.
func memStat(statname string) int {
	content, err := ioutil.ReadFile(linuxMemfile)
	if err != nil {
		return 0
	}
	lines := strings.Split(string(content), "\n")
	if len(lines) < 2 {
		return 0
	}
	for _, line := range lines {
		cols := spaceSeperator.Split(line, -1)
		if cols[0] == statname+":" {
			kbs, err := strconv.Atoi(cols[1])
			if err != nil {
				return 0
			}
			return kbs / 1024
		}
	}
	return 0
}
