package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	Trace               *log.Logger
	Info                *log.Logger
	Warning             *log.Logger
	Error               *log.Logger
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

	Status = &RedisStatus{Info: make(map[string]interface{}), Slowlog: make([]string, 0), MonitorSample: make([]*MonitorCmd, 0)}
	flag.IntVar(&CmdLimit, "cmdlimit", 100, "number of commands the monitor will store")
	flag.IntVar(&Frequency, "frequency", 10000, "Number of miliseconds to delay between samples info, slowlog")
	flag.IntVar(&MonitorSampleLength, "monsamplen", 1000, "Length of miliseconds that the monitor is sampled (0 will be coninuous however this is very costly to performance)")
	flag.IntVar(&SlowlogSize, "slogsize", 10, "slowlog size")
	flag.StringVar(&RedisHost, "h", "127.0.0.1", "redis host ")
	flag.StringVar(&RedisPassword, "a", "", "Redis password")
	flag.IntVar(&RedisPort, "p", 6379, "redis port")
}

type CommandStats struct {
	CommandCounts map[string]int64 `json:"command_counts"`
	BadCommands   map[string]int64 `json:"bad_commands"`
}
type RedisStatus struct {
	CommandStats  CommandStats           `json:"stats"`
	Slowlog       []string               `json:"slowlog"`
	Info          map[string]interface{} `json:"info"`
	MonitorSample []*MonitorCmd          `json:"monitor_sample"`
	stats         *Stats
}

func rcon() (redis.Conn, error) {
	Trace.Printf("Connecting to redis host %s on port %d", RedisHost, RedisPort)
	c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", RedisHost, RedisPort))
	if RedisPassword != "" {
		if _, err := c.Do("AUTH", RedisPassword); err != nil {
			c.Close()
			return c, err
		}
	}
	return c, err
}

func main() {
	flag.Parse()
	go func() {
		for {
			stat := &RedisStatus{Info: make(map[string]interface{}), Slowlog: make([]string, 0), MonitorSample: make([]*MonitorCmd, 0), CommandStats: CommandStats{CommandCounts: make(map[string]int64), BadCommands: make(map[string]int64)}, stats: NewStats()}
			c, err := rcon()
			if err != nil {
				Error.Printf("Failed to make connection %s", err)
				continue
			}
			go SampleMonitor(stat)
			_, err = SampleInfo(c, stat)
			if err != nil {
				Error.Println(err)
			}
			_, err = SampleSlowlog(c, stat)
			if err != nil {
				Error.Printf("Error with slowlog %s", err)
			}
			Status = stat
			c.Close()
			time.Sleep(time.Duration(Frequency) * time.Second)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		// TODO: Fix this just added before hack session timed out
		Status.CommandStats.CommandCounts = Status.stats.Counts()
		Status.CommandStats.BadCommands = Status.stats.BadCmds()
		enc.Encode(Status)
	})
	log.Fatal(http.ListenAndServe(":8000", nil))
}
