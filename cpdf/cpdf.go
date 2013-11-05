// Package cpdf contains methods to manipulate pdf files.
package cpdf

import (
	"fmt"
	"github.com/PeopleAdmin/pdfcombiner/testmode"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

// Panic at import time if cpdf not found.
var cpdfbin = cpdfPath()

type InfoCmd interface {
	ListBookmarks() ([]byte, error)
	PageCount() (int, error)
	Validate() error
}

type ManipulatorCmd interface {
	Merge([]string, string) error
	AddBookmarks(string)
}

type Cpdf struct {
	File    string
	command *exec.Cmd
	id      string
}

func New(filePath string, id string) (c *Cpdf) {
	c = &Cpdf{File: filePath}
	c.command = exec.Command(cpdfbin)
	c.id = id
	return
}

func (c *Cpdf) ListBookmarks() (out []byte, err error) {
	c.setArgs("-list-bookmarks", c.File)
	return c.run()
}

// Validate checks whether the file can be processed by cpdf.  In particular,
// it should catch encrypted documents before they get to the combine phase.
func (c *Cpdf) Validate() (err error) {
	c.setArgs("-i", c.File)
	_, err = c.run()
	return
}

// PageCount returns the number of pages in the document.
func (c *Cpdf) PageCount() (result int, err error) {
	if testmode.IsEnabled() {
		return 1, nil
	}

	c.setArgs("-pages", c.File)
	out, err := c.run()
	if err != nil {
		return
	}
	trimmed := strings.Trim(string(out), " \n")
	result, err = strconv.Atoi(trimmed)
	return
}

func (c *Cpdf) run() (output []byte, err error) {
	log.Println(c.id, `CPDF: "`+strings.Join(c.command.Args, " ")+`"`)
	if testmode.IsEnabled() {
		return
	}
	output, err = c.command.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%v - %s", err, output)
	}
	return
}

// Clear any existing args before setting the new ones.  Useful if you're
// reusing an existing cpdf object.
func (c *Cpdf) setArgs(newArgs ...string) {
	c.command = exec.Command(cpdfbin)
	c.addArgs(newArgs...)
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

// Intended for use on a blank file used to indicate an error to the user.
func (c *Cpdf) WriteErrorMessage(title string) {
	message := fmt.Sprintf("There was an error including the file: '%s'.", title)
	c.addArgs("-add-text", message, "-topleft", "50", "-font-size", "16", c.File,
		"-o", c.File)
	c.run()
}
