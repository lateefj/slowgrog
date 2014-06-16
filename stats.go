package main

var (
	cmdCounts map[string]int64
	BrokenCmd []string
	BadCmd    []string
)

func init() {
	cmdCounts = make(map[string]int64)
	BadCmd = []string{"SMEMBERS"}
	BrokenCmd = []string{"KEYS"}
}

func IncCmdCount(cmd string) {
	v, exists := cmdCounts[cmd]
	if !exists {
		v = 0
	}
	v++
	cmdCounts[cmd] = v
}

func contains(s string, l []string) bool {
	for _, x := range l {
		if s == x {
			return true
		}
	}
	return false
}

func matchCmds(cmdList []string) map[string]int64 {
	cmds := make(map[string]int64)
	for k, v := range cmdCounts {
		if contains(k, cmdList) {
			cmds[k] = v
		}
	}
	return cmds
}

func GetBrokenCmds() map[string]int64 {
	return matchCmds(BrokenCmd)
}
func GetBadComds() map[string]int64 {
	return matchCmds(BadCmd)
}
