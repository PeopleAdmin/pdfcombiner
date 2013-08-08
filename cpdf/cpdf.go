// Package cpdf contains methods to manipulate pdf files using cpdf.
package cpdf

import(
  "os/exec"
  "log"
)

// Given a slice of documents, construct an args slice suitable for
// passing to 'cpdf merge'.
func cpdfMergeArgs(doclist []string) (args []string) {
  args = []string{"merge"}
  for _,doc := range doclist{
    args = append(args, doc)
  }
  args = append(args, "-o", "combined.pdf")
  return
}

// Combine the files located at the specified paths into a single pdf.
func Merge(doclist []string) {
  cpdf, err := exec.LookPath("cpdf")
  if err != nil { log.Fatal("no cpdf") }
  combine_cmd := exec.Command(cpdf)
  combine_cmd.Dir = "/tmp"
  combine_cmd.Args = cpdfMergeArgs(doclist)
  out, err := combine_cmd.Output()
  log.Println("stdout of command:",string(out))
}

