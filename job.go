package main

import (
  "io/ioutil"
  "fmt"
  "encoding/json"
)

type Job struct {
  Name       string
  User       string
  Password   string
  Host       string
  PhantomJS  string
  ScreenGrab string
  Pages      []string
}


// load a job file
func LoadJobFile(jobFile string) Job {
  var job Job

  file, e := ioutil.ReadFile(jobFile)
  if e != nil {
    fmt.Println("Failed to load job file", jobFile, e)
  }
  json.Unmarshal(file, &job)

  return job
}
