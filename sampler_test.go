package main

import (
	"testing"
)

// Clean this up don't need all this info
var infoStuff string

func init() {
	infoStuff = `
# Server
redis_version:2.8.11
redis_git_sha1:00000000
redis_git_dirty:0
redis_build_id:cd43e547b41f72f
redis_mode:standalone
os:Darwin 13.2.0 x86_64
arch_bits:64
multiplexing_api:kqueue
gcc_version:4.2.1
process_id:60108
run_id:8168a2a6810d6f985642717179318fe6d86e4a77
tcp_port:6379
uptime_in_seconds:83676
uptime_in_days:0
hz:10
lru_clock:11286124
config_file:/usr/local/etc/redis.conf

# Clients
connected_clients:1
client_longest_output_list:0
client_biggest_input_buf:0
blocked_clients:0

# Memory
used_memory:1581488
used_memory_human:1.51M
used_memory_rss:16547840
used_memory_peak:9486960
used_memory_peak_human:9.05M
used_memory_lua:33792
mem_fragmentation_ratio:10.46
mem_allocator:libc

# Persistence
loading:0
rdb_changes_since_last_save:2190
rdb_bgsave_in_progress:0
rdb_last_save_time:1403794893
rdb_last_bgsave_status:ok
rdb_last_bgsave_time_sec:0
rdb_current_bgsave_time_sec:-1
aof_enabled:0
aof_rewrite_in_progress:0
aof_rewrite_scheduled:0
aof_last_rewrite_time_sec:-1
aof_current_rewrite_time_sec:-1
aof_last_bgrewrite_status:ok
aof_last_write_status:ok

# Stats
total_connections_received:69214
total_commands_processed:14383
instantaneous_ops_per_sec:0
rejected_connections:0
sync_full:0
sync_partial_ok:0
sync_partial_err:0
expired_keys:0
evicted_keys:0
keyspace_hits:0
keyspace_misses:0
pubsub_channels:0
pubsub_patterns:0
latest_fork_usec:1914

# Replication
role:master
connected_slaves:0
master_repl_offset:0
repl_backlog_active:0
repl_backlog_size:1048576
repl_backlog_first_byte_offset:0
repl_backlog_histlen:0

# CPU
used_cpu_sys:7.41
used_cpu_user:3.41
used_cpu_sys_children:0.00
used_cpu_user_children:0.01

# Keyspace
db0:keys=8016,expires=0,avg_ttl=0
`
}

type mockDataCmds struct {
	Info string
}

func (mdc *mockDataCmds) InfoCmd() (string, error) {
	return mdc.Info, nil
}

func (mdc *mockDataCmds) SlowlogCmd() ([]Slowlog, error) {
	logs := make([]Slowlog, 0)
	return logs, nil
}
func (mdc *mockDataCmds) MonitorCmd(chan bool) chan string {
	return make(chan string, MONITOR_BUFFER_SIZE)
}

func TestSampleInfo(t *testing.T) {
	mdc := &mockDataCmds{Info: infoStuff}
	status := &RedisStatus{Info: make(map[string]interface{}), Slowlogs: make([]Slowlog, 0), MonitorSample: make([]*MonitorCmd, 0), stats: NewStats()}
	SampleInfo(mdc, status)
	if status.Info["role"].(string) != "master" {
		t.Fatalf("Failed to find 'role' equal 'master in Info'")
	}
}
func TestMonitorSample(t *testing.T) {
	// Place holder...
}
