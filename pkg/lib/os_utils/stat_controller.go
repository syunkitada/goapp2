package os_utils

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/syunkitada/goapp2/pkg/lib/runner"
)

type StatControllerConfig struct {
	runner.Config
	HandleProcesses func(processes []Process)
	HandleCpuStat   func(cpuStat *CpuStat)
	HandleMemStat   func(memStat *MemStat)
}

type StatController struct {
	runner.Runner
	statRunner *StatRunner
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
		clkTck:          clkTck,
		handleCpuStat:   conf.HandleCpuStat,
		handleMemStat:   conf.HandleMemStat,
		handleProcesses: conf.HandleProcesses,
		interval:        conf.Config.Interval,
	}
	statController = &StatController{
		Runner:     *runner.New(&conf.Config, &statRunner),
		statRunner: &statRunner,
	}
	return
}

type StatRunner struct {
	clkTck            int
	interval          int
	beforeCpuStat     *CpuStat
	handleCpuStat     func(cpuStat *CpuStat)
	beforeMemStat     *MemStat
	handleMemStat     func(memStat *MemStat)
	handleProcesses   func(processes []Process)
	beforeProcesses   []Process
	beforePidIndexMap map[int]int
}

func (self *StatRunner) syncCpuStat() {
	cpuStat, err := GetCpuStat()
	if err != nil {
		return
	}

	if self.beforeCpuStat == nil {
		self.beforeCpuStat = cpuStat
		return
	}

	interval := self.interval

	cpuStat.IntrPerSec = (cpuStat.Intr - self.beforeCpuStat.Intr) / interval
	cpuStat.CtxPerSec = (cpuStat.Ctx - self.beforeCpuStat.Ctx) / interval
	cpuStat.BtimePerSec = (cpuStat.Btime - self.beforeCpuStat.Btime) / interval
	cpuStat.ProcessesPerSec = (cpuStat.Processes - self.beforeCpuStat.Processes) / interval
	cpuStat.SoftirqPerSec = (cpuStat.Softirq - self.beforeCpuStat.Softirq) / interval

	if self.handleCpuStat != nil {
		self.handleCpuStat(cpuStat)
	}

	self.beforeCpuStat = cpuStat
}

func (self *StatRunner) syncMemStat() {
	var memStat *MemStat
	var err error
	if memStat, err = GetMemStat("/"); err != nil {
		return
	}

	if self.beforeMemStat == nil {
		self.beforeMemStat = memStat
		return
	}

	interval := self.interval

	memStat.Vmstat.PgscanKswapdPerSec = (memStat.Vmstat.PgscanKswapd - self.beforeMemStat.Vmstat.PgscanKswapd) / interval
	memStat.Vmstat.PgscanDirectPerSec = (memStat.Vmstat.PgscanDirect - self.beforeMemStat.Vmstat.PgscanDirect) / interval
	memStat.Vmstat.PgfaultPerSec = (memStat.Vmstat.Pgfault - self.beforeMemStat.Vmstat.Pgfault) / interval
	memStat.Vmstat.PswapinPerSec = (memStat.Vmstat.Pswapin - self.beforeMemStat.Vmstat.Pswapin) / interval
	memStat.Vmstat.PswapoutPerSec = (memStat.Vmstat.Pswapout - self.beforeMemStat.Vmstat.Pswapout) / interval

	if self.handleMemStat != nil {
		self.handleMemStat(memStat)
	}

	self.beforeMemStat = memStat
	return

}

func (self *StatRunner) syncProcessStat() {
	processes, pidIndexMap, err := GetProcesses("/", true)
	if err != nil {
		return
	}
	if self.beforeProcesses == nil {
		self.beforeProcesses = processes
		self.beforePidIndexMap = pidIndexMap
		return
	}

	interval := self.interval

	for i, process := range processes {
		beforeProcessIndex, ok := self.beforePidIndexMap[process.Pid]
		if !ok {
			continue
		}
		beforeProcess := self.beforeProcesses[beforeProcessIndex]
		stat := process.Stat
		bstat := beforeProcess.Stat

		stat.UserUtil = (stat.Utime - bstat.Utime) / interval
		stat.SystemUtil = (stat.Stime - bstat.Stime) / interval
		stat.GuestUtil = (stat.Gtime - bstat.Gtime) / interval
		stat.CguestUtil = (stat.Cgtime - bstat.Cgtime) / interval

		stat.SchedTimeSlicesPerSec = (stat.SchedTimeSlices - bstat.SchedTimeSlices) / interval
		stat.SchedCpuTimePerSec = (stat.SchedCpuTime - bstat.SchedCpuTime) / interval
		stat.WaitUtil = ((stat.SchedWaitTime - bstat.SchedWaitTime) * self.clkTck) / 1000000000
		stat.VoluntaryCtxtSwitchesPerSec = (stat.VoluntaryCtxtSwitches - bstat.VoluntaryCtxtSwitches) / interval
		stat.NonvoluntaryCtxtSwitches = (stat.NonvoluntaryCtxtSwitches - bstat.NonvoluntaryCtxtSwitches) / interval

		stat.SyscrPerSec = (stat.Syscr - bstat.Syscr) / interval
		stat.SyscwPerSec = (stat.Syscw - bstat.Syscw) / interval
		stat.ReadBytesPerSec = (stat.ReadBytes - bstat.ReadBytes) / interval
		stat.WriteBytesPerSec = (stat.WriteBytes - bstat.WriteBytes) / interval
		processes[i].Stat = stat
	}

	if self.handleProcesses != nil {
		self.handleProcesses(processes)
	}

	self.beforeProcesses = processes
	self.beforePidIndexMap = pidIndexMap
}

func (self *StatRunner) Run(runAt time.Time) {
	self.syncCpuStat()
	self.syncMemStat()
	self.syncProcessStat()
}

func (self *StatRunner) StopTimeout() {
}
