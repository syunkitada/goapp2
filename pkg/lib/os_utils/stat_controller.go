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
	HandleStats func(runAt time.Time, stats *Stats)
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
		clkTck:      clkTck,
		handleStats: conf.HandleStats,
		interval:    conf.Config.Interval,
	}
	statController = &StatController{
		Runner:     *runner.New(&conf.Config, &statRunner),
		statRunner: &statRunner,
	}
	return
}

type StatRunner struct {
	clkTck               int
	interval             int
	handleStats          func(runAt time.Time, stats *Stats)
	currentCpuStat       *CpuStat
	currentMemStat       *MemStat
	currentDiskStat      *DiskStat
	currentNetStat       *NetStat
	currentLoginUserStat *LoginUserStat
	currentUptimeStat    *UptimeStat
	currentProcesses     []Process
	currentPidIndexMap   map[int]int
	currentStats         *Stats
}

type Stats struct {
	CpuStat       *CpuStat
	MemStat       *MemStat
	DiskStat      *DiskStat
	NetStat       *NetStat
	Processes     []Process
	LoginUserStat *LoginUserStat
	UptimeStat    *UptimeStat
}

func (self *StatRunner) syncCpuStat() {
	cpuStat, err := GetCpuStat()
	if err != nil {
		return
	}

	if self.currentCpuStat == nil {
		self.currentCpuStat = cpuStat
		return
	}

	interval := self.interval

	cpuStat.IntrPerSec = (cpuStat.Intr - self.currentCpuStat.Intr) / interval
	cpuStat.CtxPerSec = (cpuStat.Ctx - self.currentCpuStat.Ctx) / interval
	cpuStat.BtimePerSec = (cpuStat.Btime - self.currentCpuStat.Btime) / interval
	cpuStat.ProcessesPerSec = (cpuStat.Processes - self.currentCpuStat.Processes) / interval
	cpuStat.SoftirqPerSec = (cpuStat.Softirq - self.currentCpuStat.Softirq) / interval

	self.currentCpuStat = cpuStat
}

func (self *StatRunner) syncMemStat() {
	var memStat *MemStat
	var err error
	if memStat, err = GetMemStat("/"); err != nil {
		return
	}

	if self.currentMemStat == nil {
		self.currentMemStat = memStat
		return
	}

	interval := self.interval

	memStat.Vmstat.PgscanKswapdPerSec = (memStat.Vmstat.PgscanKswapd - self.currentMemStat.Vmstat.PgscanKswapd) / interval
	memStat.Vmstat.PgscanDirectPerSec = (memStat.Vmstat.PgscanDirect - self.currentMemStat.Vmstat.PgscanDirect) / interval
	memStat.Vmstat.PgfaultPerSec = (memStat.Vmstat.Pgfault - self.currentMemStat.Vmstat.Pgfault) / interval
	memStat.Vmstat.PswapinPerSec = (memStat.Vmstat.Pswapin - self.currentMemStat.Vmstat.Pswapin) / interval
	memStat.Vmstat.PswapoutPerSec = (memStat.Vmstat.Pswapout - self.currentMemStat.Vmstat.Pswapout) / interval

	self.currentMemStat = memStat
	return
}

func (self *StatRunner) syncDiskStat() {
	var diskStat *DiskStat
	var err error
	if diskStat, err = GetDiskStat(); err != nil {
		return
	}

	if self.currentDiskStat == nil {
		self.currentDiskStat = diskStat
		return
	}

	interval := self.interval

	for deviceName, cstat := range diskStat.DiskDeviceStatMap {
		bstat, ok := diskStat.DiskDeviceStatMap[deviceName]
		if !ok {
			continue
		}
		cstat.ReadsPerSec = (cstat.ReadsCompleted - bstat.ReadsCompleted) / interval
		cstat.RmergesPerSec = (cstat.ReadsMerges - bstat.ReadsMerges) / interval
		cstat.ReadBytesPerSec = ((cstat.ReadSectors - bstat.ReadSectors) * cstat.PblockSize) / interval
		cstat.ReadMsPerSec = (cstat.ReadMs - bstat.ReadMs) / interval

		cstat.WritesPerSec = (cstat.WritesCompleted - bstat.WritesCompleted) / interval
		cstat.WmergesPerSec = (cstat.WritesMerges - bstat.WritesMerges) / interval
		cstat.WriteBytesPerSec = ((cstat.WriteSectors - bstat.WriteSectors) * cstat.PblockSize) / interval
		cstat.WriteMsPerSec = (cstat.WriteMs - bstat.WriteMs) / interval

		cstat.DiscardsPerSec = (cstat.DiscardsCompleted - bstat.DiscardsCompleted) / interval
		cstat.DmergesPerSec = (cstat.DiscardsMerges - bstat.DiscardsMerges) / interval
		cstat.DiscardBytesPerSec = ((cstat.DiscardSectors - bstat.DiscardSectors) * cstat.PblockSize) / interval
		cstat.DiscardMsPerSec = (cstat.DiscardMs - bstat.DiscardMs) / interval

		cstat.IosMsPerSec = (cstat.IosMs - bstat.IosMs) / interval
		cstat.WeightedIosMsPerSec = (cstat.WeightedIosMs - bstat.WeightedIosMs) / interval

		diskStat.DiskDeviceStatMap[deviceName] = cstat
	}

	self.currentDiskStat = diskStat
	return
}

func (self *StatRunner) syncNetStat() {
	netStat, err := GetNetStat()
	if err != nil {
		return
	}
	if self.currentNetStat == nil {
		self.currentNetStat = netStat
		return
	}

	interval := self.interval

	for dev, cstat := range netStat.NetDevStatMap {
		bstat, ok := self.currentNetStat.NetDevStatMap[dev]
		if !ok {
			continue
		}
		cstat.ReceiveBytesPerSec = (cstat.ReceiveBytes - bstat.ReceiveBytes) / interval
		cstat.ReceivePacketsPerSec = (cstat.ReceivePackets - bstat.ReceivePackets) / interval
		cstat.ReceiveErrorsPerSec = (cstat.ReceiveErrors - bstat.ReceiveErrors) / interval
		cstat.ReceiveDropsPerSec = (cstat.ReceiveDrops - bstat.ReceiveDrops) / interval
		cstat.TransmitBytesPerSec = (cstat.TransmitBytes - bstat.TransmitBytes) / interval
		cstat.TransmitPacketsPerSec = (cstat.TransmitPackets - bstat.TransmitPackets) / interval
		cstat.TransmitErrorsPerSec = (cstat.TransmitErrors - bstat.TransmitErrors) / interval
		cstat.TransmitDropsPerSec = (cstat.TransmitDrops - bstat.TransmitDrops) / interval

		netStat.NetDevStatMap[dev] = cstat
	}

	netStat.TcpExtStat.SyncookiesSentPerSec = (netStat.TcpExtStat.SyncookiesSent - self.currentNetStat.TcpExtStat.SyncookiesSent) / interval
	netStat.TcpExtStat.SyncookiesRecvPerSec = (netStat.TcpExtStat.SyncookiesRecv - self.currentNetStat.TcpExtStat.SyncookiesRecv) / interval
	netStat.TcpExtStat.SyncookiesFailedPerSec = (netStat.TcpExtStat.SyncookiesFailed - self.currentNetStat.TcpExtStat.SyncookiesFailed) / interval
	netStat.TcpExtStat.EmbryonicRstsPerSec = (netStat.TcpExtStat.EmbryonicRsts - self.currentNetStat.TcpExtStat.EmbryonicRsts) / interval
	netStat.TcpExtStat.PruneCalledPerSec = (netStat.TcpExtStat.PruneCalled - self.currentNetStat.TcpExtStat.PruneCalled) / interval
	netStat.TcpExtStat.RcvPrunedPerSec = (netStat.TcpExtStat.RcvPruned - self.currentNetStat.TcpExtStat.RcvPruned) / interval
	netStat.TcpExtStat.OfoPrunedPerSec = (netStat.TcpExtStat.OfoPruned - self.currentNetStat.TcpExtStat.OfoPruned) / interval
	netStat.TcpExtStat.OutOfWindowIcmpsPerSec = (netStat.TcpExtStat.OutOfWindowIcmps - self.currentNetStat.TcpExtStat.OutOfWindowIcmps) / interval
	netStat.TcpExtStat.LockDroppedIcmpsPerSec = (netStat.TcpExtStat.LockDroppedIcmps - self.currentNetStat.TcpExtStat.LockDroppedIcmps) / interval
	netStat.TcpExtStat.ArpFilterPerSec = (netStat.TcpExtStat.ArpFilter - self.currentNetStat.TcpExtStat.ArpFilter) / interval
	netStat.TcpExtStat.TwPerSec = (netStat.TcpExtStat.Tw - self.currentNetStat.TcpExtStat.Tw) / interval
	netStat.TcpExtStat.TwRecycledPerSec = (netStat.TcpExtStat.TwRecycled - self.currentNetStat.TcpExtStat.TwRecycled) / interval
	netStat.TcpExtStat.TwKilledPerSec = (netStat.TcpExtStat.TwKilled - self.currentNetStat.TcpExtStat.TwKilled) / interval
	netStat.TcpExtStat.PawsActivePerSec = (netStat.TcpExtStat.PawsActive - self.currentNetStat.TcpExtStat.PawsActive) / interval
	netStat.TcpExtStat.PawsEstabPerSec = (netStat.TcpExtStat.PawsEstab - self.currentNetStat.TcpExtStat.PawsEstab) / interval
	netStat.TcpExtStat.DelayedAcksPerSec = (netStat.TcpExtStat.DelayedAcks - self.currentNetStat.TcpExtStat.DelayedAcks) / interval
	netStat.TcpExtStat.DelayedAckLockedPerSec = (netStat.TcpExtStat.DelayedAckLocked - self.currentNetStat.TcpExtStat.DelayedAckLocked) / interval
	netStat.TcpExtStat.DelayedAckLostPerSec = (netStat.TcpExtStat.DelayedAckLost - self.currentNetStat.TcpExtStat.DelayedAckLost) / interval
	netStat.TcpExtStat.ListenOverflowsPerSec = (netStat.TcpExtStat.ListenOverflows - self.currentNetStat.TcpExtStat.ListenOverflows) / interval
	netStat.TcpExtStat.ListenDropsPerSec = (netStat.TcpExtStat.ListenDrops - self.currentNetStat.TcpExtStat.ListenDrops) / interval
	netStat.TcpExtStat.TcpHpHitsPerSec = (netStat.TcpExtStat.TcpHpHits - self.currentNetStat.TcpExtStat.TcpHpHits) / interval
	netStat.TcpExtStat.TcpPureAcksPerSec = (netStat.TcpExtStat.TcpPureAcks - self.currentNetStat.TcpExtStat.TcpPureAcks) / interval
	netStat.TcpExtStat.TcpHpAcksPerSec = (netStat.TcpExtStat.TcpHpAcks - self.currentNetStat.TcpExtStat.TcpHpAcks) / interval
	netStat.TcpExtStat.TcpRenoRecoveryPerSec = (netStat.TcpExtStat.TcpRenoRecovery - self.currentNetStat.TcpExtStat.TcpRenoRecovery) / interval
	netStat.TcpExtStat.TcpSackRecoveryPerSec = (netStat.TcpExtStat.TcpSackRecovery - self.currentNetStat.TcpExtStat.TcpSackRecovery) / interval
	netStat.TcpExtStat.TcpSackRenegingPerSec = (netStat.TcpExtStat.TcpSackReneging - self.currentNetStat.TcpExtStat.TcpSackReneging) / interval
	netStat.TcpExtStat.TcpSackReorderPerSec = (netStat.TcpExtStat.TcpSackReorder - self.currentNetStat.TcpExtStat.TcpSackReorder) / interval
	netStat.TcpExtStat.TcpRenoReorderPerSec = (netStat.TcpExtStat.TcpRenoReorder - self.currentNetStat.TcpExtStat.TcpRenoReorder) / interval
	netStat.TcpExtStat.TcpTsReorderPerSec = (netStat.TcpExtStat.TcpTsReorder - self.currentNetStat.TcpExtStat.TcpTsReorder) / interval
	netStat.TcpExtStat.TcpFullUndoPerSec = (netStat.TcpExtStat.TcpFullUndo - self.currentNetStat.TcpExtStat.TcpFullUndo) / interval
	netStat.TcpExtStat.TcpPartialUndoPerSec = (netStat.TcpExtStat.TcpPartialUndo - self.currentNetStat.TcpExtStat.TcpPartialUndo) / interval
	netStat.TcpExtStat.TcpDsackUndoPerSec = (netStat.TcpExtStat.TcpDsackUndo - self.currentNetStat.TcpExtStat.TcpDsackUndo) / interval
	netStat.TcpExtStat.TcpLossUndoPerSec = (netStat.TcpExtStat.TcpLossUndo - self.currentNetStat.TcpExtStat.TcpLossUndo) / interval
	netStat.TcpExtStat.TcpLostRetransmitPerSec = (netStat.TcpExtStat.TcpLostRetransmit - self.currentNetStat.TcpExtStat.TcpLostRetransmit) / interval
	netStat.TcpExtStat.TcpRenoFailuresPerSec = (netStat.TcpExtStat.TcpRenoFailures - self.currentNetStat.TcpExtStat.TcpRenoFailures) / interval
	netStat.TcpExtStat.TcpSackFailuresPerSec = (netStat.TcpExtStat.TcpSackFailures - self.currentNetStat.TcpExtStat.TcpSackFailures) / interval
	netStat.TcpExtStat.TcpLossFailuresPerSec = (netStat.TcpExtStat.TcpLossFailures - self.currentNetStat.TcpExtStat.TcpLossFailures) / interval
	netStat.TcpExtStat.TcpFastRetransPerSec = (netStat.TcpExtStat.TcpFastRetrans - self.currentNetStat.TcpExtStat.TcpFastRetrans) / interval
	netStat.TcpExtStat.TcpSlowStartRetransPerSec = (netStat.TcpExtStat.TcpSlowStartRetrans - self.currentNetStat.TcpExtStat.TcpSlowStartRetrans) / interval
	netStat.TcpExtStat.TcpTimeoutsPerSec = (netStat.TcpExtStat.TcpTimeouts - self.currentNetStat.TcpExtStat.TcpTimeouts) / interval
	netStat.TcpExtStat.TcpLossProbesPerSec = (netStat.TcpExtStat.TcpLossProbes - self.currentNetStat.TcpExtStat.TcpLossProbes) / interval
	netStat.TcpExtStat.TcpLossProbeRecoveryPerSec = (netStat.TcpExtStat.TcpLossProbeRecovery - self.currentNetStat.TcpExtStat.TcpLossProbeRecovery) / interval
	netStat.TcpExtStat.TcpRenoRecoveryFailPerSec = (netStat.TcpExtStat.TcpRenoRecoveryFail - self.currentNetStat.TcpExtStat.TcpRenoRecoveryFail) / interval
	netStat.TcpExtStat.TcpSackRecoveryFailPerSec = (netStat.TcpExtStat.TcpSackRecoveryFail - self.currentNetStat.TcpExtStat.TcpSackRecoveryFail) / interval
	netStat.TcpExtStat.TcpRcvCollapsedPerSec = (netStat.TcpExtStat.TcpRcvCollapsed - self.currentNetStat.TcpExtStat.TcpRcvCollapsed) / interval
	netStat.TcpExtStat.TcpBacklogCoalescePerSec = (netStat.TcpExtStat.TcpBacklogCoalesce - self.currentNetStat.TcpExtStat.TcpBacklogCoalesce) / interval
	netStat.TcpExtStat.TcpDsackOldSentPerSec = (netStat.TcpExtStat.TcpDsackOldSent - self.currentNetStat.TcpExtStat.TcpDsackOldSent) / interval
	netStat.TcpExtStat.TcpDsackOfoSentPerSec = (netStat.TcpExtStat.TcpDsackOfoSent - self.currentNetStat.TcpExtStat.TcpDsackOfoSent) / interval
	netStat.TcpExtStat.TcpDsackRecvPerSec = (netStat.TcpExtStat.TcpDsackRecv - self.currentNetStat.TcpExtStat.TcpDsackRecv) / interval
	netStat.TcpExtStat.TcpDsackOfoRecvPerSec = (netStat.TcpExtStat.TcpDsackOfoRecv - self.currentNetStat.TcpExtStat.TcpDsackOfoRecv) / interval
	netStat.TcpExtStat.TcpAbortOnDataPerSec = (netStat.TcpExtStat.TcpAbortOnData - self.currentNetStat.TcpExtStat.TcpAbortOnData) / interval
	netStat.TcpExtStat.TcpAbortOnClosePerSec = (netStat.TcpExtStat.TcpAbortOnClose - self.currentNetStat.TcpExtStat.TcpAbortOnClose) / interval
	netStat.TcpExtStat.TcpAbortOnMemoryPerSec = (netStat.TcpExtStat.TcpAbortOnMemory - self.currentNetStat.TcpExtStat.TcpAbortOnMemory) / interval
	netStat.TcpExtStat.TcpAbortOnTimeoutPerSec = (netStat.TcpExtStat.TcpAbortOnTimeout - self.currentNetStat.TcpExtStat.TcpAbortOnTimeout) / interval
	netStat.TcpExtStat.TcpAbortOnLingerPerSec = (netStat.TcpExtStat.TcpAbortOnLinger - self.currentNetStat.TcpExtStat.TcpAbortOnLinger) / interval
	netStat.TcpExtStat.TcpAbortFailedPerSec = (netStat.TcpExtStat.TcpAbortFailed - self.currentNetStat.TcpExtStat.TcpAbortFailed) / interval
	netStat.TcpExtStat.TcpMemoryPressuresPerSec = (netStat.TcpExtStat.TcpMemoryPressures - self.currentNetStat.TcpExtStat.TcpMemoryPressures) / interval
	netStat.TcpExtStat.TcpMemoryPressuresChronoPerSec = (netStat.TcpExtStat.TcpMemoryPressuresChrono - self.currentNetStat.TcpExtStat.TcpMemoryPressuresChrono) / interval
	netStat.TcpExtStat.TcpSackDiscardPerSec = (netStat.TcpExtStat.TcpSackDiscard - self.currentNetStat.TcpExtStat.TcpSackDiscard) / interval
	netStat.TcpExtStat.TcpDsackIgnoredOldPerSec = (netStat.TcpExtStat.TcpDsackIgnoredOld - self.currentNetStat.TcpExtStat.TcpDsackIgnoredOld) / interval
	netStat.TcpExtStat.TcpDsackIgnoredNoUndoPerSec = (netStat.TcpExtStat.TcpDsackIgnoredNoUndo - self.currentNetStat.TcpExtStat.TcpDsackIgnoredNoUndo) / interval
	netStat.TcpExtStat.TcpSpuriousRTOsPerSec = (netStat.TcpExtStat.TcpSpuriousRTOs - self.currentNetStat.TcpExtStat.TcpSpuriousRTOs) / interval
	netStat.TcpExtStat.TcpMd5NotFoundPerSec = (netStat.TcpExtStat.TcpMd5NotFound - self.currentNetStat.TcpExtStat.TcpMd5NotFound) / interval
	netStat.TcpExtStat.TcpMd5UnexpectedPerSec = (netStat.TcpExtStat.TcpMd5Unexpected - self.currentNetStat.TcpExtStat.TcpMd5Unexpected) / interval
	netStat.TcpExtStat.TcpMd5FailurePerSec = (netStat.TcpExtStat.TcpMd5Failure - self.currentNetStat.TcpExtStat.TcpMd5Failure) / interval
	netStat.TcpExtStat.TcpSackShiftedPerSec = (netStat.TcpExtStat.TcpSackShifted - self.currentNetStat.TcpExtStat.TcpSackShifted) / interval
	netStat.TcpExtStat.TcpSackMergedPerSec = (netStat.TcpExtStat.TcpSackMerged - self.currentNetStat.TcpExtStat.TcpSackMerged) / interval
	netStat.TcpExtStat.TcpSackShiftFallbackPerSec = (netStat.TcpExtStat.TcpSackShiftFallback - self.currentNetStat.TcpExtStat.TcpSackShiftFallback) / interval
	netStat.TcpExtStat.TcpBacklogDropPerSec = (netStat.TcpExtStat.TcpBacklogDrop - self.currentNetStat.TcpExtStat.TcpBacklogDrop) / interval
	netStat.TcpExtStat.PfMemallocDropPerSec = (netStat.TcpExtStat.PfMemallocDrop - self.currentNetStat.TcpExtStat.PfMemallocDrop) / interval
	netStat.TcpExtStat.TcpMinTtlDropPerSec = (netStat.TcpExtStat.TcpMinTtlDrop - self.currentNetStat.TcpExtStat.TcpMinTtlDrop) / interval
	netStat.TcpExtStat.TcpDeferAcceptDropPerSec = (netStat.TcpExtStat.TcpDeferAcceptDrop - self.currentNetStat.TcpExtStat.TcpDeferAcceptDrop) / interval
	netStat.TcpExtStat.IpReversePathFilterPerSec = (netStat.TcpExtStat.IpReversePathFilter - self.currentNetStat.TcpExtStat.IpReversePathFilter) / interval
	netStat.TcpExtStat.TcpTimeWaitOverflowPerSec = (netStat.TcpExtStat.TcpTimeWaitOverflow - self.currentNetStat.TcpExtStat.TcpTimeWaitOverflow) / interval
	netStat.TcpExtStat.TcpReqQFullDoCookiesPerSec = (netStat.TcpExtStat.TcpReqQFullDoCookies - self.currentNetStat.TcpExtStat.TcpReqQFullDoCookies) / interval
	netStat.TcpExtStat.TcpReqQFullDropPerSec = (netStat.TcpExtStat.TcpReqQFullDrop - self.currentNetStat.TcpExtStat.TcpReqQFullDrop) / interval
	netStat.TcpExtStat.TcpRetransFailPerSec = (netStat.TcpExtStat.TcpRetransFail - self.currentNetStat.TcpExtStat.TcpRetransFail) / interval
	netStat.TcpExtStat.TcpRcvCoalescePerSec = (netStat.TcpExtStat.TcpRcvCoalesce - self.currentNetStat.TcpExtStat.TcpRcvCoalesce) / interval
	netStat.TcpExtStat.TcpOfoQueuePerSec = (netStat.TcpExtStat.TcpOfoQueue - self.currentNetStat.TcpExtStat.TcpOfoQueue) / interval
	netStat.TcpExtStat.TcpOfoDropPerSec = (netStat.TcpExtStat.TcpOfoDrop - self.currentNetStat.TcpExtStat.TcpOfoDrop) / interval
	netStat.TcpExtStat.TcpOfoMergePerSec = (netStat.TcpExtStat.TcpOfoMerge - self.currentNetStat.TcpExtStat.TcpOfoMerge) / interval
	netStat.TcpExtStat.TcpChallengeACKPerSec = (netStat.TcpExtStat.TcpChallengeACK - self.currentNetStat.TcpExtStat.TcpChallengeACK) / interval
	netStat.TcpExtStat.TcpSynChallengePerSec = (netStat.TcpExtStat.TcpSynChallenge - self.currentNetStat.TcpExtStat.TcpSynChallenge) / interval
	netStat.TcpExtStat.TcpFastOpenActivePerSec = (netStat.TcpExtStat.TcpFastOpenActive - self.currentNetStat.TcpExtStat.TcpFastOpenActive) / interval
	netStat.TcpExtStat.TcpFastOpenActiveFailPerSec = (netStat.TcpExtStat.TcpFastOpenActiveFail - self.currentNetStat.TcpExtStat.TcpFastOpenActiveFail) / interval
	netStat.TcpExtStat.TcpFastOpenPassivePerSec = (netStat.TcpExtStat.TcpFastOpenPassive - self.currentNetStat.TcpExtStat.TcpFastOpenPassive) / interval
	netStat.TcpExtStat.TcpFastOpenPassiveFailPerSec = (netStat.TcpExtStat.TcpFastOpenPassiveFail - self.currentNetStat.TcpExtStat.TcpFastOpenPassiveFail) / interval
	netStat.TcpExtStat.TcpFastOpenListenOverflowPerSec = (netStat.TcpExtStat.TcpFastOpenListenOverflow - self.currentNetStat.TcpExtStat.TcpFastOpenListenOverflow) / interval
	netStat.TcpExtStat.TcpFastOpenCookieReqdPerSec = (netStat.TcpExtStat.TcpFastOpenCookieReqd - self.currentNetStat.TcpExtStat.TcpFastOpenCookieReqd) / interval
	netStat.TcpExtStat.TcpFastOpenBlackholePerSec = (netStat.TcpExtStat.TcpFastOpenBlackhole - self.currentNetStat.TcpExtStat.TcpFastOpenBlackhole) / interval
	netStat.TcpExtStat.TcpSpuriousRtxHostQueuesPerSec = (netStat.TcpExtStat.TcpSpuriousRtxHostQueues - self.currentNetStat.TcpExtStat.TcpSpuriousRtxHostQueues) / interval
	netStat.TcpExtStat.BusyPollRxPacketsPerSec = (netStat.TcpExtStat.BusyPollRxPackets - self.currentNetStat.TcpExtStat.BusyPollRxPackets) / interval
	netStat.TcpExtStat.TcpAutoCorkingPerSec = (netStat.TcpExtStat.TcpAutoCorking - self.currentNetStat.TcpExtStat.TcpAutoCorking) / interval
	netStat.TcpExtStat.TcpFromZeroWindowAdvPerSec = (netStat.TcpExtStat.TcpFromZeroWindowAdv - self.currentNetStat.TcpExtStat.TcpFromZeroWindowAdv) / interval
	netStat.TcpExtStat.TcpToZeroWindowAdvPerSec = (netStat.TcpExtStat.TcpToZeroWindowAdv - self.currentNetStat.TcpExtStat.TcpToZeroWindowAdv) / interval
	netStat.TcpExtStat.TcpWantZeroWindowAdvPerSec = (netStat.TcpExtStat.TcpWantZeroWindowAdv - self.currentNetStat.TcpExtStat.TcpWantZeroWindowAdv) / interval
	netStat.TcpExtStat.TcpSynRetransPerSec = (netStat.TcpExtStat.TcpSynRetrans - self.currentNetStat.TcpExtStat.TcpSynRetrans) / interval
	netStat.TcpExtStat.TcpOrigDataSentPerSec = (netStat.TcpExtStat.TcpOrigDataSent - self.currentNetStat.TcpExtStat.TcpOrigDataSent) / interval
	netStat.TcpExtStat.TcpHystartTrainDetectPerSec = (netStat.TcpExtStat.TcpHystartTrainDetect - self.currentNetStat.TcpExtStat.TcpHystartTrainDetect) / interval
	netStat.TcpExtStat.TcpHystartTrainCwndPerSec = (netStat.TcpExtStat.TcpHystartTrainCwnd - self.currentNetStat.TcpExtStat.TcpHystartTrainCwnd) / interval
	netStat.TcpExtStat.TcpHystartDelayDetectPerSec = (netStat.TcpExtStat.TcpHystartDelayDetect - self.currentNetStat.TcpExtStat.TcpHystartDelayDetect) / interval
	netStat.TcpExtStat.TcpHystartDelayCwndPerSec = (netStat.TcpExtStat.TcpHystartDelayCwnd - self.currentNetStat.TcpExtStat.TcpHystartDelayCwnd) / interval
	netStat.TcpExtStat.TcpAckSkippedSynRecvPerSec = (netStat.TcpExtStat.TcpAckSkippedSynRecv - self.currentNetStat.TcpExtStat.TcpAckSkippedSynRecv) / interval
	netStat.TcpExtStat.TcpAckSkippedPAWSPerSec = (netStat.TcpExtStat.TcpAckSkippedPAWS - self.currentNetStat.TcpExtStat.TcpAckSkippedPAWS) / interval
	netStat.TcpExtStat.TcpAckSkippedSeqPerSec = (netStat.TcpExtStat.TcpAckSkippedSeq - self.currentNetStat.TcpExtStat.TcpAckSkippedSeq) / interval
	netStat.TcpExtStat.TcpAckSkippedFinWait2PerSec = (netStat.TcpExtStat.TcpAckSkippedFinWait2 - self.currentNetStat.TcpExtStat.TcpAckSkippedFinWait2) / interval
	netStat.TcpExtStat.TcpAckSkippedTimeWaitPerSec = (netStat.TcpExtStat.TcpAckSkippedTimeWait - self.currentNetStat.TcpExtStat.TcpAckSkippedTimeWait) / interval
	netStat.TcpExtStat.TcpAckSkippedChallengePerSec = (netStat.TcpExtStat.TcpAckSkippedChallenge - self.currentNetStat.TcpExtStat.TcpAckSkippedChallenge) / interval
	netStat.TcpExtStat.TcpWinProbePerSec = (netStat.TcpExtStat.TcpWinProbe - self.currentNetStat.TcpExtStat.TcpWinProbe) / interval
	netStat.TcpExtStat.TcpKeepAlivePerSec = (netStat.TcpExtStat.TcpKeepAlive - self.currentNetStat.TcpExtStat.TcpKeepAlive) / interval
	netStat.TcpExtStat.TcpMtupFailPerSec = (netStat.TcpExtStat.TcpMtupFail - self.currentNetStat.TcpExtStat.TcpMtupFail) / interval
	netStat.TcpExtStat.TcpMtupSuccessPerSec = (netStat.TcpExtStat.TcpMtupSuccess - self.currentNetStat.TcpExtStat.TcpMtupSuccess) / interval
	netStat.TcpExtStat.TcpDeliveredPerSec = (netStat.TcpExtStat.TcpDelivered - self.currentNetStat.TcpExtStat.TcpDelivered) / interval
	netStat.TcpExtStat.TcpDeliveredCEPerSec = (netStat.TcpExtStat.TcpDeliveredCE - self.currentNetStat.TcpExtStat.TcpDeliveredCE) / interval
	netStat.TcpExtStat.TcpAckCompressedPerSec = (netStat.TcpExtStat.TcpAckCompressed - self.currentNetStat.TcpExtStat.TcpAckCompressed) / interval
	netStat.TcpExtStat.TcpZeroWindowDropPerSec = (netStat.TcpExtStat.TcpZeroWindowDrop - self.currentNetStat.TcpExtStat.TcpZeroWindowDrop) / interval
	netStat.TcpExtStat.TcpRcvQDropPerSec = (netStat.TcpExtStat.TcpRcvQDrop - self.currentNetStat.TcpExtStat.TcpRcvQDrop) / interval
	netStat.TcpExtStat.TcpWqueueTooBigPerSec = (netStat.TcpExtStat.TcpWqueueTooBig - self.currentNetStat.TcpExtStat.TcpWqueueTooBig) / interval
	netStat.TcpExtStat.TcpFastOpenPassiveAltKeyPerSec = (netStat.TcpExtStat.TcpFastOpenPassiveAltKey - self.currentNetStat.TcpExtStat.TcpFastOpenPassiveAltKey) / interval

	netStat.IpExtStat.InNoRoutesPerSec = (netStat.IpExtStat.InNoRoutes - self.currentNetStat.IpExtStat.InNoRoutes) / interval
	netStat.IpExtStat.InTruncatedPktsPerSec = (netStat.IpExtStat.InTruncatedPkts - self.currentNetStat.IpExtStat.InTruncatedPkts) / interval
	netStat.IpExtStat.InCsumErrorsPerSec = (netStat.IpExtStat.InCsumErrors - self.currentNetStat.IpExtStat.InCsumErrors) / interval
	netStat.IpExtStat.InNoRoutesPerSec = (netStat.IpExtStat.InNoRoutes - self.currentNetStat.IpExtStat.InNoRoutes) / interval
	netStat.IpExtStat.InTruncatedPktsPerSec = (netStat.IpExtStat.InTruncatedPkts - self.currentNetStat.IpExtStat.InTruncatedPkts) / interval
	netStat.IpExtStat.InMcastPktsPerSec = (netStat.IpExtStat.InMcastPkts - self.currentNetStat.IpExtStat.InMcastPkts) / interval
	netStat.IpExtStat.OutMcastPktsPerSec = (netStat.IpExtStat.OutMcastPkts - self.currentNetStat.IpExtStat.OutMcastPkts) / interval
	netStat.IpExtStat.InBcastPktsPerSec = (netStat.IpExtStat.InBcastPkts - self.currentNetStat.IpExtStat.InBcastPkts) / interval
	netStat.IpExtStat.OutBcastPktsPerSec = (netStat.IpExtStat.OutBcastPkts - self.currentNetStat.IpExtStat.OutBcastPkts) / interval
	netStat.IpExtStat.InOctetsPerSec = (netStat.IpExtStat.InOctets - self.currentNetStat.IpExtStat.InOctets) / interval
	netStat.IpExtStat.OutOctetsPerSec = (netStat.IpExtStat.OutOctets - self.currentNetStat.IpExtStat.OutOctets) / interval
	netStat.IpExtStat.InMcastOctetsPerSec = (netStat.IpExtStat.InMcastOctets - self.currentNetStat.IpExtStat.InMcastOctets) / interval
	netStat.IpExtStat.OutMcastOctetsPerSec = (netStat.IpExtStat.OutMcastOctets - self.currentNetStat.IpExtStat.OutMcastOctets) / interval
	netStat.IpExtStat.InBcastOctetsPerSec = (netStat.IpExtStat.InBcastOctets - self.currentNetStat.IpExtStat.InBcastOctets) / interval
	netStat.IpExtStat.OutBcastOctetsPerSec = (netStat.IpExtStat.OutBcastOctets - self.currentNetStat.IpExtStat.OutBcastOctets) / interval
	netStat.IpExtStat.InCsumErrorsPerSec = (netStat.IpExtStat.InCsumErrors - self.currentNetStat.IpExtStat.InCsumErrors) / interval
	netStat.IpExtStat.InNoECTPktsPerSec = (netStat.IpExtStat.InNoECTPkts - self.currentNetStat.IpExtStat.InNoECTPkts) / interval
	netStat.IpExtStat.InECT1PktsPerSec = (netStat.IpExtStat.InECT1Pkts - self.currentNetStat.IpExtStat.InECT1Pkts) / interval
	netStat.IpExtStat.InECT0PktsPerSec = (netStat.IpExtStat.InECT0Pkts - self.currentNetStat.IpExtStat.InECT0Pkts) / interval
	netStat.IpExtStat.InCEPktsPerSec = (netStat.IpExtStat.InCEPkts - self.currentNetStat.IpExtStat.InCEPkts) / interval
	netStat.IpExtStat.ReasmOverlapsPerSec = (netStat.IpExtStat.ReasmOverlaps - self.currentNetStat.IpExtStat.ReasmOverlaps) / interval

	self.currentNetStat = netStat
}

func (self *StatRunner) syncProcessStat() {
	processes, pidIndexMap, err := GetProcesses("/", true)
	if err != nil {
		return
	}
	if self.currentProcesses == nil {
		self.currentProcesses = processes
		self.currentPidIndexMap = pidIndexMap
		return
	}

	interval := self.interval

	for i, process := range processes {
		currentProcessIndex, ok := self.currentPidIndexMap[process.Pid]
		if !ok {
			continue
		}
		currentProcess := self.currentProcesses[currentProcessIndex]
		stat := process.Stat
		bstat := currentProcess.Stat

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

	self.currentProcesses = processes
	self.currentPidIndexMap = pidIndexMap
}

func (self *StatRunner) syncLoginUserStat() {
	var loginUserStat *LoginUserStat
	var err error
	if loginUserStat, err = GetLoginUserStat(); err != nil {
		return
	}

	self.currentLoginUserStat = loginUserStat
	return
}

func (self *StatRunner) syncUptimeStat() {
	var uptimeStat *UptimeStat
	var err error
	if uptimeStat, err = GetUptimeStat("/"); err != nil {
		return
	}

	self.currentUptimeStat = uptimeStat
	return
}

func (self *StatRunner) Run(runAt time.Time) {
	self.syncCpuStat()
	self.syncMemStat()
	self.syncDiskStat()
	self.syncProcessStat()
	self.syncNetStat()
	self.syncLoginUserStat()
	self.syncUptimeStat()

	stats := &Stats{
		CpuStat:       self.currentCpuStat,
		MemStat:       self.currentMemStat,
		DiskStat:      self.currentDiskStat,
		Processes:     self.currentProcesses,
		NetStat:       self.currentNetStat,
		LoginUserStat: self.currentLoginUserStat,
		UptimeStat:    self.currentUptimeStat,
	}

	if self.currentStats != nil {
		self.currentStats = stats
		if self.handleStats != nil {
			self.handleStats(runAt, stats)
		}
	} else {
		self.currentStats = stats
	}
}

func (self *StatRunner) StopTimeout() {
}
