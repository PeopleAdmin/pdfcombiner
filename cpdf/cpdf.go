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

type PdfInfoCmd interface {
	PageCount() (int, error)
	ListBookmarks() string
}

type PdfManipulatorCmd interface {
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

// Merge concatenates the given files and adds a header in one pass.
func (c *Cpdf) Merge(docPaths []string, title string) (err error) {
	c.addMergeArgs(docPaths)
	c.addArgs("AND")
	c.addHeaderArgs(title)
	c.addArgs("-o", c.File)
	_, err = c.run()
	return
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

func (c *Cpdf) addMergeArgs(docPaths []string) {
	c.addArgs("-merge")
	for _, doc := range docPaths {
		c.addArgs(doc)
	}
}

func (c *Cpdf) addHeaderArgs(title string) {
	headerText := "Page %Page of %EndPage | Created %m-%d-%Y %T | " + title
	c.addArgs("-add-text", headerText, "-top", "15", "-font", "Courier")
}

func (c *Cpdf) addArgs(newArgs ...string) {
	c.command.Args = append(c.command.Args, newArgs...)
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

func cpdfPath() string {
	pathToCmd, err := exec.LookPath("cpdf")
	if err != nil {
		panic("no cpdf found in path")
	}
	return pathToCmd
}
