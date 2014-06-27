package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

var Logger = logrus.New()

var (
	Status              *RedisStatus
	CmdLimit            int
	Frequency           int
	MonitorSampleLength int
	SlowlogSize         int
	RedisHost           string
	RedisPort           int
	RedisPassword       string
	WebPort             int
)

func init() {

	Logger.Formatter = new(logrus.TextFormatter)
	//Logger.Level = logrus.Info
	Logger.Level = logrus.Debug

	Status = &RedisStatus{Info: make(map[string]interface{}), Slowlogs: make([]Slowlog, 0), MonitorSample: make([]*MonitorCmd, 0), stats: NewStats()}
	flag.StringVar(&RedisHost, "h", "127.0.0.1", "Redis host ")
	flag.IntVar(&RedisPort, "p", 6379, "Redis port")
	flag.StringVar(&RedisPassword, "a", "", "Redis password")

	flag.IntVar(&WebPort, "webport", 7071, "Port to run the http server on")
	flag.IntVar(&CmdLimit, "cmdlimit", 100, "number of commands the MONITOR will store")
	flag.IntVar(&Frequency, "frequency", 60, "Number of seconds to delay between samples INFO, SLOWLOG")
	flag.IntVar(&MonitorSampleLength, "monsamplen", 1000, "Length of miliseconds that the monitor is sampled (0 will be coninuous however this is very costly to performance)")
	flag.IntVar(&SlowlogSize, "slogsize", 10, "SLOWLOG size")
}

type CommandStats struct {
	CommandCounts map[string]int64 `json:"command_counts"`
	BadCommands   map[string]int64 `json:"bad_commands"`
}

type RedisStatus struct {
	CommandStats  CommandStats           `json:"stats"`
	Slowlogs      []Slowlog              `json:"slowlogs"`
	Info          map[string]interface{} `json:"info"`
	MonitorSample []*MonitorCmd          `json:"monitor_sample"`
	stats         *Stats
}

func main() {
	flag.Parse()
	rc := NewRedisCmds()
	stopper := make(chan bool, 1)
	go SampleMonitor(rc, stopper, Status)
	go func() {
		for {
			_, err := SampleInfo(rc, Status)
			if err != nil {
				Logger.Errorf("SampleInfo error %s\n", err)
			}
			_, err = SampleSlowlog(rc, Status)
			if err != nil {
				Logger.Errorf("Error with slowLogger %s", err)
			}
			time.Sleep(time.Duration(Frequency) * time.Second)
		}
	}()

	Logger.Debug("Starting up the http handler")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		// TODO: Fix this just added before hack session timed out
		Status.CommandStats.CommandCounts = Status.stats.Counts()
		Status.CommandStats.BadCommands = Status.stats.BadCmds()
		enc.Encode(Status)
	})
	Logger.Fatalf("Failed %s", http.ListenAndServe(fmt.Sprintf(":%d", WebPort), nil))
}
