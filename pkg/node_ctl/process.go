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

var interval int
var stopTimeout int
var isStat bool
var process string
var pid int

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "process",
	Run: func(cmd *cobra.Command, args []string) {
		conf := os_utils.StatControllerConfig{
			Config: runner.Config{
				Interval:    interval,
				StopTimeout: stopTimeout,
			},
			HandleProcesses: func(processes []os_utils.Process) {
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
	psCmd.PersistentFlags().IntVarP(&interval, "interval", "i", 1, "interval")
	psCmd.PersistentFlags().BoolVarP(&isStat, "stat", "s", false, "stat")
	psCmd.PersistentFlags().IntVarP(&stopTimeout, "stop-timeout", "T", 5, "timeout for stopping process")
	psCmd.PersistentFlags().IntVarP(&pid, "process pid", "p", 0, "timeout for stopping process")
	psCmd.PersistentFlags().StringVarP(&process, "process", "P", "", "timeout for stopping process")
	rootCmd.AddCommand(psCmd)
}
