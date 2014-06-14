package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Status  *RedisStatus
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
}

type RedisStatus struct {
	Info          map[string]interface{} `json:"info"`
	Slowlog       []string               `json:"slowlog"`
	MonitorSample []*MonitorCmd          `json:"monitor_sample"`
}

func main() {
	go func() {
		for {
			stat := &RedisStatus{Info: make(map[string]interface{}), Slowlog: make([]string, 0), MonitorSample: make([]*MonitorCmd, 0)}
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				// handle error
			}
			defer c.Close()

			_, err = SampleInfo(c, stat)
			if err != nil {
				Error.Println(err)
			}
			fmt.Printf("Info: %s\n", stat.Info)
			SampleMonitor(c, stat)
			_, err = SampleSlowlog(c, stat)
			if err != nil {
				Error.Println(err)
			}
			Status = stat
			time.Sleep(5 * time.Second)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		enc.Encode(Status)
	})
	log.Fatal(http.ListenAndServe(":8000", nil))
}
