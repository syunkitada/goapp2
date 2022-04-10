package node_ctl

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/syunkitada/goapp2/pkg/lib/process_utils"
	"github.com/syunkitada/goapp2/pkg/lib/runner"
)

var interval int
var stopTimeout int
var isStat bool
var process string
var pid int

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "process",
	Run: func(cmd *cobra.Command, args []string) {
		ecmd := exec.Command("getconf", "CLK_TCK")
		out := new(bytes.Buffer)
		ecmd.Stdout = out
		_ = ecmd.Run()
		clkTck, _ := strconv.Atoi(out.String())
		psRunner := &psRunner{
			ClkTck:        clkTck,
			TargetProcess: process,
			TargetPid:     pid,
		}
		if isStat {
			conf := runner.Config{
				Interval:    interval,
				StopTimeout: stopTimeout,
			}
			runner := runner.New(&conf, psRunner)
			runner.Start()
		} else {
			psRunner.ShowPs()
		}
	},
}

type psRunner struct {
	ClkTck            int
	TargetProcess     string
	TargetPid         int
	beforeProcesses   []process_utils.Process
	beforePidIndexMap map[int]int
}

func (self *psRunner) Run(runAt time.Time) {
	processes, pidIndexMap, err := process_utils.GetProcesses("/", true)
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
		stat.WaitUtil = (stat.SchedWaitTime - bstat.SchedWaitTime) * self.ClkTck / 1000000000

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

func (self *psRunner) ShowPs() {
	processes, _, err := process_utils.GetProcesses("/", true)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Pid", "Name", "Cmd"})
	for _, process := range processes {
		table.Append([]string{strconv.Itoa(process.Pid), process.Name, strings.Join(process.Cmds, " ")})
	}
	table.Render()
}

func (self *psRunner) StopTimeout() {
}

func init() {
	psCmd.PersistentFlags().IntVarP(&interval, "interval", "i", 1, "interval")
	psCmd.PersistentFlags().BoolVarP(&isStat, "stat", "s", false, "stat")
	psCmd.PersistentFlags().IntVarP(&stopTimeout, "stop-timeout", "T", 5, "timeout for stopping process")
	psCmd.PersistentFlags().IntVarP(&pid, "process pid", "p", 0, "timeout for stopping process")
	psCmd.PersistentFlags().StringVarP(&process, "process", "P", "", "timeout for stopping process")
	rootCmd.AddCommand(psCmd)
}
