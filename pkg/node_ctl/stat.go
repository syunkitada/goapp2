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
					for name, stat := range stats.NetStat.NetDevStatMap {
						// TODO optionでフィルタリングを制御できるようにする
						strs := []string{
							"net:",
							"dev=" + name,
							"rbps=" + strconv.Itoa(stat.ReceiveBytesPerSec),
							"rpps=" + strconv.Itoa(stat.ReceivePacketsPerSec),
							"reps=" + strconv.Itoa(stat.ReceiveErrorsPerSec),
							"rdps=" + strconv.Itoa(stat.ReceiveDropsPerSec),
							"tbps=" + strconv.Itoa(stat.TransmitBytesPerSec),
							"tpps=" + strconv.Itoa(stat.TransmitPacketsPerSec),
							"teps=" + strconv.Itoa(stat.TransmitErrorsPerSec),
							"tdps=" + strconv.Itoa(stat.TransmitDropsPerSec),
						}
						fmt.Println(strings.Join(strs, " "))
					}
					strs := []string{"tcpExt:"}
					if stats.NetStat.TcpExtStat.SyncookiesSentPerSec != 0 {
						strs = append(strs, "SyncookiesSent="+strconv.Itoa(stats.NetStat.TcpExtStat.SyncookiesSentPerSec))
					}
					if stats.NetStat.TcpExtStat.SyncookiesRecvPerSec != 0 {
						strs = append(strs, "SyncookiesRecv="+strconv.Itoa(stats.NetStat.TcpExtStat.SyncookiesRecvPerSec))
					}
					if stats.NetStat.TcpExtStat.SyncookiesFailedPerSec != 0 {
						strs = append(strs, "SyncookiesFailed="+strconv.Itoa(stats.NetStat.TcpExtStat.SyncookiesFailedPerSec))
					}
					if stats.NetStat.TcpExtStat.EmbryonicRstsPerSec != 0 {
						strs = append(strs, "EmbryonicRsts="+strconv.Itoa(stats.NetStat.TcpExtStat.EmbryonicRstsPerSec))
					}
					if stats.NetStat.TcpExtStat.PruneCalledPerSec != 0 {
						strs = append(strs, "PruneCalled="+strconv.Itoa(stats.NetStat.TcpExtStat.PruneCalledPerSec))
					}
					if stats.NetStat.TcpExtStat.RcvPrunedPerSec != 0 {
						strs = append(strs, "RcvPruned="+strconv.Itoa(stats.NetStat.TcpExtStat.RcvPrunedPerSec))
					}
					if stats.NetStat.TcpExtStat.OfoPrunedPerSec != 0 {
						strs = append(strs, "OfoPruned="+strconv.Itoa(stats.NetStat.TcpExtStat.OfoPrunedPerSec))
					}
					if stats.NetStat.TcpExtStat.OutOfWindowIcmpsPerSec != 0 {
						strs = append(strs, "OutOfWindowIcmps="+strconv.Itoa(stats.NetStat.TcpExtStat.OutOfWindowIcmpsPerSec))
					}
					if stats.NetStat.TcpExtStat.LockDroppedIcmpsPerSec != 0 {
						strs = append(strs, "LockDroppedIcmps="+strconv.Itoa(stats.NetStat.TcpExtStat.LockDroppedIcmpsPerSec))
					}
					if stats.NetStat.TcpExtStat.ArpFilterPerSec != 0 {
						strs = append(strs, "ArpFilter="+strconv.Itoa(stats.NetStat.TcpExtStat.ArpFilterPerSec))
					}
					if stats.NetStat.TcpExtStat.TwPerSec != 0 {
						strs = append(strs, "Tw="+strconv.Itoa(stats.NetStat.TcpExtStat.TwPerSec))
					}
					if stats.NetStat.TcpExtStat.TwRecycledPerSec != 0 {
						strs = append(strs, "TwRecycled="+strconv.Itoa(stats.NetStat.TcpExtStat.TwRecycledPerSec))
					}
					if stats.NetStat.TcpExtStat.TwKilledPerSec != 0 {
						strs = append(strs, "TwKilled="+strconv.Itoa(stats.NetStat.TcpExtStat.TwKilledPerSec))
					}
					if stats.NetStat.TcpExtStat.PawsActivePerSec != 0 {
						strs = append(strs, "PawsActive="+strconv.Itoa(stats.NetStat.TcpExtStat.PawsActivePerSec))
					}
					if stats.NetStat.TcpExtStat.PawsEstabPerSec != 0 {
						strs = append(strs, "PawsEstab="+strconv.Itoa(stats.NetStat.TcpExtStat.PawsEstabPerSec))
					}
					if stats.NetStat.TcpExtStat.DelayedAcksPerSec != 0 {
						strs = append(strs, "DelayedAcks="+strconv.Itoa(stats.NetStat.TcpExtStat.DelayedAcksPerSec))
					}
					if stats.NetStat.TcpExtStat.DelayedAckLockedPerSec != 0 {
						strs = append(strs, "DelayedAckLocked="+strconv.Itoa(stats.NetStat.TcpExtStat.DelayedAckLockedPerSec))
					}
					if stats.NetStat.TcpExtStat.DelayedAckLostPerSec != 0 {
						strs = append(strs, "DelayedAckLost="+strconv.Itoa(stats.NetStat.TcpExtStat.DelayedAckLostPerSec))
					}
					if stats.NetStat.TcpExtStat.ListenOverflowsPerSec != 0 {
						strs = append(strs, "ListenOverflows="+strconv.Itoa(stats.NetStat.TcpExtStat.ListenOverflowsPerSec))
					}
					if stats.NetStat.TcpExtStat.ListenDropsPerSec != 0 {
						strs = append(strs, "ListenDrops="+strconv.Itoa(stats.NetStat.TcpExtStat.ListenDropsPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpHpHitsPerSec != 0 {
						strs = append(strs, "TcpHpHits="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpHpHitsPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpPureAcksPerSec != 0 {
						strs = append(strs, "TcpPureAcks="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpPureAcksPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpHpAcksPerSec != 0 {
						strs = append(strs, "TcpHpAcks="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpHpAcksPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpRenoRecoveryPerSec != 0 {
						strs = append(strs, "TcpRenoRecovery="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpRenoRecoveryPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpSackRecoveryPerSec != 0 {
						strs = append(strs, "TcpSackRecovery="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpSackRecoveryPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpSackRenegingPerSec != 0 {
						strs = append(strs, "TcpSackReneging="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpSackRenegingPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpSackReorderPerSec != 0 {
						strs = append(strs, "TcpSackReorder="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpSackReorderPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpRenoReorderPerSec != 0 {
						strs = append(strs, "TcpRenoReorder="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpRenoReorderPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpTsReorderPerSec != 0 {
						strs = append(strs, "TcpTsReorder="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpTsReorderPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpFullUndoPerSec != 0 {
						strs = append(strs, "TcpFullUndo="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpFullUndoPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpPartialUndoPerSec != 0 {
						strs = append(strs, "TcpPartialUndo="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpPartialUndoPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpDsackUndoPerSec != 0 {
						strs = append(strs, "TcpDsackUndo="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpDsackUndoPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpLossUndoPerSec != 0 {
						strs = append(strs, "TcpLossUndo="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpLossUndoPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpLostRetransmitPerSec != 0 {
						strs = append(strs, "TcpLostRetransmit="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpLostRetransmitPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpRenoFailuresPerSec != 0 {
						strs = append(strs, "TcpRenoFailures="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpRenoFailuresPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpSackFailuresPerSec != 0 {
						strs = append(strs, "TcpSackFailures="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpSackFailuresPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpLossFailuresPerSec != 0 {
						strs = append(strs, "TcpLossFailures="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpLossFailuresPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpFastRetransPerSec != 0 {
						strs = append(strs, "TcpFastRetrans="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpFastRetransPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpSlowStartRetransPerSec != 0 {
						strs = append(strs, "TcpSlowStartRetrans="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpSlowStartRetransPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpTimeoutsPerSec != 0 {
						strs = append(strs, "TcpTimeouts="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpTimeoutsPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpLossProbesPerSec != 0 {
						strs = append(strs, "TcpLossProbes="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpLossProbesPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpLossProbeRecoveryPerSec != 0 {
						strs = append(strs, "TcpLossProbeRecovery="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpLossProbeRecoveryPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpRenoRecoveryFailPerSec != 0 {
						strs = append(strs, "TcpRenoRecoveryFail="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpRenoRecoveryFailPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpSackRecoveryFailPerSec != 0 {
						strs = append(strs, "TcpSackRecoveryFail="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpSackRecoveryFailPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpRcvCollapsedPerSec != 0 {
						strs = append(strs, "TcpRcvCollapsed="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpRcvCollapsedPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpBacklogCoalescePerSec != 0 {
						strs = append(strs, "TcpBacklogCoalesce="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpBacklogCoalescePerSec))
					}
					if stats.NetStat.TcpExtStat.TcpDsackOldSentPerSec != 0 {
						strs = append(strs, "TcpDsackOldSent="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpDsackOldSentPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpDsackOfoSentPerSec != 0 {
						strs = append(strs, "TcpDsackOfoSent="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpDsackOfoSentPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpDsackRecvPerSec != 0 {
						strs = append(strs, "TcpDsackRecv="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpDsackRecvPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpDsackOfoRecvPerSec != 0 {
						strs = append(strs, "TcpDsackOfoRecv="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpDsackOfoRecvPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpAbortOnDataPerSec != 0 {
						strs = append(strs, "TcpAbortOnData="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpAbortOnDataPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpAbortOnClosePerSec != 0 {
						strs = append(strs, "TcpAbortOnClose="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpAbortOnClosePerSec))
					}
					if stats.NetStat.TcpExtStat.TcpAbortOnMemoryPerSec != 0 {
						strs = append(strs, "TcpAbortOnMemory="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpAbortOnMemoryPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpAbortOnTimeoutPerSec != 0 {
						strs = append(strs, "TcpAbortOnTimeout="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpAbortOnTimeoutPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpAbortOnLingerPerSec != 0 {
						strs = append(strs, "TcpAbortOnLinger="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpAbortOnLingerPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpAbortFailedPerSec != 0 {
						strs = append(strs, "TcpAbortFailed="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpAbortFailedPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpMemoryPressuresPerSec != 0 {
						strs = append(strs, "TcpMemoryPressures="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpMemoryPressuresPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpMemoryPressuresChronoPerSec != 0 {
						strs = append(strs, "TcpMemoryPressuresChrono="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpMemoryPressuresChronoPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpSackDiscardPerSec != 0 {
						strs = append(strs, "TcpSackDiscard="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpSackDiscardPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpDsackIgnoredOldPerSec != 0 {
						strs = append(strs, "TcpDsackIgnoredOld="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpDsackIgnoredOldPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpDsackIgnoredNoUndoPerSec != 0 {
						strs = append(strs, "TcpDsackIgnoredNoUndo="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpDsackIgnoredNoUndoPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpSpuriousRTOsPerSec != 0 {
						strs = append(strs, "TcpSpuriousRTOs="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpSpuriousRTOsPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpMd5NotFoundPerSec != 0 {
						strs = append(strs, "TcpMd5NotFound="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpMd5NotFoundPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpMd5UnexpectedPerSec != 0 {
						strs = append(strs, "TcpMd5Unexpected="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpMd5UnexpectedPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpMd5FailurePerSec != 0 {
						strs = append(strs, "TcpMd5Failure="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpMd5FailurePerSec))
					}
					if stats.NetStat.TcpExtStat.TcpSackShiftedPerSec != 0 {
						strs = append(strs, "TcpSackShifted="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpSackShiftedPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpSackMergedPerSec != 0 {
						strs = append(strs, "TcpSackMerged="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpSackMergedPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpSackShiftFallbackPerSec != 0 {
						strs = append(strs, "TcpSackShiftFallback="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpSackShiftFallbackPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpBacklogDropPerSec != 0 {
						strs = append(strs, "TcpBacklogDrop="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpBacklogDropPerSec))
					}
					if stats.NetStat.TcpExtStat.PfMemallocDropPerSec != 0 {
						strs = append(strs, "PfMemallocDrop="+strconv.Itoa(stats.NetStat.TcpExtStat.PfMemallocDropPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpMinTtlDropPerSec != 0 {
						strs = append(strs, "TcpMinTtlDrop="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpMinTtlDropPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpDeferAcceptDropPerSec != 0 {
						strs = append(strs, "TcpDeferAcceptDrop="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpDeferAcceptDropPerSec))
					}
					if stats.NetStat.TcpExtStat.IpReversePathFilterPerSec != 0 {
						strs = append(strs, "IpReversePathFilter="+strconv.Itoa(stats.NetStat.TcpExtStat.IpReversePathFilterPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpTimeWaitOverflowPerSec != 0 {
						strs = append(strs, "TcpTimeWaitOverflow="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpTimeWaitOverflowPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpReqQFullDoCookiesPerSec != 0 {
						strs = append(strs, "TcpReqQFullDoCookies="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpReqQFullDoCookiesPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpReqQFullDropPerSec != 0 {
						strs = append(strs, "TcpReqQFullDrop="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpReqQFullDropPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpRetransFailPerSec != 0 {
						strs = append(strs, "TcpRetransFail="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpRetransFailPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpRcvCoalescePerSec != 0 {
						strs = append(strs, "TcpRcvCoalesce="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpRcvCoalescePerSec))
					}
					if stats.NetStat.TcpExtStat.TcpOfoQueuePerSec != 0 {
						strs = append(strs, "TcpOfoQueue="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpOfoQueuePerSec))
					}
					if stats.NetStat.TcpExtStat.TcpOfoDropPerSec != 0 {
						strs = append(strs, "TcpOfoDrop="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpOfoDropPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpOfoMergePerSec != 0 {
						strs = append(strs, "TcpOfoMerge="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpOfoMergePerSec))
					}
					if stats.NetStat.TcpExtStat.TcpChallengeACKPerSec != 0 {
						strs = append(strs, "TcpChallengeACK="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpChallengeACKPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpSynChallengePerSec != 0 {
						strs = append(strs, "TcpSynChallenge="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpSynChallengePerSec))
					}
					if stats.NetStat.TcpExtStat.TcpFastOpenActivePerSec != 0 {
						strs = append(strs, "TcpFastOpenActive="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpFastOpenActivePerSec))
					}
					if stats.NetStat.TcpExtStat.TcpFastOpenActiveFailPerSec != 0 {
						strs = append(strs, "TcpFastOpenActiveFail="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpFastOpenActiveFailPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpFastOpenPassivePerSec != 0 {
						strs = append(strs, "TcpFastOpenPassive="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpFastOpenPassivePerSec))
					}
					if stats.NetStat.TcpExtStat.TcpFastOpenPassiveFailPerSec != 0 {
						strs = append(strs, "TcpFastOpenPassiveFail="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpFastOpenPassiveFailPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpFastOpenListenOverflowPerSec != 0 {
						strs = append(strs, "TcpFastOpenListenOverflow="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpFastOpenListenOverflowPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpFastOpenCookieReqdPerSec != 0 {
						strs = append(strs, "TcpFastOpenCookieReqd="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpFastOpenCookieReqdPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpFastOpenBlackholePerSec != 0 {
						strs = append(strs, "TcpFastOpenBlackhole="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpFastOpenBlackholePerSec))
					}
					if stats.NetStat.TcpExtStat.TcpSpuriousRtxHostQueuesPerSec != 0 {
						strs = append(strs, "TcpSpuriousRtxHostQueues="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpSpuriousRtxHostQueuesPerSec))
					}
					if stats.NetStat.TcpExtStat.BusyPollRxPacketsPerSec != 0 {
						strs = append(strs, "BusyPollRxPackets="+strconv.Itoa(stats.NetStat.TcpExtStat.BusyPollRxPacketsPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpAutoCorkingPerSec != 0 {
						strs = append(strs, "TcpAutoCorking="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpAutoCorkingPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpFromZeroWindowAdvPerSec != 0 {
						strs = append(strs, "TcpFromZeroWindowAdv="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpFromZeroWindowAdvPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpToZeroWindowAdvPerSec != 0 {
						strs = append(strs, "TcpToZeroWindowAdv="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpToZeroWindowAdvPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpWantZeroWindowAdvPerSec != 0 {
						strs = append(strs, "TcpWantZeroWindowAdv="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpWantZeroWindowAdvPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpSynRetransPerSec != 0 {
						strs = append(strs, "TcpSynRetrans="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpSynRetransPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpOrigDataSentPerSec != 0 {
						strs = append(strs, "TcpOrigDataSent="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpOrigDataSentPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpHystartTrainDetectPerSec != 0 {
						strs = append(strs, "TcpHystartTrainDetect="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpHystartTrainDetectPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpHystartTrainCwndPerSec != 0 {
						strs = append(strs, "TcpHystartTrainCwnd="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpHystartTrainCwndPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpHystartDelayDetectPerSec != 0 {
						strs = append(strs, "TcpHystartDelayDetect="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpHystartDelayDetectPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpHystartDelayCwndPerSec != 0 {
						strs = append(strs, "TcpHystartDelayCwnd="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpHystartDelayCwndPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpAckSkippedSynRecvPerSec != 0 {
						strs = append(strs, "TcpAckSkippedSynRecv="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpAckSkippedSynRecvPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpAckSkippedPAWSPerSec != 0 {
						strs = append(strs, "TcpAckSkippedPAWS="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpAckSkippedPAWSPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpAckSkippedSeqPerSec != 0 {
						strs = append(strs, "TcpAckSkippedSeq="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpAckSkippedSeqPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpAckSkippedFinWait2PerSec != 0 {
						strs = append(strs, "TcpAckSkippedFinWait2="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpAckSkippedFinWait2PerSec))
					}
					if stats.NetStat.TcpExtStat.TcpAckSkippedTimeWaitPerSec != 0 {
						strs = append(strs, "TcpAckSkippedTimeWait="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpAckSkippedTimeWaitPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpAckSkippedChallengePerSec != 0 {
						strs = append(strs, "TcpAckSkippedChallenge="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpAckSkippedChallengePerSec))
					}
					if stats.NetStat.TcpExtStat.TcpWinProbePerSec != 0 {
						strs = append(strs, "TcpWinProbe="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpWinProbePerSec))
					}
					if stats.NetStat.TcpExtStat.TcpKeepAlivePerSec != 0 {
						strs = append(strs, "TcpKeepAlive="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpKeepAlivePerSec))
					}
					if stats.NetStat.TcpExtStat.TcpMtupFailPerSec != 0 {
						strs = append(strs, "TcpMtupFail="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpMtupFailPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpMtupSuccessPerSec != 0 {
						strs = append(strs, "TcpMtupSuccess="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpMtupSuccessPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpDeliveredPerSec != 0 {
						strs = append(strs, "TcpDelivered="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpDeliveredPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpDeliveredCEPerSec != 0 {
						strs = append(strs, "TcpDeliveredCE="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpDeliveredCEPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpAckCompressedPerSec != 0 {
						strs = append(strs, "TcpAckCompressed="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpAckCompressedPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpZeroWindowDropPerSec != 0 {
						strs = append(strs, "TcpZeroWindowDrop="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpZeroWindowDropPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpRcvQDropPerSec != 0 {
						strs = append(strs, "TcpRcvQDrop="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpRcvQDropPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpWqueueTooBigPerSec != 0 {
						strs = append(strs, "TcpWqueueTooBig="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpWqueueTooBigPerSec))
					}
					if stats.NetStat.TcpExtStat.TcpFastOpenPassiveAltKeyPerSec != 0 {
						strs = append(strs, "TcpFastOpenPassiveAltKey="+strconv.Itoa(stats.NetStat.TcpExtStat.TcpFastOpenPassiveAltKeyPerSec))
					}

					fmt.Println(strings.Join(strs, " "))

					strs = []string{"ipExt:"}

					if stats.NetStat.IpExtStat.InNoRoutesPerSec != 0 {
						strs = append(strs, "InNoRoutes="+strconv.Itoa(stats.NetStat.IpExtStat.InNoRoutesPerSec))
					}
					if stats.NetStat.IpExtStat.InTruncatedPktsPerSec != 0 {
						strs = append(strs, "InTruncatedPkts="+strconv.Itoa(stats.NetStat.IpExtStat.InTruncatedPktsPerSec))
					}
					if stats.NetStat.IpExtStat.InCsumErrorsPerSec != 0 {
						strs = append(strs, "InCsumErrors="+strconv.Itoa(stats.NetStat.IpExtStat.InCsumErrorsPerSec))
					}
					if stats.NetStat.IpExtStat.InNoRoutesPerSec != 0 {
						strs = append(strs, "InNoRoutes="+strconv.Itoa(stats.NetStat.IpExtStat.InNoRoutesPerSec))
					}
					if stats.NetStat.IpExtStat.InTruncatedPktsPerSec != 0 {
						strs = append(strs, "InTruncatedPkts="+strconv.Itoa(stats.NetStat.IpExtStat.InTruncatedPktsPerSec))
					}
					if stats.NetStat.IpExtStat.InMcastPktsPerSec != 0 {
						strs = append(strs, "InMcastPkts="+strconv.Itoa(stats.NetStat.IpExtStat.InMcastPktsPerSec))
					}
					if stats.NetStat.IpExtStat.OutMcastPktsPerSec != 0 {
						strs = append(strs, "OutMcastPkts="+strconv.Itoa(stats.NetStat.IpExtStat.OutMcastPktsPerSec))
					}
					if stats.NetStat.IpExtStat.InBcastPktsPerSec != 0 {
						strs = append(strs, "InBcastPkts="+strconv.Itoa(stats.NetStat.IpExtStat.InBcastPktsPerSec))
					}
					if stats.NetStat.IpExtStat.OutBcastPktsPerSec != 0 {
						strs = append(strs, "OutBcastPkts="+strconv.Itoa(stats.NetStat.IpExtStat.OutBcastPktsPerSec))
					}
					if stats.NetStat.IpExtStat.InOctetsPerSec != 0 {
						strs = append(strs, "InOctets="+strconv.Itoa(stats.NetStat.IpExtStat.InOctetsPerSec))
					}
					if stats.NetStat.IpExtStat.OutOctetsPerSec != 0 {
						strs = append(strs, "OutOctets="+strconv.Itoa(stats.NetStat.IpExtStat.OutOctetsPerSec))
					}
					if stats.NetStat.IpExtStat.InMcastOctetsPerSec != 0 {
						strs = append(strs, "InMcastOctets="+strconv.Itoa(stats.NetStat.IpExtStat.InMcastOctetsPerSec))
					}
					if stats.NetStat.IpExtStat.OutMcastOctetsPerSec != 0 {
						strs = append(strs, "OutMcastOctets="+strconv.Itoa(stats.NetStat.IpExtStat.OutMcastOctetsPerSec))
					}
					if stats.NetStat.IpExtStat.InBcastOctetsPerSec != 0 {
						strs = append(strs, "InBcastOctets="+strconv.Itoa(stats.NetStat.IpExtStat.InBcastOctetsPerSec))
					}
					if stats.NetStat.IpExtStat.OutBcastOctetsPerSec != 0 {
						strs = append(strs, "OutBcastOctets="+strconv.Itoa(stats.NetStat.IpExtStat.OutBcastOctetsPerSec))
					}
					if stats.NetStat.IpExtStat.InCsumErrorsPerSec != 0 {
						strs = append(strs, "InCsumErrors="+strconv.Itoa(stats.NetStat.IpExtStat.InCsumErrorsPerSec))
					}
					if stats.NetStat.IpExtStat.InNoECTPktsPerSec != 0 {
						strs = append(strs, "InNoECTPkts="+strconv.Itoa(stats.NetStat.IpExtStat.InNoECTPktsPerSec))
					}
					if stats.NetStat.IpExtStat.InECT1PktsPerSec != 0 {
						strs = append(strs, "InECT1Pkts="+strconv.Itoa(stats.NetStat.IpExtStat.InECT1PktsPerSec))
					}
					if stats.NetStat.IpExtStat.InECT0PktsPerSec != 0 {
						strs = append(strs, "InECT0Pkts="+strconv.Itoa(stats.NetStat.IpExtStat.InECT0PktsPerSec))
					}
					if stats.NetStat.IpExtStat.InCEPktsPerSec != 0 {
						strs = append(strs, "InCEPkts="+strconv.Itoa(stats.NetStat.IpExtStat.InCEPktsPerSec))
					}
					if stats.NetStat.IpExtStat.ReasmOverlapsPerSec != 0 {
						strs = append(strs, "ReasmOverlaps="+strconv.Itoa(stats.NetStat.IpExtStat.ReasmOverlapsPerSec))
					}

					fmt.Println(strings.Join(strs, " "))

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
