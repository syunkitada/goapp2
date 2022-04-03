package process_utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/syunkitada/goapp2/pkg/lib/str_utils"
)

type Thread struct {
	Pid int
}

type Process struct {
	Name     string
	Pid      int
	Tgid     int
	Ppid     int
	Cmd      string
	Cmds     []string
	Children []int
	Threads  []Thread
	State    int
	Stat     ProcessStat
}

type ProcessStat struct {
	Timestamp                time.Time
	VmSizeKb                 int
	VmRssKb                  int
	State                    int
	SchedCpuTime             int
	SchedWaitTime            int
	SchedTimeSlices          int
	HugetlbPages             int
	Threads                  int
	VoluntaryCtxtSwitches    int
	NonvoluntaryCtxtSwitches int

	Utime     int
	Stime     int
	Gtime     int
	Cgtime    int
	StartTime int

	Syscr      int
	Syscw      int
	ReadBytes  int
	WriteBytes int
}

const ProcDir = "proc/"

func GetProcesses(rootDir string, isVerbose bool) (processes []Process, pidIndexMap map[int]int, err error) {
	var procDirFile *os.File
	procDir := rootDir + ProcDir
	if procDirFile, err = os.Open(procDir); err != nil {
		return
	}

	var procFileInfos []os.FileInfo
	procFileInfos, err = procDirFile.Readdir(-1)
	procDirFile.Close()
	if err != nil {
		return
	}

	timestamp := time.Now()
	index := 0
	pidIndexMap = map[int]int{}
	for _, procFileInfo := range procFileInfos {
		if !procFileInfo.IsDir() {
			continue
		}

		var process *Process
		if process, err = getProcess(procDir, procFileInfo.Name(), isVerbose); err != nil {
			return
		}

		if process == nil {
			// fs, bus, irq, sys, tty, acpi, scsi, asound, driver, sysvipc, pressure はスキップされる
			continue
		}

		process.Stat.Timestamp = timestamp
		pidIndexMap[process.Pid] = index
		index += 1
		processes = append(processes, *process)
	}

	for _, process := range processes {
		if index, ok := pidIndexMap[process.Ppid]; ok {
			processes[index].Children = append(processes[index].Children, process.Pid)
		}
	}

	return
}

func getProcess(rootProcDir string, pidStr string, isVerbose bool) (process *Process, err error) {
	var tmpFile *os.File
	var tmpBytes []byte
	var tmpTexts []string
	var tmpReader *bufio.Reader

	// /proc/self, /proc/fs などのPID(int)でないものは除外する
	var pid int
	var tmpErr error
	if pid, tmpErr = strconv.Atoi(pidStr); tmpErr != nil {
		return
	}

	procDir := rootProcDir + pidStr + "/"

	// ----------------------------------------------------------------------------------------------------
	// Parse cmdline
	cmdlinePath := procDir + "cmdline"
	if tmpBytes, err = ioutil.ReadFile(cmdlinePath); err != nil {
		return
	}
	cmds := strings.Split(string(tmpBytes), string(byte(0)))
	if cmds[len(cmds)-1] == "" {
		cmds = cmds[0 : len(cmds)-1]
	}

	// ----------------------------------------------------------------------------------------------------
	// Parse status
	if tmpFile, err = os.Open(procDir + "status"); err != nil {
		return
	}
	defer tmpFile.Close()
	tmpReader = bufio.NewReader(tmpFile)
	statusMap := map[string][]string{}
	for {
		tmpBytes, _, tmpErr := tmpReader.ReadLine()
		if tmpErr != nil {
			break
		}

		fields := str_utils.SplitSpace(string(tmpBytes))
		statusMap[fields[0]] = fields[1:]
	}

	// Name:   kworker/6:2-events
	name := statusMap["Name:"][0]

	// Umask:  0000
	// State:  I (idle)
	var stateInt int
	switch statusMap["State:"][0] {
	case "R":
		stateInt = 3
	case "D":
		stateInt = 2
	case "S":
		stateInt = 1
	case "I":
		stateInt = 0
	case "Z":
		stateInt = -1
	default:
		stateInt = 0
	}

	// Tgid:   23550
	tgid, _ := strconv.Atoi(statusMap["Tgid:"][0])
	// Ngid:   0
	// Pid:    23550
	// PPid:   23547
	ppid, _ := strconv.Atoi(statusMap["PPid:"][0])
	// TracerPid:      0
	// Uid:    0       0       0       0
	// Gid:    0       0       0       0
	// FDSize: 256
	// Groups:
	// NStgid: 23550
	// NSpid:  23550
	// NSpgid: 23550
	// NSsid:  23547
	// VmPeak:  3235840 kB
	// VmSize:  2461756 kB

	// VmSize
	var vmSizeKb int
	if value, ok := statusMap["VmSize:"]; ok {
		vmSizeKb, _ = strconv.Atoi(value[0])
	}
	// VmLck:         0 kB
	// VmPin:         0 kB
	// VmHWM:     31584 kB
	// VmRSS:     28784 kB
	var vmRssKb int
	if value, ok := statusMap["VmRSS:"]; ok {
		vmRssKb, _ = strconv.Atoi(value[0])
	}
	// RssAnon:           16256 kB
	// RssFile:           12528 kB
	// RssShmem:              0 kB
	// VmData:  2399228 kB
	// VmStk:       132 kB
	// VmExe:     11452 kB
	// VmLib:      7424 kB
	// VmPTE:       572 kB
	// VmSwap:        0 kB
	// HugetlbPages:    2097152 kB
	var hugetlbPages int
	if value, ok := statusMap["HugetlbPages:"]; ok {
		hugetlbPages, _ = strconv.Atoi(value[0])
	}
	// CoreDumping:    0
	// THP_enabled:    1
	// Threads:        4
	threads, _ := strconv.Atoi(statusMap["Threads:"][0])
	// SigQ:   0/62468
	// SigPnd: 0000000000000000
	// ShdPnd: 0000000000000000
	// SigBlk: 0000000010002240
	// SigIgn: 0000000000001000
	// SigCgt: 0000000180004243
	// CapInh: 0000000000000000
	// CapPrm: 0000003fffffffff
	// CapEff: 0000003fffffffff
	// CapBnd: 0000003fffffffff
	// CapAmb: 0000000000000000
	// NoNewPrivs:     0
	// Seccomp:        0
	// Speculation_Store_Bypass:       thread vulnerable
	// Cpus_allowed:   ffff
	// Cpus_allowed_list:      0-15
	// Mems_allowed:   00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000001
	// Mems_allowed_list:      0
	// voluntary_ctxt_switches:        14415
	voluntaryCtxtSwitches, _ := strconv.Atoi(statusMap["voluntary_ctxt_switches:"][0])
	// nonvoluntary_ctxt_switches:     219
	nonvoluntaryCtxtSwitches, _ := strconv.Atoi(statusMap["nonvoluntary_ctxt_switches:"][0])

	process = &Process{
		Name:  name,
		Cmds:  cmds,
		Pid:   pid,
		Tgid:  tgid,
		Ppid:  ppid,
		State: stateInt,
	}
	if !isVerbose {
		return
	}

	// ----------------------------------------------------------------------------------------------------
	// Parse /proc/[pid]/schedstat
	// 2554841551 177487694 35200
	// [time spent on the cpu] [time spent waiting on a runqueue] [timeslices run on this cpu]
	if tmpFile, err = os.Open(procDir + "schedstat"); err != nil {
		return
	}
	tmpReader = bufio.NewReader(tmpFile)
	if tmpBytes, _, err = tmpReader.ReadLine(); err != nil {
		return
	}
	tmpTexts = strings.Split(string(tmpBytes), " ")
	if len(tmpTexts) != 3 {
		err = fmt.Errorf("Unexpected Format: path=/proc/[pid]/schedstat, text=%s", string(tmpBytes))
		return
	}
	// time spent on the cpu
	schedCpuTime, _ := strconv.Atoi(tmpTexts[0])
	// time spent waiting on a runqueue
	schedWaitTime, _ := strconv.Atoi(tmpTexts[1])
	// # of timeslices run on this cpu
	schedTimeSlices, _ := strconv.Atoi(tmpTexts[2])

	// ----------------------------------------------------------------------------------------------------
	// $ cat /proc/24120/stat
	// 24120 (qemu-system-x86) S 24119 24120 24119 0 -1 138412416 23189 0 0 0 2227 753 0 0 20 0 6 0 251962 4969209856 7743 18446744073709551615 1 1 0 0 0 0 268444224 4096 16963 0 0 0 17 9 0 0 0 2041 0 0 0 0 0 0 0 0 0
	if tmpFile, err = os.Open(procDir + "stat"); err != nil {
		return
	}
	tmpReader = bufio.NewReader(tmpFile)
	if tmpBytes, _, err = tmpReader.ReadLine(); err != nil {
		return
	}
	tmpTexts = strings.Split(string(tmpBytes), " ")
	utime, _ := strconv.Atoi(tmpTexts[13])
	stime, _ := strconv.Atoi(tmpTexts[14])
	gtime, _ := strconv.Atoi(tmpTexts[42])
	cgtime, _ := strconv.Atoi(tmpTexts[43])
	startTime, _ := strconv.Atoi(tmpTexts[21])

	// ----------------------------------------------------------------------------------------------------
	// $ cat /proc/24120/io
	// rchar: 160323858
	// wchar: 14532026
	// syscr: 48257
	// syscw: 37187
	// read_bytes: 163528704
	// write_bytes: 15466496
	// cancelled_write_bytes: 0
	// root権限がないと見れない
	if tmpFile, err = os.Open(procDir + "io"); err != nil {
		return
	}
	tmpReader = bufio.NewReader(tmpFile)
	_, _, _ = tmpReader.ReadLine()
	_, _, _ = tmpReader.ReadLine()
	tmpBytes, _, _ = tmpReader.ReadLine()
	syscr, _ := strconv.Atoi(str_utils.ParseLastValue(string(tmpBytes)))
	tmpBytes, _, _ = tmpReader.ReadLine()
	syscw, _ := strconv.Atoi(str_utils.ParseLastValue(string(tmpBytes)))
	tmpBytes, _, _ = tmpReader.ReadLine()
	readBytes, _ := strconv.Atoi(str_utils.ParseLastValue(string(tmpBytes)))
	tmpBytes, _, _ = tmpReader.ReadLine()
	writeBytes, _ := strconv.Atoi(str_utils.ParseLastValue(string(tmpBytes)))

	process.Stat = ProcessStat{
		SchedCpuTime:             schedCpuTime,
		SchedWaitTime:            schedWaitTime,
		SchedTimeSlices:          schedTimeSlices,
		VmSizeKb:                 vmSizeKb,
		VmRssKb:                  vmRssKb,
		HugetlbPages:             hugetlbPages,
		Threads:                  threads,
		VoluntaryCtxtSwitches:    voluntaryCtxtSwitches,
		NonvoluntaryCtxtSwitches: nonvoluntaryCtxtSwitches,

		Utime:     utime,
		Stime:     stime,
		Gtime:     gtime,
		Cgtime:    cgtime,
		StartTime: startTime,

		Syscr:      syscr,
		Syscw:      syscw,
		ReadBytes:  readBytes,
		WriteBytes: writeBytes,
	}

	return
}
