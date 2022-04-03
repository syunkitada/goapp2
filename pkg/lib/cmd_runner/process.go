package cmd_runner

import (
	"fmt"

	"github.com/syunkitada/goapp2/pkg/lib/process_utils"
)

type Process struct {
	Pid      int
	Cmds     []string
	Children []Process
}

func GetProcess(pid int) (process *Process, err error) {
	var processes []process_utils.Process
	var pidIndexMap map[int]int
	if processes, pidIndexMap, err = process_utils.GetProcesses("/"); err != nil {
		return
	}

	process, err = getProcess(processes, pidIndexMap, pid)
	return
}

func getProcess(processes []process_utils.Process, pidIndexMap map[int]int, pid int) (process *Process, err error) {
	pidIndex, ok := pidIndexMap[pid]
	if !ok {
		return
	}

	tmpProcess := processes[pidIndex]

	children := []Process{}
	for _, cpid := range tmpProcess.Children {
		var child *Process
		if child, err = getProcess(processes, pidIndexMap, cpid); err != nil {
			return
		} else if child == nil {
			err = fmt.Errorf("Unexpected child is found: pid=%d, child=%d", pid, cpid)
		}
		children = append(children, *child)
	}

	process = &Process{
		Pid:      pid,
		Cmds:     tmpProcess.Cmds,
		Children: children,
	}
	return
}
