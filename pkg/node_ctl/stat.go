package node_ctl

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/syunkitada/goapp2/pkg/lib/os_utils"
	"github.com/syunkitada/goapp2/pkg/lib/runner"
)

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "stat",
	Run: func(cmd *cobra.Command, args []string) {
		conf := os_utils.StatControllerConfig{
			Config: runner.Config{
				Interval:    interval,
				StopTimeout: stopTimeout,
			},
			HandleCpuStat: func(cpuStat *os_utils.CpuStat) {
				return
				fmt.Println("DEBUG cpu",
					cpuStat.ProcsRunning,
					cpuStat.ProcsBlocked,
					cpuStat.IntrPerSec,
					cpuStat.CtxPerSec,
					cpuStat.BtimePerSec,
					cpuStat.ProcessesPerSec,
					cpuStat.SoftirqPerSec,
				)
			},
			HandleMemStat: func(memStat *os_utils.MemStat) {
				fmt.Println("DEBUG memstat", memStat)
			},
			HandleProcesses: func(processes []os_utils.Process) {
				return
				for _, p := range processes {
					if pid != 0 && p.Pid != pid {
						continue
					}
					if process != "" && !strings.Contains(p.Name, process) {
						continue
					}
					fmt.Println(strconv.Itoa(p.Pid), p.Name, strconv.Itoa(p.Stat.UserUtil), strconv.Itoa(p.Stat.WaitUtil))
				}

				if !isStat {
					os.Exit(0)
				}
			},
		}
		statCtl := os_utils.NewStatController(&conf)
		statCtl.Start()
	},
}

func init() {
	statCmd.PersistentFlags().IntVarP(&interval, "interval", "i", 1, "interval")
	statCmd.PersistentFlags().BoolVarP(&isStat, "stat", "s", false, "stat")
	statCmd.PersistentFlags().IntVarP(&stopTimeout, "stop-timeout", "T", 5, "timeout for stopping process")
	statCmd.PersistentFlags().IntVarP(&pid, "process pid", "p", 0, "timeout for stopping process")
	statCmd.PersistentFlags().StringVarP(&process, "process", "P", "", "timeout for stopping process")

	rootCmd.AddCommand(statCmd)
}
