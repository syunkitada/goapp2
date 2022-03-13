package cmd_runner

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProcess(t *testing.T) {
	a := assert.New(t)

	pidCh := make(chan int)
	command := exec.Command("sleep", "1")
	go func(pidCh chan int) {
		command.Start()
		pidCh <- command.Process.Pid
		command.Wait()
	}(pidCh)

	{
		pid := <-pidCh
		process, err := GetProcess(pid)
		a.Equal(err, nil)
		expected := Process{
			Pid:  process.Pid,
			Cmds: []string{"sleep", "1"},
		}
		a.Equal(expected, *process)
	}

	{
		process, err := GetProcess(-1)
		a.Equal(err, nil)
		var expectedProcess *Process
		a.Equal(expectedProcess, process)
	}
}
