// Package cpdf contains methods to manipulate pdf files using cpdf.
// TODO this does not currently handle demo versions of cpdf, as a
// workaround, rename and wrap it with a script that strips demo text:
//   #!/usr/bin/env bash
//   _cpdf $@ |grep -v "This demo is for evaluation only"
package cpdf

import (
	"log"
	"os/exec"
	"strings"
	"strconv"
  "github.com/kless/shutil/sh"
)

// Combine the files located at the specified paths into a single pdf.
func Merge(doclist []string) {
	cpdf, err := exec.LookPath("cpdf")
	if err != nil {
		log.Fatal("no cpdf")
	}
	combine_cmd := exec.Command(cpdf)
	combine_cmd.Dir = "/tmp"
	combine_cmd.Args = cpdfMergeArgs(doclist)
	out, err := combine_cmd.Output()
	log.Println("stdout of command:", string(out))
}

func GetPageCount(filename string) (count int, err error) {
	out, err := runCpdfCommand("pages",filename)
	countStr := strings.Trim(string(out)," \n")
	count, err = strconv.Atoi(countStr)
	return
}

func runCpdfCommand(command,filename string) (output string, err error) {
	out, err := sh.Run(cpdf()+" -"+command+" " +filename)
	output = string(out)
	return
}

func cpdf() string {
	cpdf, err := exec.LookPath("cpdf")
	if err != nil {
		log.Fatal("no cpdf")
	}
	return cpdf
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
