package main

import (
	"encoding/json"
	"flag"
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
)

func init() {

	Logger.Formatter = new(logrus.TextFormatter)

	Status = &RedisStatus{Info: make(map[string]interface{}), Slowlogs: make([]Slowlog, 0), MonitorSample: make([]*MonitorCmd, 0), stats: NewStats()}
	flag.StringVar(&RedisHost, "h", "127.0.0.1", "redis host ")
	flag.IntVar(&RedisPort, "p", 6379, "redis port")
	flag.StringVar(&RedisPassword, "a", "", "Redis password")

	flag.IntVar(&CmdLimit, "cmdlimit", 100, "number of commands the monitor will store")
	flag.IntVar(&Frequency, "frequency", 10000, "Number of miliseconds to delay between samples info, slowLogger")
	flag.IntVar(&MonitorSampleLength, "monsamplen", 1000, "Length of miliseconds that the monitor is sampled (0 will be coninuous however this is very costly to performance)")
	flag.IntVar(&SlowlogSize, "sLoggersize", 10, "slowLogger size")
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
	go func() {
		go SampleMonitor(Status)
		for {
			c, err := rcon()
			if err != nil {

				Logger.Errorf("Failed to make connection %s", err)
				continue
			}
			_, err = SampleInfo(c, Status)
			if err != nil {
				Logger.Errorf("SampleInfo error %s\n", err)
			}
			_, err = SampleSlowlog(c, Status)
			if err != nil {
				Logger.Errorf("Error with slowLogger %s", err)
			}
			c.Close()
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
	Logger.Fatalf("Failed %s", http.ListenAndServe(":8000", nil))
}
