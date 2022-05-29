package os_utils

import (
	"bufio"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/syunkitada/goapp2/pkg/lib/str_utils"
)

type NetStat struct {
	TcpExtStat    TcpExtStat
	IpExtStat     IpExtStat
	NetDevStatMap map[string]NetDevStat
}

type TcpExtStat struct {
	SyncookiesSent            int
	SyncookiesRecv            int
	SyncookiesFailed          int
	EmbryonicRsts             int
	PruneCalled               int
	RcvPruned                 int
	OfoPruned                 int
	OutOfWindowIcmps          int
	LockDroppedIcmps          int
	ArpFilter                 int
	Tw                        int
	TwRecycled                int
	TwKilled                  int
	PawsActive                int
	PawsEstab                 int
	DelayedAcks               int
	DelayedAckLocked          int
	DelayedAckLost            int
	ListenOverflows           int
	ListenDrops               int
	TcpHpHits                 int
	TcpPureAcks               int
	TcpHpAcks                 int
	TcpRenoRecovery           int
	TcpSackRecovery           int
	TcpSackReneging           int
	TcpSackReorder            int
	TcpRenoReorder            int
	TcpTsReorder              int
	TcpFullUndo               int
	TcpPartialUndo            int
	TcpDsackUndo              int
	TcpLossUndo               int
	TcpLostRetransmit         int
	TcpRenoFailures           int
	TcpSackFailures           int
	TcpLossFailures           int
	TcpFastRetrans            int
	TcpSlowStartRetrans       int
	TcpTimeouts               int
	TcpLossProbes             int
	TcpLossProbeRecovery      int
	TcpRenoRecoveryFail       int
	TcpSackRecoveryFail       int
	TcpRcvCollapsed           int
	TcpBacklogCoalesce        int
	TcpDsackOldSent           int
	TcpDsackOfoSent           int
	TcpDsackRecv              int
	TcpDsackOfoRecv           int
	TcpAbortOnData            int
	TcpAbortOnClose           int
	TcpAbortOnMemory          int
	TcpAbortOnTimeout         int
	TcpAbortOnLinger          int
	TcpAbortFailed            int
	TcpMemoryPressures        int
	TcpMemoryPressuresChrono  int
	TcpSackDiscard            int
	TcpDsackIgnoredOld        int
	TcpDsackIgnoredNoUndo     int
	TcpSpuriousRTOs           int
	TcpMd5NotFound            int
	TcpMd5Unexpected          int
	TcpMd5Failure             int
	TcpSackShifted            int
	TcpSackMerged             int
	TcpSackShiftFallback      int
	TcpBacklogDrop            int
	PfMemallocDrop            int
	TcpMinTtlDrop             int
	TcpDeferAcceptDrop        int
	IpReversePathFilter       int
	TcpTimeWaitOverflow       int
	TcpReqQFullDoCookies      int
	TcpReqQFullDrop           int
	TcpRetransFail            int
	TcpRcvCoalesce            int
	TcpOfoQueue               int
	TcpOfoDrop                int
	TcpOfoMerge               int
	TcpChallengeACK           int
	TcpSynChallenge           int
	TcpFastOpenActive         int
	TcpFastOpenActiveFail     int
	TcpFastOpenPassive        int
	TcpFastOpenPassiveFail    int
	TcpFastOpenListenOverflow int
	TcpFastOpenCookieReqd     int
	TcpFastOpenBlackhole      int
	TcpSpuriousRtxHostQueues  int
	BusyPollRxPackets         int
	TcpAutoCorking            int
	TcpFromZeroWindowAdv      int
	TcpToZeroWindowAdv        int
	TcpWantZeroWindowAdv      int
	TcpSynRetrans             int
	TcpOrigDataSent           int
	TcpHystartTrainDetect     int
	TcpHystartTrainCwnd       int
	TcpHystartDelayDetect     int
	TcpHystartDelayCwnd       int
	TcpAckSkippedSynRecv      int
	TcpAckSkippedPAWS         int
	TcpAckSkippedSeq          int
	TcpAckSkippedFinWait2     int
	TcpAckSkippedTimeWait     int
	TcpAckSkippedChallenge    int
	TcpWinProbe               int
	TcpKeepAlive              int
	TcpMtupFail               int
	TcpMtupSuccess            int
	TcpDelivered              int
	TcpDeliveredCE            int
	TcpAckCompressed          int
	TcpZeroWindowDrop         int
	TcpRcvQDrop               int
	TcpWqueueTooBig           int
	TcpFastOpenPassiveAltKey  int

	// PerSec
	SyncookiesSentPerSec            int
	SyncookiesRecvPerSec            int
	SyncookiesFailedPerSec          int
	EmbryonicRstsPerSec             int
	PruneCalledPerSec               int
	RcvPrunedPerSec                 int
	OfoPrunedPerSec                 int
	OutOfWindowIcmpsPerSec          int
	LockDroppedIcmpsPerSec          int
	ArpFilterPerSec                 int
	TwPerSec                        int
	TwRecycledPerSec                int
	TwKilledPerSec                  int
	PawsActivePerSec                int
	PawsEstabPerSec                 int
	DelayedAcksPerSec               int
	DelayedAckLockedPerSec          int
	DelayedAckLostPerSec            int
	ListenOverflowsPerSec           int
	ListenDropsPerSec               int
	TcpHpHitsPerSec                 int
	TcpPureAcksPerSec               int
	TcpHpAcksPerSec                 int
	TcpRenoRecoveryPerSec           int
	TcpSackRecoveryPerSec           int
	TcpSackRenegingPerSec           int
	TcpSackReorderPerSec            int
	TcpRenoReorderPerSec            int
	TcpTsReorderPerSec              int
	TcpFullUndoPerSec               int
	TcpPartialUndoPerSec            int
	TcpDsackUndoPerSec              int
	TcpLossUndoPerSec               int
	TcpLostRetransmitPerSec         int
	TcpRenoFailuresPerSec           int
	TcpSackFailuresPerSec           int
	TcpLossFailuresPerSec           int
	TcpFastRetransPerSec            int
	TcpSlowStartRetransPerSec       int
	TcpTimeoutsPerSec               int
	TcpLossProbesPerSec             int
	TcpLossProbeRecoveryPerSec      int
	TcpRenoRecoveryFailPerSec       int
	TcpSackRecoveryFailPerSec       int
	TcpRcvCollapsedPerSec           int
	TcpBacklogCoalescePerSec        int
	TcpDsackOldSentPerSec           int
	TcpDsackOfoSentPerSec           int
	TcpDsackRecvPerSec              int
	TcpDsackOfoRecvPerSec           int
	TcpAbortOnDataPerSec            int
	TcpAbortOnClosePerSec           int
	TcpAbortOnMemoryPerSec          int
	TcpAbortOnTimeoutPerSec         int
	TcpAbortOnLingerPerSec          int
	TcpAbortFailedPerSec            int
	TcpMemoryPressuresPerSec        int
	TcpMemoryPressuresChronoPerSec  int
	TcpSackDiscardPerSec            int
	TcpDsackIgnoredOldPerSec        int
	TcpDsackIgnoredNoUndoPerSec     int
	TcpSpuriousRTOsPerSec           int
	TcpMd5NotFoundPerSec            int
	TcpMd5UnexpectedPerSec          int
	TcpMd5FailurePerSec             int
	TcpSackShiftedPerSec            int
	TcpSackMergedPerSec             int
	TcpSackShiftFallbackPerSec      int
	TcpBacklogDropPerSec            int
	PfMemallocDropPerSec            int
	TcpMinTtlDropPerSec             int
	TcpDeferAcceptDropPerSec        int
	IpReversePathFilterPerSec       int
	TcpTimeWaitOverflowPerSec       int
	TcpReqQFullDoCookiesPerSec      int
	TcpReqQFullDropPerSec           int
	TcpRetransFailPerSec            int
	TcpRcvCoalescePerSec            int
	TcpOfoQueuePerSec               int
	TcpOfoDropPerSec                int
	TcpOfoMergePerSec               int
	TcpChallengeACKPerSec           int
	TcpSynChallengePerSec           int
	TcpFastOpenActivePerSec         int
	TcpFastOpenActiveFailPerSec     int
	TcpFastOpenPassivePerSec        int
	TcpFastOpenPassiveFailPerSec    int
	TcpFastOpenListenOverflowPerSec int
	TcpFastOpenCookieReqdPerSec     int
	TcpFastOpenBlackholePerSec      int
	TcpSpuriousRtxHostQueuesPerSec  int
	BusyPollRxPacketsPerSec         int
	TcpAutoCorkingPerSec            int
	TcpFromZeroWindowAdvPerSec      int
	TcpToZeroWindowAdvPerSec        int
	TcpWantZeroWindowAdvPerSec      int
	TcpSynRetransPerSec             int
	TcpOrigDataSentPerSec           int
	TcpHystartTrainDetectPerSec     int
	TcpHystartTrainCwndPerSec       int
	TcpHystartDelayDetectPerSec     int
	TcpHystartDelayCwndPerSec       int
	TcpAckSkippedSynRecvPerSec      int
	TcpAckSkippedPAWSPerSec         int
	TcpAckSkippedSeqPerSec          int
	TcpAckSkippedFinWait2PerSec     int
	TcpAckSkippedTimeWaitPerSec     int
	TcpAckSkippedChallengePerSec    int
	TcpWinProbePerSec               int
	TcpKeepAlivePerSec              int
	TcpMtupFailPerSec               int
	TcpMtupSuccessPerSec            int
	TcpDeliveredPerSec              int
	TcpDeliveredCEPerSec            int
	TcpAckCompressedPerSec          int
	TcpZeroWindowDropPerSec         int
	TcpRcvQDropPerSec               int
	TcpWqueueTooBigPerSec           int
	TcpFastOpenPassiveAltKeyPerSec  int
}

type IpExtStat struct {
	InNoRoutes      int
	InTruncatedPkts int
	InMcastPkts     int
	OutMcastPkts    int
	InBcastPkts     int
	OutBcastPkts    int
	InOctets        int
	OutOctets       int
	InMcastOctets   int
	OutMcastOctets  int
	InBcastOctets   int
	OutBcastOctets  int
	InCsumErrors    int
	InNoECTPkts     int
	InECT1Pkts      int
	InECT0Pkts      int
	InCEPkts        int
	ReasmOverlaps   int

	// stat
	InNoRoutesPerSec      int
	InTruncatedPktsPerSec int
	InMcastPktsPerSec     int
	OutMcastPktsPerSec    int
	InBcastPktsPerSec     int
	OutBcastPktsPerSec    int
	InOctetsPerSec        int
	OutOctetsPerSec       int
	InMcastOctetsPerSec   int
	OutMcastOctetsPerSec  int
	InBcastOctetsPerSec   int
	OutBcastOctetsPerSec  int
	InCsumErrorsPerSec    int
	InNoECTPktsPerSec     int
	InECT1PktsPerSec      int
	InECT0PktsPerSec      int
	InCEPktsPerSec        int
	ReasmOverlapsPerSec   int
}

type NetDevStat struct {
	ReceiveBytes    int
	ReceivePackets  int
	ReceiveErrors   int
	ReceiveDrops    int
	TransmitBytes   int
	TransmitPackets int
	TransmitErrors  int
	TransmitDrops   int

	ReceiveBytesPerSec    int
	ReceivePacketsPerSec  int
	ReceiveErrorsPerSec   int
	ReceiveDropsPerSec    int
	TransmitBytesPerSec   int
	TransmitPacketsPerSec int
	TransmitErrorsPerSec  int
	TransmitDropsPerSec   int
}

func GetNetStat() (netStat *NetStat, err error) {
	// $ cat /proc/net/snmp
	netstatFile, _ := os.Open("/proc/net/netstat")
	defer netstatFile.Close()
	tmpReader := bufio.NewReader(netstatFile)

	tmpBytes, _, _ := tmpReader.ReadLine()
	tcpExtKeys := strings.Split(string(tmpBytes), " ")
	lenKeys := len(tcpExtKeys)

	tmpBytes, _, _ = tmpReader.ReadLine()
	tcpExtValues := strings.Split(string(tmpBytes), " ")

	tcpExtMap := map[string]int{}
	for i := 1; i < lenKeys; i++ {
		tcpExtMap[tcpExtKeys[i]], _ = strconv.Atoi(tcpExtValues[i])
	}

	tcpExtStat := TcpExtStat{
		SyncookiesSent:            tcpExtMap["SyncookiesSent"],
		SyncookiesRecv:            tcpExtMap["SyncookiesRecv"],
		SyncookiesFailed:          tcpExtMap["SyncookiesFailed"],
		EmbryonicRsts:             tcpExtMap["EmbryonicRsts"],
		PruneCalled:               tcpExtMap["PruneCalled"],
		RcvPruned:                 tcpExtMap["RcvPruned"],
		OfoPruned:                 tcpExtMap["OfoPruned"],
		OutOfWindowIcmps:          tcpExtMap["OutOfWindowIcmps"],
		LockDroppedIcmps:          tcpExtMap["LockDroppedIcmps"],
		ArpFilter:                 tcpExtMap["ArpFilter"],
		Tw:                        tcpExtMap["TW"],
		TwRecycled:                tcpExtMap["TWRecycled"],
		TwKilled:                  tcpExtMap["TWKilled"],
		PawsActive:                tcpExtMap["PAWSActive"],
		PawsEstab:                 tcpExtMap["PAWSEstab"],
		DelayedAcks:               tcpExtMap["DelayedACKs"],
		DelayedAckLocked:          tcpExtMap["DelayedACKLocked"],
		DelayedAckLost:            tcpExtMap["DelayedACKLost"],
		ListenOverflows:           tcpExtMap["ListenOverflows"],
		ListenDrops:               tcpExtMap["ListenDrops"],
		TcpHpHits:                 tcpExtMap["TCPHPHits"],
		TcpPureAcks:               tcpExtMap["TCPPureAcks"],
		TcpHpAcks:                 tcpExtMap["TCPHPAcks"],
		TcpRenoRecovery:           tcpExtMap["TCPRenoRecovery"],
		TcpSackRecovery:           tcpExtMap["TCPSackRecovery"],
		TcpSackReneging:           tcpExtMap["TCPSACKReneging"],
		TcpSackReorder:            tcpExtMap["TCPSACKReorder"],
		TcpRenoReorder:            tcpExtMap["TCPRenoReorder"],
		TcpTsReorder:              tcpExtMap["TCPTSReorder"],
		TcpFullUndo:               tcpExtMap["TCPFullUndo"],
		TcpPartialUndo:            tcpExtMap["TCPPartialUndo"],
		TcpDsackUndo:              tcpExtMap["TCPDSACKUndo"],
		TcpLossUndo:               tcpExtMap["TCPLossUndo"],
		TcpLostRetransmit:         tcpExtMap["TCPLostRetransmit"],
		TcpRenoFailures:           tcpExtMap["TCPRenoFailures"],
		TcpSackFailures:           tcpExtMap["TCPSackFailures"],
		TcpLossFailures:           tcpExtMap["TCPLossFailures"],
		TcpFastRetrans:            tcpExtMap["TCPFastRetrans"],
		TcpSlowStartRetrans:       tcpExtMap["TCPSlowStartRetrans"],
		TcpTimeouts:               tcpExtMap["TCPTimeouts"],
		TcpLossProbes:             tcpExtMap["TCPLossProbes"],
		TcpLossProbeRecovery:      tcpExtMap["TCPLossProbeRecovery"],
		TcpRenoRecoveryFail:       tcpExtMap["TCPRenoRecoveryFail"],
		TcpSackRecoveryFail:       tcpExtMap["TCPSackRecoveryFail"],
		TcpRcvCollapsed:           tcpExtMap["TCPRcvCollapsed"],
		TcpBacklogCoalesce:        tcpExtMap["TCPBacklogCoalesce"],
		TcpDsackOldSent:           tcpExtMap["TCPDSACKOldSent"],
		TcpDsackOfoSent:           tcpExtMap["TCPDSACKOfoSent"],
		TcpDsackRecv:              tcpExtMap["TCPDSACKRecv"],
		TcpDsackOfoRecv:           tcpExtMap["TCPDSACKOfoRecv"],
		TcpAbortOnData:            tcpExtMap["TCPAbortOnData"],
		TcpAbortOnClose:           tcpExtMap["TCPAbortOnClose"],
		TcpAbortOnMemory:          tcpExtMap["TCPAbortOnMemory"],
		TcpAbortOnTimeout:         tcpExtMap["TCPAbortOnTimeout"],
		TcpAbortOnLinger:          tcpExtMap["TCPAbortOnLinger"],
		TcpAbortFailed:            tcpExtMap["TCPAbortFailed"],
		TcpMemoryPressures:        tcpExtMap["TCPMemoryPressures"],
		TcpMemoryPressuresChrono:  tcpExtMap["TCPMemoryPressuresChrono"],
		TcpSackDiscard:            tcpExtMap["TCPSACKDiscard"],
		TcpDsackIgnoredOld:        tcpExtMap["TCPDSACKIgnoredOld"],
		TcpDsackIgnoredNoUndo:     tcpExtMap["TCPDSACKIgnoredNoUndo"],
		TcpSpuriousRTOs:           tcpExtMap["TCPSpuriousRTOs"],
		TcpMd5NotFound:            tcpExtMap["TCPMD5NotFound"],
		TcpMd5Unexpected:          tcpExtMap["TCPMD5Unexpected"],
		TcpMd5Failure:             tcpExtMap["TCPMD5Failure"],
		TcpSackShifted:            tcpExtMap["TCPSackShifted"],
		TcpSackMerged:             tcpExtMap["TCPSackMerged"],
		TcpSackShiftFallback:      tcpExtMap["TCPSackShiftFallback"],
		TcpBacklogDrop:            tcpExtMap["TCPBacklogDrop"],
		PfMemallocDrop:            tcpExtMap["PFMemallocDrop"],
		TcpMinTtlDrop:             tcpExtMap["TCPMinTTLDrop"],
		TcpDeferAcceptDrop:        tcpExtMap["TCPDeferAcceptDrop"],
		IpReversePathFilter:       tcpExtMap["IPReversePathFilter"],
		TcpTimeWaitOverflow:       tcpExtMap["TCPTimeWaitOverflow"],
		TcpReqQFullDoCookies:      tcpExtMap["TCPReqQFullDoCookies"],
		TcpReqQFullDrop:           tcpExtMap["TCPReqQFullDrop"],
		TcpRetransFail:            tcpExtMap["TCPRetransFail"],
		TcpRcvCoalesce:            tcpExtMap["TCPRcvCoalesce"],
		TcpOfoQueue:               tcpExtMap["TCPOFOQueue"],
		TcpOfoDrop:                tcpExtMap["TCPOFODrop"],
		TcpOfoMerge:               tcpExtMap["TCPOFOMerge"],
		TcpChallengeACK:           tcpExtMap["TCPChallengeACK"],
		TcpSynChallenge:           tcpExtMap["TCPSYNChallenge"],
		TcpFastOpenActive:         tcpExtMap["TCPFastOpenActive"],
		TcpFastOpenActiveFail:     tcpExtMap["TCPFastOpenActiveFail"],
		TcpFastOpenPassive:        tcpExtMap["TCPFastOpenPassive"],
		TcpFastOpenPassiveFail:    tcpExtMap["TCPFastOpenPassiveFail"],
		TcpFastOpenListenOverflow: tcpExtMap["TCPFastOpenListenOverflow"],
		TcpFastOpenCookieReqd:     tcpExtMap["TCPFastOpenCookieReqd"],
		TcpFastOpenBlackhole:      tcpExtMap["TCPFastOpenBlackhole"],
		TcpSpuriousRtxHostQueues:  tcpExtMap["TCPSpuriousRtxHostQueues"],
		BusyPollRxPackets:         tcpExtMap["BusyPollRxPackets"],
		TcpAutoCorking:            tcpExtMap["TCPAutoCorking"],
		TcpFromZeroWindowAdv:      tcpExtMap["TCPFromZeroWindowAdv"],
		TcpToZeroWindowAdv:        tcpExtMap["TCPToZeroWindowAdv"],
		TcpWantZeroWindowAdv:      tcpExtMap["TCPWantZeroWindowAdv"],
		TcpSynRetrans:             tcpExtMap["TCPSynRetrans"],
		TcpOrigDataSent:           tcpExtMap["TCPOrigDataSent"],
		TcpHystartTrainDetect:     tcpExtMap["TCPHystartTrainDetect"],
		TcpHystartTrainCwnd:       tcpExtMap["TCPHystartTrainCwnd"],
		TcpHystartDelayDetect:     tcpExtMap["TCPHystartDelayDetect"],
		TcpHystartDelayCwnd:       tcpExtMap["TCPHystartDelayCwnd"],
		TcpAckSkippedSynRecv:      tcpExtMap["TCPACKSkippedSynRecv"],
		TcpAckSkippedPAWS:         tcpExtMap["TCPACKSkippedPAWS"],
		TcpAckSkippedSeq:          tcpExtMap["TCPACKSkippedSeq"],
		TcpAckSkippedFinWait2:     tcpExtMap["TCPACKSkippedFinWait2"],
		TcpAckSkippedTimeWait:     tcpExtMap["TCPACKSkippedTimeWait"],
		TcpAckSkippedChallenge:    tcpExtMap["TCPACKSkippedChallenge"],
		TcpWinProbe:               tcpExtMap["TCPWinProbe"],
		TcpKeepAlive:              tcpExtMap["TCPKeepAlive"],
		TcpMtupFail:               tcpExtMap["TCPMTUPFail"],
		TcpMtupSuccess:            tcpExtMap["TCPMTUPSuccess"],
		TcpDelivered:              tcpExtMap["TCPDelivered"],
		TcpDeliveredCE:            tcpExtMap["TCPDeliveredCE"],
		TcpAckCompressed:          tcpExtMap["TCPAckCompressed"],
		TcpZeroWindowDrop:         tcpExtMap["TCPZeroWindowDrop"],
		TcpRcvQDrop:               tcpExtMap["TCPRcvQDrop"],
		TcpWqueueTooBig:           tcpExtMap["TCPWqueueTooBig"],
		TcpFastOpenPassiveAltKey:  tcpExtMap["TCPFastOpenPassiveAltKey"],
	}

	// ipExt
	tmpBytes, _, _ = tmpReader.ReadLine()
	ipExtKeys := strings.Split(string(tmpBytes), " ")
	lenKeys = len(ipExtKeys)

	tmpBytes, _, _ = tmpReader.ReadLine()
	ipExtValues := strings.Split(string(tmpBytes), " ")

	ipExtMap := map[string]int{}
	for i := 1; i < lenKeys; i++ {
		ipExtMap[ipExtKeys[i]], _ = strconv.Atoi(ipExtValues[i])
	}

	ipExtStat := IpExtStat{
		InNoRoutes:      ipExtMap["InNoRoutes"],
		InTruncatedPkts: ipExtMap["InTruncatedPkts"],
		InMcastPkts:     ipExtMap["InMcastPkts"],
		OutMcastPkts:    ipExtMap["OutMcastPkts"],
		InBcastPkts:     ipExtMap["InBcastPkts"],
		OutBcastPkts:    ipExtMap["OutBcastPkts"],
		InOctets:        ipExtMap["InOctets"],
		OutOctets:       ipExtMap["OutOctets"],
		InMcastOctets:   ipExtMap["InMcastOctets"],
		OutMcastOctets:  ipExtMap["OutMcastOctets"],
		InBcastOctets:   ipExtMap["InBcastOctets"],
		OutBcastOctets:  ipExtMap["OutBcastOctets"],
		InCsumErrors:    ipExtMap["InCsumErrors"],
		InNoECTPkts:     ipExtMap["InNoECTPkts"],
		InECT1Pkts:      ipExtMap["InECT1Pkts"],
		InECT0Pkts:      ipExtMap["InECT0Pkts"],
		InCEPkts:        ipExtMap["InCEPkts"],
		ReasmOverlaps:   ipExtMap["ReasmOverlaps"],
	}

	// $ cat /proc/net/dev
	// Inter-|   Receive                                                |  Transmit
	//  face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
	//  com-1-ex:    1426      19    0    0    0     0          0         0     4616      43    0    0    0     0       0          0
	//  enp31s0: 7855580   30554    0    0    0     0          0      1408 19677375   42829    0    0    0     0       0          0
	//      lo: 1442597782 3051437    0    0    0     0          0         0 1442597782 3051437    0    0    0     0       0          0
	//   com-0-ex:   29026     447    0    0    0     0          0         0    34621     471    0    0    0     0       0          0
	//   com-2-ex:   26578     383    0    0    0     0          0         0    32083     406    0    0    0     0       0          0
	//   com-4-ex:   28084     420    0    0    0     0          0         0    33499     442    0    0    0     0       0          0
	//   docker0:       0       0    0    0    0     0          0         0        0       0    0    0    0     0       0          0
	bytes, tmpErr := ioutil.ReadFile("/proc/net/dev")
	if tmpErr != nil {
		return
	}
	netDevStatMap := parseNetDev(string(bytes))

	netStat = &NetStat{
		TcpExtStat:    tcpExtStat,
		IpExtStat:     ipExtStat,
		NetDevStatMap: netDevStatMap,
	}
	return
}

func parseNetDev(out string) (netDevStatMap map[string]NetDevStat) {
	// $ cat /proc/net/dev
	// Inter-|   Receive                                                |  Transmit
	//  face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
	//  com-1-ex:    1426      19    0    0    0     0          0         0     4616      43    0    0    0     0       0          0
	//  enp31s0: 7855580   30554    0    0    0     0          0      1408 19677375   42829    0    0    0     0       0          0
	//      lo: 1442597782 3051437    0    0    0     0          0         0 1442597782 3051437    0    0    0     0       0          0
	//   com-0-ex:   29026     447    0    0    0     0          0         0    34621     471    0    0    0     0       0          0
	//   com-2-ex:   26578     383    0    0    0     0          0         0    32083     406    0    0    0     0       0          0
	//   com-4-ex:   28084     420    0    0    0     0          0         0    33499     442    0    0    0     0       0          0
	//   docker0:       0       0    0    0    0     0          0         0        0       0    0    0    0     0       0          0

	splited := strings.Split(out, "\n")

	netDevStatMap = map[string]NetDevStat{}
	lenSplited := len(splited)
	for i := 2; i < lenSplited; i++ {
		splitedStr := str_utils.SplitColon(splited[i])
		if len(splitedStr) < 2 {
			continue
		}
		columns := str_utils.SplitSpace(splitedStr[1])

		receiveBytes, _ := strconv.Atoi(columns[0])
		receivePackets, _ := strconv.Atoi(columns[1])
		receiveErrors, _ := strconv.Atoi(columns[2])
		receiveDrops, _ := strconv.Atoi(columns[3])

		transmitBytes, _ := strconv.Atoi(columns[8])
		transmitPackets, _ := strconv.Atoi(columns[9])
		transmitErrors, _ := strconv.Atoi(columns[10])
		transmitDrops, _ := strconv.Atoi(columns[11])

		netDevStatMap[splitedStr[0]] = NetDevStat{
			ReceiveBytes:    receiveBytes,
			ReceivePackets:  receivePackets,
			ReceiveErrors:   receiveErrors,
			ReceiveDrops:    receiveDrops,
			TransmitBytes:   transmitBytes,
			TransmitPackets: transmitPackets,
			TransmitErrors:  transmitErrors,
			TransmitDrops:   transmitDrops,
		}
	}

	return
}
