// Package cpdf contains methods to manipulate pdf files.
package cpdf

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Panic at import time if cpdf not found.
var cpdf = cmdPath()

// Merge combines the files located at the specified paths into a single pdf.
func Merge(doclist []string, outfile string) (err error) {
	combine_cmd := exec.Command(cpdf)
	combine_cmd.Args = cpdfMergeArgs(doclist, outfile)
	out, failed := combine_cmd.CombinedOutput()
	if failed != nil {
		err = fmt.Errorf("%v - %s", failed, out)
	}
	return
}

func PageCount(filePath string) (result int) {
	cmd := exec.Command(cpdf, "-pages", filePath)
	out, err := cmd.Output()
	if err != nil {
		return -1
	}
	trimmed := strings.Trim(string(out), " \n")
	result, _ = strconv.Atoi(trimmed)
	return
}

// cmdPath returns the absolute path to the cpdf executable.
func cmdPath() string {
	pathToCmd, err := exec.LookPath("cpdf")
	if err != nil {
		panic("no cpdf found in path")
	}
	return pathToCmd
}

// Given a slice of documents, construct an args slice suitable for
// passing to 'cpdf merge'.
func cpdfMergeArgs(doclist []string, outfile string) (args []string) {
	args = []string{"merge"}
	for _, doc := range doclist {
		args = append(args, doc)
	}
	args = append(args, "-o", outfile)
	return
}
