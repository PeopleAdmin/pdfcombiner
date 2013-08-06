package cpdf

import(
  "os/exec"
  "log"
)

func cpdfArgs(doclist []string) (args []string) {
  args = []string{"merge"}
  for _,doc := range doclist{
    args = append(args, doc)
  }
  args = append(args, "-o", "combined.pdf")
  return
}

func Merge(doclist []string) {
  cpdf, err := exec.LookPath("cpdf")
  if err != nil { log.Fatal("no cpdf") }
  combine_cmd := exec.Command(cpdf)
  combine_cmd.Dir = "/tmp"
  combine_cmd.Args = cpdfArgs(doclist)
  out, err := combine_cmd.Output()
  log.Println("stdout of command:",string(out))
}

