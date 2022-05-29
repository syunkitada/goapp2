package os_utils

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/syunkitada/goapp2/pkg/lib/str_utils"
)

type DiskStat struct {
	DiskDeviceStatMap map[string]DiskDeviceStat
	DiskFsStatMap     map[string]DiskFsStat
}

type DiskDeviceStat struct {
	PblockSize        int
	ReadsCompleted    int
	ReadsMerges       int
	ReadSectors       int
	ReadMs            int
	WritesCompleted   int
	WritesMerges      int
	WriteSectors      int
	WriteMs           int
	ProgressIos       int
	IosMs             int
	WeightedIosMs     int
	DiscardsCompleted int
	DiscardsMerges    int
	DiscardSectors    int
	DiscardMs         int

	ReadsPerSec         int
	RmergesPerSec       int
	ReadBytesPerSec     int
	ReadMsPerSec        int
	WritesPerSec        int
	WmergesPerSec       int
	WriteBytesPerSec    int
	WriteMsPerSec       int
	DiscardsPerSec      int
	DmergesPerSec       int
	DiscardBytesPerSec  int
	DiscardMsPerSec     int
	IosMsPerSec         int
	WeightedIosMsPerSec int
}

type DiskFsStat struct {
	Path      string
	Type      string
	MountPath string
	TotalSize int
	FreeSize  int
	UsedSize  int
	Files     int
}

func GetDiskStat() (diskStat *DiskStat, err error) {
	// Read /proc/diskstats

	// 259       0 nvme0n1 94360 70783 6403078 67950 136558 90723 6419592 38105 0 97140 59208 0 0 0 0
	// 259       0 nvme0n1 94360 70783 6403078 67950 136611 90751 6423880 38111 0 97200 59208 0 0 0 0
	// 259       0 nvme0n1 94364 70783 6403230 67951 155638 101247 7087392 41420 0 107356 59208 0 0 0 0

	// Field  1 -- # of reads completed
	// Field  2 -- # of reads merged, field 6 -- # of writes merged
	// Field  3 -- # of sectors read
	// Field  4 -- # of milliseconds spent reading
	// Field  5 -- # of writes completed
	// Field  6 -- # of writes merged
	// Field  7 -- # of sectors written
	// Field  8 -- # of milliseconds spent writing
	// Field  9 -- # of I/Os currently in progress
	// Field 10 -- # of milliseconds spent doing I/Os
	// Field 11 -- weighted # of milliseconds spent doing I/Os
	// Field 12 -- # of discards completed
	// Field 13 -- # of discards merged
	// Field 14 -- # of sectors discarded
	// Field 15 -- # of milliseconds spent discarding
	diskDeviceStatMap := map[string]DiskDeviceStat{}

	f, _ := os.Open("/proc/diskstats")
	defer f.Close()
	tmpReader := bufio.NewReader(f)
	for {
		tmpBytes, _, tmpErr := tmpReader.ReadLine()
		if tmpErr != nil {
			break
		}
		columns := str_utils.SplitSpace(string(tmpBytes))

		pblockSizeFile, tmpErr := os.Open("/sys/block/" + columns[2] + "/queue/physical_block_size")
		if tmpErr != nil {
			continue
		}
		pblockSizeReader := bufio.NewReader(pblockSizeFile)
		pblockSizeBytes, _, tmpErr := pblockSizeReader.ReadLine()
		pblockSizeFile.Close()
		if tmpErr != nil {
			continue
		}
		pblockSize, _ := strconv.Atoi(string(pblockSizeBytes))

		readsCompleted, _ := strconv.Atoi(columns[3])
		readsMerges, _ := strconv.Atoi(columns[4])
		readSectors, _ := strconv.Atoi(columns[5])
		readMs, _ := strconv.Atoi(columns[6])
		writesCompleted, _ := strconv.Atoi(columns[7])
		writesMerges, _ := strconv.Atoi(columns[8])
		writeSectors, _ := strconv.Atoi(columns[9])
		writeMs, _ := strconv.Atoi(columns[10])

		progressIos, _ := strconv.Atoi(columns[11])
		iosMs, _ := strconv.Atoi(columns[12])
		weightedIosMs, _ := strconv.Atoi(columns[13])

		discardsCompleted, _ := strconv.Atoi(columns[14])
		discardsMerges, _ := strconv.Atoi(columns[15])
		discardSectors, _ := strconv.Atoi(columns[16])
		discardMs, _ := strconv.Atoi(columns[17])

		diskDeviceStatMap[columns[2]] = DiskDeviceStat{
			PblockSize:        pblockSize,
			ReadsCompleted:    readsCompleted,
			ReadsMerges:       readsMerges,
			ReadSectors:       readSectors,
			ReadMs:            readMs,
			WritesCompleted:   writesCompleted,
			WritesMerges:      writesMerges,
			WriteSectors:      writeSectors,
			WriteMs:           writeMs,
			ProgressIos:       progressIos,
			IosMs:             iosMs,
			WeightedIosMs:     weightedIosMs,
			DiscardsCompleted: discardsCompleted,
			DiscardsMerges:    discardsMerges,
			DiscardSectors:    discardSectors,
			DiscardMs:         discardMs,
		}
	}

	// read /proc/self/mounts
	// MEMO: /etc/mtab is symbolic link to /proc/self/mounts
	diskFsStatMap := map[string]DiskFsStat{}
	mountsFile, _ := os.Open("/proc/self/mounts")
	defer mountsFile.Close()
	tmpReader = bufio.NewReader(mountsFile)
	var splitedLine []string
	var tmpBytes []byte
	var tmpErr error
	for {
		tmpBytes, _, tmpErr = tmpReader.ReadLine()
		if tmpErr != nil {
			break
		}
		splitedLine = strings.Split(string(tmpBytes), " ")
		var statfs syscall.Statfs_t
		if tmpErr = syscall.Statfs(splitedLine[1], &statfs); tmpErr != nil {
			continue
		}
		totalSize := int(statfs.Blocks) * int(statfs.Bsize)
		freeSize := int(statfs.Bavail) * int(statfs.Bsize)

		diskFsStatMap[splitedLine[0]] = DiskFsStat{
			Path:      splitedLine[0],
			MountPath: splitedLine[1],
			Type:      splitedLine[2],
			TotalSize: totalSize,
			FreeSize:  freeSize,
			UsedSize:  totalSize - freeSize,
			Files:     int(statfs.Files),
		}
	}

	diskStat = &DiskStat{
		DiskDeviceStatMap: diskDeviceStatMap,
		DiskFsStatMap:     diskFsStatMap,
	}
	return
}
