package process_utils

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Thread struct {
	Pid int
}

type Process struct {
	Pid      int
	Cmds     []string
	Children []Process
	Threads  []Thread
}

func GetProcess(pid int) (process *Process, err error) {
	pidStr := strconv.Itoa(pid)
	path := "/proc/" + pidStr
	_, tmpErr := os.Stat(path)
	if tmpErr == nil {
		// get cmdline
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

		// get threads
		var threads []Thread
		var pidDirs []os.FileInfo
		if pidDirs, err = ioutil.ReadDir(path + "/task"); err != nil {
			return
		}
		for _, pidDir := range pidDirs {
			if pidStr != pidDir.Name() {
				var threadPid int
				if threadPid, err = strconv.Atoi(pidDir.Name()); err != nil {
					return
				}
				threads = append(threads, Thread{
					Pid: threadPid,
				})
			}
		}

		// get children
		var children []Process
		childrenPath := path + "/task/" + pidStr + "/children"
		if bytes, err = ioutil.ReadFile(childrenPath); err != nil {
			return
		}
		splitedChildren := strings.Fields(string(bytes))
		for _, child := range splitedChildren {
			var childPid int
			if childPid, err = strconv.Atoi(child); err != nil {
				return
			}
			var childProcess *Process
			if childProcess, err = GetProcess(childPid); err != nil {
				return
			}
			children = append(children, *childProcess)
		}

		process = &Process{
			Pid:      pid,
			Cmds:     cmds,
			Threads:  threads,
			Children: children,
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
