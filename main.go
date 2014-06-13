package main

import (
  "fmt"
  "log"
  "os"

  "github.com/garyburd/redigo/redis"
)

var (
  Trace   *log.Logger
  Info    *log.Logger
  Warning *log.Logger
  Error   *log.Logger
)

func init() {

  Trace = log.New(os.Stdout,
    "TRACE: ",
    log.Ldate|log.Ltime|log.Lshortfile)

  Info = log.New(os.Stdout,
    "INFO: ",
    log.Ldate|log.Ltime|log.Lshortfile)

  Warning = log.New(os.Stderr,
    "WARNING: ",
    log.Ldate|log.Ltime|log.Lshortfile)

  Error = log.New(os.Stderr,
    "ERROR: ",
    log.Ldate|log.Ltime|log.Lshortfile)
}

type Status struct {
  Info          map[string]interface{} `json:"info"`
  Slowlog       []string               `json:"slowlog"`
  MonitorSample []string               `json:"monitor_sample"`
}

func main() {
  status := &Status{Info: make(map[string]interface{}), Slowlog: make([]string, 0), MonitorSample: make([]string, 0)}
  c, err := redis.Dial("tcp", ":6379")
  if err != nil {
    // handle error
  }
  defer c.Close()

  _, err = SampleInfo(c, status)
  if err != nil {
    Error.Println(err)
  }
  fmt.Printf("Info: %s\n", status.Info)
  SampleMonitor(c, status)
  ss, err := SampleSlowlog(c, status)
  if err != nil {
    Error.Println(err)
  }
  fmt.Printf("Slowlog: %s\n", ss)
}
