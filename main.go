package main

import (
  "flag"
  "fmt"
  "os"
  "os/exec"
  "runtime"
)

/*
  use phantomJS to capture screen shots of a website
*/
func main() {
  threads := flag.Int("t", runtime.NumCPU(), "the number of concurrent pages to download")

  config := flag.String("c", "./job.json", "the job configuration file")

  var outDir string
  flag.StringVar(&outDir, "o", "", "directory in which to save the screen grabs")

  flag.Parse()


  jobSettings := LoadJobFile(*config)

  // make sure the output directory exists
  os.MkdirAll(outDir, os.ModePerm)

  chGen := grabJobChannelGenerator(outDir, jobSettings)

  chFanned := fanOut(*threads, grabImageToFile, chGen)

  fmt.Println("Capturing page images:")
  for n := range fanGrabJobsIn(chFanned...) {
    fmt.Println("\t", n.outFile)
  }

}

func grabImageToFile(grabJob GrabJob) GrabJob {

  cmd := fmt.Sprintf("%s %s %s \"%s\" %s %s", grabJob.PhantomJS, grabJob.ScreenGrab,
    grabJob.User, grabJob.Password, grabJob.page, grabJob.outFile)

  out, err := exec.Command("sh", "-c", cmd).Output()
  if err != nil {
    fmt.Println(cmd, out)
    fmt.Println(err)
  }

  return grabJob
}
