package node_ctl

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/syunkitada/goapp2/pkg/lib/os_utils"
	"github.com/syunkitada/goapp2/pkg/lib/runner"
)

var target string
var interval int
var stopTimeout int
var isStat bool
var process string
var pid int

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "stat",
	Run: func(cmd *cobra.Command, args []string) {
		showCpu := strings.Contains(target, "c")
		showCpuWide := strings.Contains(target, "C")
		showMem := strings.Contains(target, "m")
		showMemWide := strings.Contains(target, "m")
		showBuddyinfo := strings.Contains(target, "b")
		showDisk := strings.Contains(target, "d")
		showDiskWide := strings.Contains(target, "D")
		showFs := strings.Contains(target, "f")
		showNet := strings.Contains(target, "n")
		showUser := strings.Contains(target, "u")
		// TODO 異常値をカラーリングできるようにする（設定値はファイルからも読み取れるようにする）

		conf := os_utils.StatControllerConfig{
			Config: runner.Config{
				Interval:    interval,
				StopTimeout: stopTimeout,
			},
			HandleStats: func(runAt time.Time, stats *os_utils.Stats) {
				fmt.Println("time:", runAt)
				strs := []string{}
				if showCpu || showCpuWide {
					strs = append(strs,
						"cpu:",
						"run="+strconv.Itoa(stats.CpuStat.ProcsRunning),
						"blocked="+strconv.Itoa(stats.CpuStat.ProcsBlocked),
					)
					if showCpuWide {
						strs = append(strs,
							"intr="+strconv.Itoa(stats.CpuStat.IntrPerSec),
							"ctx="+strconv.Itoa(stats.CpuStat.CtxPerSec),
							"btime="+strconv.Itoa(stats.CpuStat.BtimePerSec),
							"process="+strconv.Itoa(stats.CpuStat.ProcessesPerSec),
							"sirq="+strconv.Itoa(stats.CpuStat.SoftirqPerSec),
						)
					}
				}
				fmt.Println(strings.Join(strs, " "))

				if showMem || showMemWide {
					for _, node := range stats.MemStat.Nodes {
						strs := []string{
							"mem:",
							"node=" + strconv.Itoa(node.NodeId),
							"tota=" + strconv.Itoa(node.MemTotal),
							"free=" + strconv.Itoa(node.MemFree),
							"used=" + strconv.Itoa(node.MemUsed),
							"avai=" + strconv.Itoa(node.MemAvailable),
						}
						fmt.Println(strings.Join(strs, " "))
						if showBuddyinfo {
							strs := []string{
								"buddyinfo:",
								"node=" + strconv.Itoa(node.NodeId),
								"4k=" + strconv.Itoa(node.Buddyinfo.M4K),
								"8k=" + strconv.Itoa(node.Buddyinfo.M8K),
								"16k=" + strconv.Itoa(node.Buddyinfo.M16K),
								"32k=" + strconv.Itoa(node.Buddyinfo.M32K),
								"64k=" + strconv.Itoa(node.Buddyinfo.M64K),
								"128k=" + strconv.Itoa(node.Buddyinfo.M128K),
								"256k=" + strconv.Itoa(node.Buddyinfo.M256K),
								"512k=" + strconv.Itoa(node.Buddyinfo.M512K),
								"1m=" + strconv.Itoa(node.Buddyinfo.M1M),
								"2m=" + strconv.Itoa(node.Buddyinfo.M2M),
								"4m=" + strconv.Itoa(node.Buddyinfo.M4M),
							}
							fmt.Println(strings.Join(strs, " "))
						}
					}
				}

				if showDisk || showDiskWide {
					for name, stat := range stats.DiskStat.DiskDeviceStatMap {
						// TODO FIME optionでフィルタリングを制御できるようにする
						if strings.Contains(name, "loop") {
							continue
						}
						strs := []string{
							"disk:",
							"device=" + name,
							"rps=" + strconv.Itoa(stat.ReadsPerSec),
							"rbps=" + strconv.Itoa(stat.ReadBytesPerSec),
							"rmsps=" + strconv.Itoa(stat.ReadMsPerSec),
							"wps=" + strconv.Itoa(stat.WritesPerSec),
							"wbps=" + strconv.Itoa(stat.WriteBytesPerSec),
							"wmsps=" + strconv.Itoa(stat.WriteMsPerSec),
							"pios=" + strconv.Itoa(stat.ProgressIos),
						}
						fmt.Println(strings.Join(strs, " "))
					}
				}
				if showFs {
					for name, stat := range stats.DiskStat.DiskFsStatMap {
						// TODO FIME optionでフィルタリングを制御できるようにする
						if strings.Contains(name, "loop") {
							continue
						}
						if !strings.Contains(stat.Type, "ext") {
							continue
						}
						strs := []string{
							"fs:",
							"path=" + name,
							"mount=" + stat.MountPath,
							"type=" + stat.Type,
							"total=" + strconv.Itoa(stat.TotalSize),
							"free=" + strconv.Itoa(stat.FreeSize),
							"used=" + strconv.Itoa(stat.UsedSize),
							"files=" + strconv.Itoa(stat.Files),
						}
						fmt.Println(strings.Join(strs, " "))
					}
				}

				if showNet {
					// TODO
				}

				if showUser {
					for name, stat := range stats.LoginUserStat.UserStatMap {
						strs := []string{
							"user:",
							"name=" + name,
							"durationSec=" + strconv.Itoa(stat.LoginDuration),
						}
						fmt.Println(strings.Join(strs, " "))
					}
				}

				if pid != 0 {
					for _, p := range stats.Processes {
						if p.Pid != pid {
							continue
						}
						fmt.Println(strconv.Itoa(p.Pid), p.Name, strconv.Itoa(p.Stat.UserUtil), strconv.Itoa(p.Stat.WaitUtil))
					}
				}
				if process != "" {
					for _, p := range stats.Processes {
						if !strings.Contains(p.Name, process) {
							continue
						}
						fmt.Println(strconv.Itoa(p.Pid), p.Name, strconv.Itoa(p.Stat.UserUtil), strconv.Itoa(p.Stat.WaitUtil))
					}
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
	statCmd.PersistentFlags().StringVarP(&target, "target", "t", "", "stat target")

	rootCmd.AddCommand(statCmd)
}
