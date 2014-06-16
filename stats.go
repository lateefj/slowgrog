package main

var (
	// Commands that are documented to be blocking, performance issues or not to be used for production
	BadCmdList []string
)

func init() {

	BadCmdList = []string{
		"KEYS",     // KEYS is a blocking command that should NOT be used for production!!! (http://redis.io/commands/keys)
		"SMEMBERS", // SMEMBERS is a blocking command that should not be used if possible instead use SCAN!! (http://redis.io/topics/latency)
	}
}

type Stats struct {
	cmdCounts map[string]int64
}

func NewStats() *Stats {
	return &Stats{cmdCounts: make(map[string]int64)}
}

func (s *Stats) IncCmdCount(cmd string) {
	v, exists := s.cmdCounts[cmd]
	if !exists {
		v = 0
	}
	v++
	s.cmdCounts[cmd] = v
}

func contains(s string, l []string) bool {
	for _, x := range l {
		if s == x {
			return true
		}
	}
	return false
}

func matchCmds(cmdList []string, counts map[string]int64) map[string]int64 {
	cmds := make(map[string]int64)
	for k, v := range counts {
		if contains(k, cmdList) {
			cmds[k] = v
		}
	}
	return cmds
}
func (s *Stats) Counts() map[string]int64 {
	return s.cmdCounts
}

func (s *Stats) BadCmds() map[string]int64 {
	return matchCmds(BadCmdList, s.cmdCounts)
}
