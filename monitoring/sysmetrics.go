// Functions to return general system metrics like disk space and free memory.
package monitoring

import (
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const DF_AVAIL_COL = 3
const MEMFILE = "/proc/meminfo"

var spaceSeperator = regexp.MustCompile("[\t ]+")

// The free space on the root partition.
func RootPartitionFreeMb() func() float64 {
	return func() float64 {
		return float64(MbFreeAtMount("/"))
	}
}

// System free memory.
func FreeMemoryMb() func() float64 {
	return func() float64 {
		return float64(MemStat("MemFree"))
	}
}

// Parses `df -m` for the given mount and returns the MB free in that dir.
func MbFreeAtMount(mountpoint string) int {
	out, _ := exec.Command("df", "-m", mountpoint).Output()
	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return 0
	}
	cols := spaceSeperator.Split(lines[1], -1)
	if len(cols) < DF_AVAIL_COL+1 {
		return 0
	}

	mbs, err := strconv.Atoi(cols[DF_AVAIL_COL])
	if err != nil {
		return 0
	}
	return mbs
}

// Parses /proc/meminfo and returns the requested line as MBs.
func MemStat(statname string) int {
	content, err := ioutil.ReadFile(MEMFILE)
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
