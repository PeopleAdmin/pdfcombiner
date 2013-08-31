// Package cpdf contains methods to manipulate pdf files.
package cpdf

import (
	"fmt"
	"github.com/PeopleAdmin/pdfcombiner/testmode"
	"os/exec"
	"strconv"
	"strings"
)

// Panic at import time if cpdf not found.
var cpdf = cmdPath()

// Merge combines the files located at the specified paths into a
// single pdf with a custom header.
func Merge(doclist []string, outfile, title string) (err error) {
	if testmode.IsEnabled() {
		return
	}
	err = execMerge(doclist, outfile)
	if err != nil {
		return
	}
	err = addHeader(outfile, title)
	return
}

func PageCount(filePath string) (result int) {
	if testmode.IsEnabled() {
		return
	}
	cmd := exec.Command(cpdf, "-pages", filePath)
	out, err := cmd.Output()
	if err != nil {
		return -1
	}
	trimmed := strings.Trim(string(out), " \n")
	result, _ = strconv.Atoi(trimmed)
	return
}

func execMerge(doclist []string, outfile string) (err error) {
	cmd := exec.Command(cpdf)
	cmd.Args = cpdfMergeArgs(doclist, outfile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%v - %s", err, out)
	}
	return
}

func addHeader(filePath, title string) (err error) {
	headerText := "Page %Page of %EndPage | Created %m-%d-%Y %T | " + title
	cmd := exec.Command(cpdf, "-add-text", headerText,
		"-top", "15", "-font", "Courier", filePath, "-o", filePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%v - %s", err, out)
	}
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
