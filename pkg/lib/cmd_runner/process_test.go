package cmd_runner

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProcess(t *testing.T) {
	a := assert.New(t)

	command1 := exec.Command("sleep", "10")
	command1.Start()
	pid1 := command1.Process.Pid

	command2 := exec.Command("sleep", "12")
	command2.Start()
	pid2 := command2.Process.Pid

	{
		// 単体プロセスのテスト
		process, err := GetProcess(pid1)
		a.NoError(err)
		a.NotNil(process)
		expected := Process{
			Pid:      pid1,
			Cmds:     []string{"sleep", "10"},
			Children: []Process{},
		}
		a.Equal(expected, *process)

		// go testの子プロセスとしてsleepプロセスがあることを確認する
		selfPid := os.Getpid()
		selfProcess, err := GetProcess(selfPid)
		a.NoError(err)
		expectedChildren := []Process{
			Process{Pid: pid1, Cmds: []string{"sleep", "10"}, Children: []Process{}},
			Process{Pid: pid2, Cmds: []string{"sleep", "12"}, Children: []Process{}},
		}
		a.ElementsMatch(expectedChildren, selfProcess.Children)
	}

	{
		// 存在しないPIDのテスト
		process, err := GetProcess(-1)
		a.Equal(err, nil)
		var expectedProcess *Process
		a.Equal(expectedProcess, process)
	}
}
