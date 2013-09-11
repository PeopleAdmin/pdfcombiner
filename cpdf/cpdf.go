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
var cpdfbin = cpdfPath()

type InfoCmd interface {
	ListBookmarks() ([]byte, error)
	PageCount() (int, error)
}

type ManipulatorCmd interface {
	Merge([]string, string) error
	AddBookmarks(string)
}

type Cpdf struct {
	File    string
	command *exec.Cmd
}

func New(filePath string) (c *Cpdf) {
	c = &Cpdf{File: filePath}
	c.command = exec.Command(cpdfbin)
	return
}

func (c *Cpdf) ListBookmarks() (out []byte, err error) {
	c.addArgs("-list-bookmarks", c.File)
	return c.run()
}

// PageCount returns the number of pages in the document.
func (c *Cpdf) PageCount() (result int, err error) {
	c.command.Args = []string{"-pages", c.File}
	out, err := c.run()
	if err != nil {
		return
	}
	trimmed := strings.Trim(string(out), " \n")
	result, err = strconv.Atoi(trimmed)
	return
}

func (c *Cpdf) run() (output []byte, err error) {
	if testmode.IsEnabled() {
		return
	}
	output, err = c.command.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%v - %s", err, output)
	}
	return
}

func (c *Cpdf) addArgs(newArgs ...string) {
	c.command.Args = append(c.command.Args, newArgs...)
}

func cpdfPath() string {
	pathToCmd, err := exec.LookPath("cpdf")
	if err != nil {
		panic("no cpdf found in path")
	}
	return pathToCmd
}
