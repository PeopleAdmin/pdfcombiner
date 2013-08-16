// Package cpdf contains methods to manipulate pdf files.
package cpdf

import (
	"fmt"
	"github.com/kless/shutil/sh"
	"log"
	"os/exec"
	"strings"
	"time"
)

// Combine the files located at the specified paths into a single pdf.
func Merge(doclist []string) (outfile string, err error) {
	outfile = getOutfile()
	cmd := mergeCmd(doclist, outfile)
	out, err := sh.Run(cmd)
	log.Println(cmd)
	log.Println(out)
	log.Println(err)
	return
}

// Constructs the command string to pdftk -- looks like
// "/bin/pdftk file1.pdf file2.pdf output 12345.pdf"
func mergeCmd(doclist []string, outfile string) string {
	cmdComponents := []string{
		CmdPath(),
		strings.Join(appendTmp(doclist), " "),
		("output /tmp/" + outfile)}
	return strings.Join(cmdComponents, " ")
}

func CmdPath() string {
	pathToCmd, err := exec.LookPath("pdftk")
	if err != nil {
		log.Fatal("no pdftk found in path")
	}
	return pathToCmd
}

func getOutfile() string {
	return fmt.Sprintf("%d.pdf", time.Now().Unix())
}

func appendTmp(list []string) []string {
	for idx, file := range list {
		list[idx] = "/tmp/" + file
	}
	return list
}

// Given a slice of documents, construct an args slice suitable for
// passing to 'cpdf merge'.
func cpdfMergeArgs(doclist []string) (args []string) {
	args = []string{"merge"}
	for _, doc := range doclist {
		args = append(args, doc)
	}
	args = append(args, "-o", "combined.pdf")
	return
}
