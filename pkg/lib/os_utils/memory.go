package os_utils

import (
	"bufio"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/syunkitada/goapp2/pkg/lib/str_utils"
	// "github.com/syunkitada/goapp/pkg/lib/logger"
	// "github.com/syunkitada/goapp/pkg/lib/str_utils"
	// "github.com/syunkitada/goapp/pkg/resource/config"
	// "github.com/syunkitada/goapp/pkg/resource/consts"
	// "github.com/syunkitada/goapp/pkg/resource/resource_api/spec"
)

type MemStat struct {
	Nodes  []MemNodeStat
	Vmstat Vmstat
}

type MemNodeStat struct {
	ReportStatus int // 0, 1(GetReport), 2(Reported)
	NodeId       int
	NodeName     string
	MemTotal     int
	MemFree      int
	MemUsed      int
	MemAvailable int
	Active       int
	Inactive     int
	ActiveAnon   int
	InactiveAnon int
	ActiveFile   int
	InactiveFile int
	Unevictable  int
	Mlocked      int
	Dirty        int
	Writeback    int
	FilePages    int
	Mapped       int
	AnonPages    int
	Shmem        int
	KernelStack  int
	PageTables   int
	NfsUnstable  int
	Bounce       int
	WritebackTmp int
	KReclaimable int
	Slab         int
	SReclaimable int
	SUnreclaim   int

	HugePages1GTotal int
	HugePates1GFree  int
	HugePages1GUsed  int
}

type Vmstat struct {
	PgscanKswapd int
	PgscanDirect int
	Pgfault      int
	Pswapin      int
	Pswapout     int

	PgscanKswapdPerSec int
	PgscanDirectPerSec int
	PgfaultPerSec      int
	PswapinPerSec      int
	PswapoutPerSec     int
}

const NodeDir = "sys/devices/system/node"

func GetMemStat(rootDir string) (stat *MemStat, err error) {
	// Read /sys/devices/system/node/node.*/hugepages
	// Read /sys/devices/system/node/node.*/meminfo
	var tmpReader *bufio.Reader
	var tmpBytes []byte
	var tmpFile *os.File
	var tmpErr error

	var nodeDirFile *os.File
	nodeDir := rootDir + NodeDir
	if nodeDirFile, err = os.Open(nodeDir); err != nil {
		return
	}
	defer nodeDirFile.Close()

	var nodeFileInfos []os.FileInfo
	nodeFileInfos, err = nodeDirFile.Readdir(-1)
	if err != nil {
		return
	}

	nodes := []MemNodeStat{}
	for _, nodeFileInfo := range nodeFileInfos {
		if !nodeFileInfo.IsDir() {
			continue
		}
		if strings.Index(nodeFileInfo.Name(), "node") == 0 {
			nodeName := nodeFileInfo.Name()
			id := len(nodes)

			tmpBytes, _ = ioutil.ReadFile("/sys/devices/system/node/" + nodeName + "/hugepages/hugepages-1048576kB/nr_hugepages")
			nr1GHugepages, _ := strconv.Atoi(string(tmpBytes))

			tmpBytes, _ = ioutil.ReadFile("/sys/devices/system/node/" + nodeName + "/hugepages/hugepages-1048576kB/free_hugepages")
			free1GHugepages, _ := strconv.Atoi(string(tmpBytes))

			if tmpFile, err = os.Open("/sys/devices/system/node/" + nodeName + "/meminfo"); err != nil {
				return
			}
			tmpReader = bufio.NewReader(tmpFile)

			tmpBytes, _, _ = tmpReader.ReadLine()
			memTotal, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			memFree, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			memUsed, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			active, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			inactive, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			activeAnon, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			inactiveAnon, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			activeFile, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			inactiveFile, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			unevictable, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			mlocked, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			dirty, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			writeback, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			filePages, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			mapped, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			anonPages, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			shmem, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			kernelStack, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			pageTables, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			nfsUnstable, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			bounce, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			writebackTmp, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			kReclaimable, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			slab, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			sReclaimable, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))
			tmpBytes, _, _ = tmpReader.ReadLine()
			sUnreclaim, _ := strconv.Atoi(str_utils.ParseLastSecondValue(string(tmpBytes)))

			memAvailable := memFree + inactive + kReclaimable + sReclaimable

			memStat := MemNodeStat{
				ReportStatus: 0,
				NodeId:       id,
				NodeName:     nodeName,
				MemTotal:     memTotal,
				MemFree:      memFree,
				MemUsed:      memUsed,
				MemAvailable: memAvailable,
				Active:       active,
				Inactive:     inactive,
				ActiveAnon:   activeAnon,
				InactiveAnon: inactiveAnon,
				ActiveFile:   activeFile,
				InactiveFile: inactiveFile,
				Unevictable:  unevictable,
				Mlocked:      mlocked,
				Dirty:        dirty,
				Writeback:    writeback,
				FilePages:    filePages,
				Mapped:       mapped,
				AnonPages:    anonPages,
				Shmem:        shmem,
				KernelStack:  kernelStack,
				PageTables:   pageTables,
				NfsUnstable:  nfsUnstable,
				Bounce:       bounce,
				WritebackTmp: writebackTmp,
				KReclaimable: kReclaimable,
				Slab:         slab,
				SReclaimable: sReclaimable,
				SUnreclaim:   sUnreclaim,

				HugePages1GTotal: nr1GHugepages,
				HugePates1GFree:  free1GHugepages,
				HugePages1GUsed:  nr1GHugepages - free1GHugepages,
			}
			nodes = append(nodes, memStat)
		}
	}

	// Read /proc/vmstat
	var vmstatFile *os.File
	if vmstatFile, err = os.Open("/proc/vmstat"); err != nil {
		return
	}
	defer vmstatFile.Close()
	tmpReader = bufio.NewReader(vmstatFile)
	vmstatMap := map[string]string{}
	for {
		tmpBytes, _, tmpErr = tmpReader.ReadLine()
		if tmpErr != nil {
			break
		}
		columns := strings.Split(string(tmpBytes), " ")
		vmstatMap[columns[0]] = columns[1]
	}

	pgscanKswapd, _ := strconv.Atoi(str_utils.ParseLastValue(vmstatMap["pgscan_kswapd"]))
	pgscanDirect, _ := strconv.Atoi(str_utils.ParseLastValue(vmstatMap["pgscan_direct"]))
	pgfault, _ := strconv.Atoi(str_utils.ParseLastValue(vmstatMap["pgfault"]))

	pswapin, _ := strconv.Atoi(str_utils.ParseLastValue(vmstatMap["pswapin"]))
	pswapout, _ := strconv.Atoi(str_utils.ParseLastValue(vmstatMap["pswapout"]))

	vmstat := Vmstat{
		PgscanKswapd: pgscanKswapd,
		PgscanDirect: pgscanDirect,
		Pgfault:      pgfault,
		Pswapin:      pswapin,
		Pswapout:     pswapout,
	}

	stat = &MemStat{
		Nodes:  nodes,
		Vmstat: vmstat,
	}

	return
}

// func ReportMetrics() (metrics []spec.ResourceMetric) {
// 	metrics = make([]spec.ResourceMetric, 0, len(reader.vmStats)+len(reader.memStats))
// 	for _, stat := range reader.vmStats {
// 		if stat.ReportStatus == ReportStatusReported {
// 			continue
// 		}
// 		metrics = append(metrics, spec.ResourceMetric{
// 			Name: "system_vmstat",
// 			Time: stat.Timestamp,
// 			Metric: map[string]interface{}{
// 				"pgscan_kswapd": stat.DiffPgscanKswapd,
// 				"pgscan_direct": stat.DiffPgscanDirect,
// 				"pgfault":       stat.DiffPgfault,
// 				"pswapin":       stat.DiffPswapin,
// 				"pswapout":      stat.DiffPswapout,
// 			},
// 		})
// 	}
//
// 	for _, stat := range reader.memStats {
// 		if stat.ReportStatus == ReportStatusReported {
// 			continue
// 		}
//
// 		reclaimable := (stat.Inactive + stat.KReclaimable + stat.SReclaimable) * 1000
// 		metrics = append(metrics, spec.ResourceMetric{
// 			Name: "system_mem",
// 			Time: stat.Timestamp,
// 			Tag: map[string]string{
// 				"node_id": strconv.Itoa(stat.NodeId),
// 			},
// 			Metric: map[string]interface{}{
// 				"reclaimable":   reclaimable,
// 				"mem_total":     stat.MemTotal * 1000,
// 				"mem_free":      stat.MemFree * 1000,
// 				"mem_used":      stat.MemUsed * 1000,
// 				"active":        stat.Active * 1000,
// 				"inactive":      stat.Inactive * 1000,
// 				"active_anon":   stat.ActiveAnon * 1000,
// 				"inactive_anon": stat.InactiveAnon * 1000,
// 				"active_file":   stat.ActiveFile * 1000,
// 				"inactive_file": stat.InactiveFile * 1000,
// 				"unevictable":   stat.Unevictable * 1000,
// 				"mlocked":       stat.Mlocked * 1000,
// 				"dirty":         stat.Dirty * 1000,
// 				"writeback":     stat.Writeback * 1000,
// 				"writeback_tmp": stat.WritebackTmp * 1000,
// 				"k_reclaimable": stat.KReclaimable * 1000,
// 				"slab":          stat.Slab * 1000,
// 				"s_reclaimable": stat.SReclaimable * 1000,
// 				"s_unreclaim":   stat.SUnreclaim * 1000,
// 			},
// 		})
// 	}
//
// 	return
// }
//
// func (reader *MemReader) ReportEvents() (events []spec.ResourceEvent) {
// 	if len(reader.vmStats) == 0 {
// 		return
// 	}
//
// 	stats := reader.memStats[len(reader.memStats)-reader.lenNodes:]
// 	msgs := []string{}
// 	eventCheckAvailableLevel := consts.EventLevelSuccess
// 	for _, stat := range stats {
// 		if reader.checkAvailableWarnNodeCounters[stat.NodeId] > reader.checkAvailableOccurences {
// 			eventCheckAvailableLevel = consts.EventLevelWarning
// 		}
// 		msgs = append(msgs,
// 			fmt.Sprintf("node:%d,total=%d,free=%d,available=%d",
// 				stat.NodeId,
// 				stat.MemTotal,
// 				stat.MemFree,
// 				stat.MemAvailable))
// 	}
//
// 	events = append(events, spec.ResourceEvent{
// 		Name:            "CheckMemAvailable",
// 		Time:            stats[0].Timestamp,
// 		Level:           eventCheckAvailableLevel,
// 		Msg:             strings.Join(msgs, ", "),
// 		ReissueDuration: reader.checkAvailableReissueDuration,
// 	})
//
// 	stat := reader.vmStats[len(reader.vmStats)-1]
// 	eventCheckPgscanLevel := consts.EventLevelSuccess
// 	if reader.checkPgscanWarnPgscanDirectCounter > reader.checkPgscanOccurences {
// 		eventCheckPgscanLevel = consts.EventLevelWarning
// 	}
// 	events = append(events, spec.ResourceEvent{
// 		Name:  "CheckMemPgscan",
// 		Time:  stat.Timestamp,
// 		Level: eventCheckPgscanLevel,
// 		Msg: fmt.Sprintf("Pgscan kswapd=%d, direct=%d",
// 			stat.DiffPgscanKswapd,
// 			stat.DiffPgscanDirect,
// 		),
// 		ReissueDuration: reader.checkPgscanReissueDuration,
// 	})
//
// 	return
// }
//
// func (reader *MemReader) Reported() {
// 	for i := range reader.vmStats {
// 		reader.vmStats[i].ReportStatus = ReportStatusReported
// 	}
//
// 	for i := range reader.memStats {
// 		reader.memStats[i].ReportStatus = ReportStatusReported
// 	}
// 	return
// }

// func (reader *MemReader) readTmpVmStat() (tmpVmStat *TmpVmStat) {
// 	// Read /proc/vmstat
// 	f, _ := os.Open("/proc/vmstat")
// 	defer f.Close()
// 	tmpReader := bufio.NewReader(f)
// 	vmstat := map[string]string{}
// 	for {
// 		tmpBytes, _, tmpErr := tmpReader.ReadLine()
// 		if tmpErr != nil {
// 			break
// 		}
// 		columns := strings.Split(string(tmpBytes), " ")
// 		vmstat[columns[0]] = columns[1]
// 	}
//
// 	pgscanKswapd, _ := strconv.Atoi(str_utils.ParseLastValue(vmstat["pgscan_kswapd"]))
// 	pgscanDirect, _ := strconv.Atoi(str_utils.ParseLastValue(vmstat["pgscan_direct"]))
// 	pgfault, _ := strconv.Atoi(str_utils.ParseLastValue(vmstat["pgfault"]))
//
// 	pswapin, _ := strconv.Atoi(str_utils.ParseLastValue(vmstat["pswapin"]))
// 	pswapout, _ := strconv.Atoi(str_utils.ParseLastValue(vmstat["pswapout"]))
//
// 	tmpVmStat = &TmpVmStat{
// 		PgscanKswapd: pgscanKswapd,
// 		PgscanDirect: pgscanDirect,
// 		Pgfault:      pgfault,
// 		Pswapin:      pswapin,
// 		Pswapout:     pswapout,
// 	}
// 	return
// }
