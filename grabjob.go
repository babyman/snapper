package main

import (
  "sync"
  "path/filepath"
  "strings"
)

type GrabJob struct {
  page       string
  outFile    string
  PhantomJS  string
  ScreenGrab string
  User       string
  Password   string
}

type GrabJobTask func(GrabJob) GrabJob

func grabJobChannelGenerator(outDir string, settings Job) <-chan GrabJob {
  chOut := make(chan GrabJob)
  go func(jobSettings Job) {
    for _, v := range settings.Pages {
      grab := GrabJob{
        jobSettings.Host + v,
        filepath.Join(outDir, filename(v)),
        jobSettings.PhantomJS,
        jobSettings.ScreenGrab,
        jobSettings.User,
        jobSettings.Password,
      }
      chOut <- grab
    }
    close(chOut)
  }(settings)
  return chOut
}

func fanOut(count int, task GrabJobTask, chIn <-chan GrabJob) []<-chan GrabJob {

  var chFanned []<-chan GrabJob

  for i := 0; i < count; i++ {
    chFanned = append(chFanned, performGrabJobTask(task, chIn))
  }

  return chFanned
}

func performGrabJobTask(task GrabJobTask, chIn <-chan GrabJob) <-chan GrabJob {
  chOut := make(chan GrabJob)
  go func() {
    for grabJob := range chIn {
      chOut <- task(grabJob)
    }
    close(chOut)
  }()
  return chOut
}

func fanGrabJobsIn(chIns ...<-chan GrabJob) <-chan GrabJob {
  chOut := make(chan GrabJob)

  var wg sync.WaitGroup
  wg.Add(len(chIns))

  go func() {
    for _, v := range chIns {
      go func(chIn <-chan GrabJob) {
        for i := range chIn {
          chOut <- i
        }
        wg.Done()
      }(v)
    }
  }()

  go func() {
    wg.Wait()
    close(chOut)
  }()

  return chOut
}

func filename(path string) string {
  out := strings.Replace(path, "/", "_", -1)
  return out + ".png"
}
