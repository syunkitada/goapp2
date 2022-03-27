package process_utils

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProcess(t *testing.T) {
	a := assert.New(t)

	command1 := exec.Command("sleep", "1")
	command1.Start()
	pid1 := command1.Process.Pid

	// command2 := exec.Command("sleep", "1")
	// command2.Start()
	// pid2 := command2.Process.Pid

	{
		// 単体プロセスのテスト
		process, err := GetProcessFromPid(pid1)
		a.NoError(err)
		expected := Process{
			Pid:  pid1,
			Cmds: []string{"sleep", "1"},
		}
		a.Equal(expected, *process)

		// go testの子プロセスとしてsleepプロセスがあることを確認する
		// selfPid := os.Getpid()
		// selfProcess, err := GetProcessFromPid(selfPid)
		// a.NoError(err)
		// a.Equal(selfPid, selfProcess.Pid)
		// expectedChildren := []Process{
		// 	Process{Pid: pid1, Cmds: []string{"sleep", "1"}},
		// 	Process{Pid: pid2, Cmds: []string{"sleep", "1"}},
		// }
		// a.Equal(expectedChildren, selfProcess.Children)

		// // go testはマルチスレッドなので一定数以上のスレッドが見つかることを確認する
		// a.Greater(len(selfProcess.Threads), 2)
		// // threadのpidは親のPIDよりも後であることを確認する
		// a.Greater(selfProcess.Threads[0].Pid, selfPid)
	}

	{
		// 存在しないPIDのテスト
		process, err := GetProcessFromPid(-1)
		a.Equal(err, nil)
		var expectedProcess *Process
		a.Equal(expectedProcess, process)
	}
}

func TestGetProcesses(t *testing.T) {
	a := assert.New(t)
	a.Equal(true, true)

	processes, err := GetProcesses()
	a.NoError(err)
	fmt.Println("process", processes)
}
