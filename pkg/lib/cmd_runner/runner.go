package cmd_runner

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/syunkitada/goapp2/pkg/lib/process_utils"
)

const (
	StatusTimeout = 124
	StatusErr     = 500
)

type Config struct {
	Timeout  int
	Interval int
	Cmd      string
	UseShell bool
}

type Runner struct {
	conf         *Config
	cmd          string
	cmdOptions   []string
	cmdTimeout   time.Duration
	killLimit    int
	killInterval time.Duration
	stopCh       chan chan bool
	isStarted    bool
}

func New(conf *Config) *Runner {
	cmds := strings.Fields(conf.Cmd)
	cmd := cmds[0]
	cmdOptions := cmds[1:]
	timeout := conf.Timeout
	if timeout == 0 {
		timeout = conf.Interval - 10
	}
	cmdTimeout := time.Duration(timeout) * time.Second
	return &Runner{
		conf:         conf,
		cmd:          cmd,
		cmdOptions:   cmdOptions,
		cmdTimeout:   cmdTimeout,
		killLimit:    2,
		killInterval: time.Duration(1) * time.Second,
		stopCh:       make(chan chan bool),
	}
}

func (self *Runner) Start() {
	fmt.Println("DEBUG Start")
	if !self.isStarted {
		go self.start()
		self.isStarted = true
	} else {
		log.Printf("Already Started")
	}
}

func (self *Runner) Stop() {
	log.Printf("stopping: %v", self.conf.Cmd)
	doneCh := make(chan bool)
	self.stopCh <- doneCh
	<-doneCh
	fmt.Println("end stop")
}

func (self *Runner) start() {
	interval := time.Duration(self.conf.Interval) * time.Second
	ticker := time.NewTicker(interval)
	log.Printf("start: %s", self.conf.Cmd)

	self.Run()
	for {
		select {
		case doneCh := <-self.stopCh:
			doneCh <- true
			log.Printf("done: %s", self.conf.Cmd)
			return
		case t := <-ticker.C:
			fmt.Println("tick at", t)
			self.Run()
		}
	}
	return
}

type Result struct {
	Cmd    string
	Err    error
	Output string
	Status int
}

func (self *Runner) MustKillProcess(pid int) {
	if pid == 0 {
		return
	}
	for i := 0; i < self.killLimit; i++ {
		process, err := process_utils.GetProcessFromPid(pid)
		if err != nil {
			log.Fatalf("Unexpected Error: %s", err.Error())
		}
		if process == nil {
			return
		}

		if self.conf.UseShell && process.Cmds[2] != self.conf.Cmd {
			log.Fatalf("Unexpected Cmd Found: expectedCmd=%s, foundCmd=%v", self.cmd, process.Cmds)
		} else if process.Cmds[0] != self.cmd {
			log.Fatalf("Unexpected Cmd Found: expectedCmd=%s, foundCmd=%v", self.cmd, process.Cmds)
		}

		log.Printf("ExistsProcess will be killed: pid=%d, cmds=%v", pid, process.Cmds)
		if err = syscall.Kill(-pid, syscall.SIGKILL); err != nil {
			log.Fatalf("Unexpected Error: %s", err.Error())
			return
		}
		time.Sleep(self.killInterval)
	}

	process, err := process_utils.GetProcessFromPid(pid)
	if err != nil {
		log.Fatalf("Unexpected Error: %s", err.Error())
	}
	if process != nil {
		if process.Cmds[0] != self.cmd {
			log.Fatalf("Failed KillProcess, and Unexpected Cmd Found: expectedCmd=%s, foundCmd=%v", self.cmd, process.Cmds)
		}
		log.Fatalf("Failed KillProcess: pid=%d, cmds=%v", pid, process.Cmds)
	}
}

func (self *Runner) Run() (result *Result, err error) {
	log.Printf("Run: %s", self.conf.Cmd)
	var pid int
	var tmpCmdResult *cmdResult
	var errStatus int

	defer func() {
		self.MustKillProcess(pid)

		result = &Result{
			Cmd: self.conf.Cmd,
		}
		if errStatus == StatusTimeout {
			result.Output = fmt.Sprintf("Failed command by timeout")
			result.Status = StatusTimeout
		} else {
			if tmpCmdResult != nil {
				if tmpCmdResult.err != nil {
					result.Output = fmt.Sprintf("Failed command: err=%s", tmpCmdResult.err.Error())
					result.Status = StatusErr
				} else {
					result.Status = tmpCmdResult.exitCode
					tmpOutput := []string{}
					if len(tmpCmdResult.stdout) > 0 {
						tmpOutput = append(tmpOutput, string(tmpCmdResult.stdout))
					}
					if len(tmpCmdResult.stderr) > 0 {
						tmpOutput = append(tmpOutput, string(tmpCmdResult.stderr))
					}
					result.Output = strings.Join(tmpOutput, "\n")
				}
			} else {
				if err != nil {
					result.Output = fmt.Sprintf("Failed command: err=%s", err.Error())
					result.Status = StatusErr
				}
			}
		}
		log.Printf("EndRun: %s", self.conf.Cmd)
		return
	}()

	var command *exec.Cmd
	if self.conf.UseShell {
		command = exec.Command("sh", "-c", self.conf.Cmd)
	} else {
		command = exec.Command(self.cmd, self.cmdOptions...)
	}
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var stdoutPipe io.ReadCloser
	if stdoutPipe, err = command.StdoutPipe(); err != nil {
		return
	}

	var stderrPipe io.ReadCloser
	if stderrPipe, err = command.StderrPipe(); err != nil {
		return
	}

	if err = command.Start(); err != nil {
		return
	}
	pid = command.Process.Pid

	cmdResultCh := make(chan *cmdResult)
	go cmdWait(command, stdoutPipe, stderrPipe, cmdResultCh)

	select {
	case tmpCmdResult = <-cmdResultCh:
		return
	case <-time.After(self.cmdTimeout):
		errStatus = StatusTimeout
	}
	return
}

type cmdResult struct {
	stdout   []byte
	stderr   []byte
	exitCode int
	err      error
}

func cmdWait(cmd *exec.Cmd, stdoutPipe io.ReadCloser, stderrPipe io.ReadCloser, resultCh chan *cmdResult) {
	var stdout []byte
	var stderr []byte
	var exitCode int
	var err error
	defer func() {
		resultCh <- &cmdResult{
			stdout:   stdout,
			stderr:   stderr,
			exitCode: exitCode,
			err:      err,
		}
	}()

	if stdout, err = ioutil.ReadAll(stdoutPipe); err != nil {
		return
	}
	if stderr, err = ioutil.ReadAll(stderrPipe); err != nil {
		return
	}
	if err = cmd.Wait(); err != nil {
		return
	}
	exitCode = cmd.ProcessState.ExitCode()
}
