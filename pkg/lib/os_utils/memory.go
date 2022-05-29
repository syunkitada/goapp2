package os_utils

import (
	"bufio"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/syunkitada/goapp2/pkg/lib/str_utils"
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

	Buddyinfo BuddyinfoStat
}

type BuddyinfoStat struct {
	M4K   int
	M8K   int
	M16K  int
	M32K  int
	M64K  int
	M128K int
	M256K int
	M512K int
	M1M   int
	M2M   int
	M4M   int
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

	// Read /proc/buddyinfo
	// Output example is below.
	// $ /proc/buddyinfo
	//                           4K     8k    16k    32k    64k   128k   256k   512k     1M     2M     4M
	// Node 0, zone      DMA      0      0      0      1      2      1      1      0      1      1      3
	// Node 0, zone    DMA32      3      3      3      3      3      2      5      6      5      2    874
	// Node 0, zone   Normal  24727  53842  18419  15120  10448   4451   1761    804    382    105    229
	buddyinfoFile, _ := os.Open("/proc/buddyinfo")
	defer buddyinfoFile.Close()
	tmpReader = bufio.NewReader(buddyinfoFile)
	for {
		tmpBytes, _, tmpErr = tmpReader.ReadLine()
		if tmpErr != nil {
			break
		}
		buddyinfo := str_utils.SplitSpace(string(tmpBytes))
		if len(buddyinfo) < 10 {
			continue
		}
		if buddyinfo[3] == "Normal" {
			nodeId, _ := strconv.Atoi(buddyinfo[1])
			m4K, _ := strconv.Atoi(buddyinfo[4])
			m8K, _ := strconv.Atoi(buddyinfo[5])
			m16K, _ := strconv.Atoi(buddyinfo[6])
			m32K, _ := strconv.Atoi(buddyinfo[7])
			m64K, _ := strconv.Atoi(buddyinfo[8])
			m128K, _ := strconv.Atoi(buddyinfo[9])
			m256K, _ := strconv.Atoi(buddyinfo[10])
			m512K, _ := strconv.Atoi(buddyinfo[11])
			m1M, _ := strconv.Atoi(buddyinfo[12])
			m2M, _ := strconv.Atoi(buddyinfo[13])
			m4M, _ := strconv.Atoi(buddyinfo[14])

			nodes[nodeId].Buddyinfo = BuddyinfoStat{
				M4K:   m4K,
				M8K:   m8K,
				M16K:  m16K,
				M32K:  m32K,
				M64K:  m64K,
				M128K: m128K,
				M256K: m256K,
				M512K: m512K,
				M1M:   m1M,
				M2M:   m2M,
				M4M:   m4M,
			}
		}
	}

	stat = &MemStat{
		Nodes:  nodes,
		Vmstat: vmstat,
	}

	return
}
