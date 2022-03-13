package cmd_runner

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Process struct {
	Pid  int
	Cmds []string
}

func GetProcess(pid int) (process *Process, err error) {
	path := "/proc/" + strconv.Itoa(pid)
	_, tmpErr := os.Stat(path)
	if tmpErr == nil {
		cmdPath := path + "/cmdline"
		bytes, tmpErr := ioutil.ReadFile(cmdPath)
		if tmpErr != nil {
			err = tmpErr
			return
		}
		cmds := strings.Split(string(bytes), string(byte(0)))
		if cmds[len(cmds)-1] == "" {
			cmds = cmds[0 : len(cmds)-1]
		}

		process = &Process{
			Pid:  pid,
			Cmds: cmds,
		}
		return
	}
	if os.IsNotExist(tmpErr) {
		return
	} else {
		err = tmpErr
	}
	return
}
