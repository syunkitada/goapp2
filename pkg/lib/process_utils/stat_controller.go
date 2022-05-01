package process_utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/syunkitada/goapp2/pkg/lib/runner"
)

type StatControllerConfig struct {
	runner.Config
	TargetProcess string
	TargetPid     int
}

type StatController struct {
	runner.Runner
	StatRunner *StatRunner
}

func NewStatController(conf *StatControllerConfig) (statController *StatController) {
	ecmd := exec.Command("getconf", "CLK_TCK")
	out := new(bytes.Buffer)
	ecmd.Stdout = out
	tmpErr := ecmd.Run()
	if tmpErr != nil {
		os.Exit(1)
	}
	clkTck, tmpErr := strconv.Atoi(strings.TrimSpace(out.String()))
	if tmpErr != nil {
		os.Exit(1)
	}
	statRunner := StatRunner{
		ClkTck:        clkTck,
		TargetProcess: conf.TargetProcess,
		TargetPid:     conf.TargetPid,
	}
	statController = &StatController{
		Runner:     *runner.New(&conf.Config, &statRunner),
		statRunner: &statRunner,
	}
	return
}

type StatRunner struct {
	ClkTck            int
	TargetProcess     string
	TargetPid         int
	beforeProcesses   []Process
	beforePidIndexMap map[int]int
}

func (self *StatRunner) Run(runAt time.Time) {
	processes, pidIndexMap, err := GetProcesses("/", true)
	if err != nil {
		return
	}
	if self.beforeProcesses == nil {
		self.beforeProcesses = processes
		self.beforePidIndexMap = pidIndexMap
		return
	}

	for i, process := range processes {
		beforeProcessIndex, ok := self.beforePidIndexMap[process.Pid]
		if !ok {
			continue
		}
		beforeProcess := self.beforeProcesses[beforeProcessIndex]
		stat := process.Stat
		bstat := beforeProcess.Stat
		interval := int(stat.Timestamp.Unix() - bstat.Timestamp.Unix())

		stat.UserUtil = (stat.Utime - bstat.Utime) / interval
		stat.SystemUtil = (stat.Stime - bstat.Stime) / interval
		stat.GuestUtil = (stat.Gtime - bstat.Gtime) / interval
		stat.CguestUtil = (stat.Cgtime - bstat.Cgtime) / interval

		stat.SchedTimeSlicesPerSec = (stat.SchedTimeSlices - bstat.SchedTimeSlices) / interval
		stat.SchedCpuTimePerSec = (stat.SchedCpuTime - bstat.SchedCpuTime) / interval
		stat.WaitUtil = ((stat.SchedWaitTime - bstat.SchedWaitTime) * self.ClkTck) / 1000000000
		stat.VoluntaryCtxtSwitchesPerSec = (stat.VoluntaryCtxtSwitches - bstat.VoluntaryCtxtSwitches) / interval
		stat.NonvoluntaryCtxtSwitches = (stat.NonvoluntaryCtxtSwitches - bstat.NonvoluntaryCtxtSwitches) / interval

		stat.SyscrPerSec = (stat.Syscr - bstat.Syscr) / interval
		stat.SyscwPerSec = (stat.Syscw - bstat.Syscw) / interval
		stat.ReadBytesPerSec = (stat.ReadBytes - bstat.ReadBytes) / interval
		stat.WriteBytesPerSec = (stat.WriteBytes - bstat.WriteBytes) / interval
		processes[i].Stat = stat
	}

	targetPid := self.TargetPid
	targetProcess := self.TargetProcess
	for _, process := range processes {
		if targetPid != 0 && process.Pid != targetPid {
			continue
		}
		if targetProcess != "" && !strings.Contains(process.Name, targetProcess) {
			continue
		}
		fmt.Println(strconv.Itoa(process.Pid), process.Name, strconv.Itoa(process.Stat.UserUtil), strconv.Itoa(process.Stat.WaitUtil))
	}

	fmt.Println("run", runAt)
	self.beforeProcesses = processes
	self.beforePidIndexMap = pidIndexMap
}

func (self *StatRunner) StopTimeout() {
}

func (self *StatController) ShowPs() {
	processes, _, err := GetProcesses("/", true)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println("DEBUG", processes)

	// table := tablewriter.NewWriter(os.Stdout)
	// table.SetHeader([]string{"Pid", "Name", "Cmd"})
	// for _, process := range processes {
	// 	table.Append([]string{strconv.Itoa(process.Pid), process.Name, strings.Join(process.Cmds, " ")})
	// }
	// table.Render()
}
