package os_utils

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/syunkitada/goapp2/pkg/lib/str_utils"
	// "github.com/syunkitada/goapp/pkg/lib/logger"
	// "github.com/syunkitada/goapp/pkg/lib/str_utils"
	// "github.com/syunkitada/goapp/pkg/resource/config"
	// "github.com/syunkitada/goapp/pkg/resource/consts"
	// "github.com/syunkitada/goapp/pkg/resource/resource_api/spec"
)

type CpuProcessorStat struct {
	Timestamp  time.Time
	Processor  int
	PhysicalId int
	CoreId     int
	Mhz        float64
	User       float64
	Nice       float64
	System     float64
	Idle       float64
	Iowait     float64
	Irq        float64
	Softirq    float64
	Steal      float64
	Guest      float64
	GuestNice  float64
	Interrupts map[string]Interrupt
}

type Interrupt struct {
	Interrupt  int
	Type       string
	DeviceName string
}

type CpuStat struct {
	Timestamp time.Time

	Intr         int
	Ctx          int
	Btime        int
	Processes    int
	ProcsRunning int
	ProcsBlocked int
	Softirq      int

	IntrPerSec      int
	CtxPerSec       int
	BtimePerSec     int
	ProcessesPerSec int
	SoftirqPerSec   int

	CpuProcessorStats []CpuProcessorStat
}

// func Read() {
// 	timestamp := time.Now()
// 	if reader.tmpCpuStat == nil {
// 		reader.tmpCpuStat, _ = reader.read()
// 	} else {
// 		tmpCpuStat, tmpCpuProcessorStats := reader.read()
// 		if len(reader.cpuStats) > reader.cacheLength {
// 			reader.cpuStats = reader.cpuStats[1:]
// 		}
// 		reader.cpuStats = append(reader.cpuStats, CpuStat{
// 			Timestamp:    timestamp,
// 			Intr:         tmpCpuStat.Intr - reader.tmpCpuStat.Intr,
// 			Ctx:          tmpCpuStat.Ctx - reader.tmpCpuStat.Ctx,
// 			Btime:        tmpCpuStat.Btime - reader.tmpCpuStat.Btime,
// 			Processes:    tmpCpuStat.Processes - reader.tmpCpuStat.Processes,
// 			ProcsRunning: tmpCpuStat.ProcsRunning,
// 			ProcsBlocked: tmpCpuStat.ProcsBlocked,
// 			Softirq:      tmpCpuStat.Softirq - reader.tmpCpuStat.Softirq,
// 		})
//
// 	}
// }

func GetCpuStat() (cpuStat *CpuStat, err error) {
	var tmpReader *bufio.Reader
	timestamp := time.Now()

	var cpuProcessorStats []CpuProcessorStat

	// Read /proc/cpuinfo
	var cpuinfo *os.File
	if cpuinfo, err = os.Open("/proc/cpuinfo"); err != nil {
		return
	}
	defer cpuinfo.Close()
	tmpReader = bufio.NewReader(cpuinfo)

	var tmpProcessor int
	var tmpBytes []byte
	var tmpErr error
	for {
		tmpBytes, _, tmpErr = tmpReader.ReadLine()
		if tmpErr != nil {
			break
		}

		splited := str_utils.SplitSpaceColon(string(tmpBytes))
		if len(splited) < 1 {
			continue
		}

		if splited[0] == "processor" {
			tmpProcessor, _ = strconv.Atoi(splited[1])

			// processorから空行までを読みこむ
			cpuinfo := map[string]string{}
			for {
				tmpBytes, _, tmpErr = tmpReader.ReadLine()
				if tmpErr != nil {
					break
				}
				splited := str_utils.SplitSpaceColon(string(tmpBytes))
				if len(splited) < 1 {
					break
				}
				cpuinfo[splited[0]] = splited[1]
			}

			cpuMhzF, _ := strconv.ParseFloat(cpuinfo["cpu MHz"], 64)
			physicalId, _ := strconv.Atoi(cpuinfo["physical id"])
			coreId, _ := strconv.Atoi(cpuinfo["core id"])

			cpuProcessorStats = append(cpuProcessorStats, CpuProcessorStat{
				Processor:  tmpProcessor,
				PhysicalId: physicalId,
				CoreId:     coreId,
				Mhz:        cpuMhzF,
				Interrupts: map[string]Interrupt{},
			})
		}
	}

	// Read /proc/stat
	//      user   nice system idle    iowait irq softirq steal guest guest_nice
	// cpu  264230 262  60792  8237284 20685  0   2652    0     0     0
	// cpu0 126387 2    30266  4124610 11105  0   1011    0     0     0
	// cpu1 137842 260  30525  4112674 9580   0   1641    0     0     0
	// intr 18316761 ...
	// ctxt 57087643
	// btime 1546819593
	// processes 227393
	// procs_running 1
	// procs_blocked 0
	// softirq 11650881 ...

	f, _ := os.Open("/proc/stat")
	defer f.Close()
	tmpReader = bufio.NewReader(f)

	tmpBytes, _, _ = tmpReader.ReadLine()

	lenCpus := len(cpuProcessorStats)
	for i := 0; i < lenCpus; i++ {
		tmpBytes, _, _ = tmpReader.ReadLine()
		cpu := strings.Split(string(tmpBytes), " ")
		user, _ := strconv.Atoi(cpu[1])
		nice, _ := strconv.Atoi(cpu[2])
		system, _ := strconv.Atoi(cpu[3])
		idle, _ := strconv.Atoi(cpu[4])
		iowait, _ := strconv.Atoi(cpu[5])
		irq, _ := strconv.Atoi(cpu[6])
		softirq, _ := strconv.Atoi(cpu[7])
		steal, _ := strconv.Atoi(cpu[8])
		guest, _ := strconv.Atoi(cpu[9])
		guestNice, _ := strconv.Atoi(cpu[10])

		total := float64(user + nice + nice + system + idle + iowait + irq + softirq + steal + guest + guestNice)

		cpuProcessorStats[i].User = float64(user) * 100 / total
		cpuProcessorStats[i].Nice = float64(nice) * 100 / total
		cpuProcessorStats[i].System = float64(system) * 100 / total
		cpuProcessorStats[i].Idle = float64(idle) * 100 / total
		cpuProcessorStats[i].Iowait = float64(iowait) * 100 / total
		cpuProcessorStats[i].Irq = float64(irq) * 100 / total
		cpuProcessorStats[i].Softirq = float64(softirq) * 100 / total
		cpuProcessorStats[i].Steal = float64(steal) * 100 / total
		cpuProcessorStats[i].Guest = float64(guest) * 100 / total
		cpuProcessorStats[i].GuestNice = float64(guestNice) * 100 / total
	}

	tmpBytes, _, _ = tmpReader.ReadLine()
	intr, _ := strconv.Atoi(strings.Split(string(tmpBytes), " ")[1])
	tmpBytes, _, _ = tmpReader.ReadLine()
	ctx, _ := strconv.Atoi(strings.Split(string(tmpBytes), " ")[1])
	tmpBytes, _, _ = tmpReader.ReadLine()
	btime, _ := strconv.Atoi(strings.Split(string(tmpBytes), " ")[1])
	tmpBytes, _, _ = tmpReader.ReadLine()
	processes, _ := strconv.Atoi(strings.Split(string(tmpBytes), " ")[1])
	tmpBytes, _, _ = tmpReader.ReadLine()
	procsRunning, _ := strconv.Atoi(strings.Split(string(tmpBytes), " ")[1])
	tmpBytes, _, _ = tmpReader.ReadLine()
	procsBlocked, _ := strconv.Atoi(strings.Split(string(tmpBytes), " ")[1])
	tmpBytes, _, _ = tmpReader.ReadLine()
	softirq, _ := strconv.Atoi(strings.Split(string(tmpBytes), " ")[1])

	// Read /proc/interrupts
	//             CPU0       CPU1       CPU2       CPU3       CPU4       CPU5       CPU6       CPU7       CPU8       CPU9       CPU10      CPU11
	//    0:         35          0          0          0          0          0          0          0          0          0          0          0  IR-IO-APIC    2-edge      timer
	//    7:          0          0          0          0          0          0          0          0          0          0          0          0  IR-IO-APIC    7-fasteoi   pinctrl_amd
	//    8:          0          0          0          0          0          1          0          0          0          0          0          0  IR-IO-APIC    8-edge      rtc0
	interruptsFile, _ := os.Open("/proc/interrupts")
	tmpReader = bufio.NewReader(interruptsFile)
	_, _, _ = tmpReader.ReadLine() // CPUの行はスキップする
	for {
		tmpBytes, _, tmpErr = tmpReader.ReadLine()
		if tmpErr != nil {
			break
		}
		splitedIntr := str_utils.SplitSpace(string(tmpBytes))
		irqNumber := splitedIntr[0][0 : len(splitedIntr[0])-1]
		if len(splitedIntr) > lenCpus {
			for i := 0; i < lenCpus; i++ {
				intr, _ := strconv.Atoi(splitedIntr[i+1])
				cpuProcessorStats[i].Interrupts[irqNumber] = Interrupt{
					Interrupt:  intr,
					Type:       splitedIntr[lenCpus+1],
					DeviceName: splitedIntr[lenCpus+2],
				}
			}
		} else {
			for i := 0; i < len(splitedIntr)-1; i++ {
				intr, _ := strconv.Atoi(splitedIntr[i+1])
				cpuProcessorStats[i].Interrupts[irqNumber] = Interrupt{
					Interrupt: intr,
				}
			}
		}
	}

	cpuStat = &CpuStat{
		Timestamp:         timestamp,
		Intr:              intr,
		Ctx:               ctx,
		Btime:             btime,
		Processes:         processes,
		ProcsRunning:      procsRunning,
		ProcsBlocked:      procsBlocked,
		Softirq:           softirq,
		CpuProcessorStats: cpuProcessorStats,
	}

	return
}

// func (reader *CpuReader) ReportMetrics() (metrics []spec.ResourceMetric) {
// 	metrics = make([]spec.ResourceMetric, 0, len(reader.cpuStats))
//
// 	for _, stat := range reader.cpuStats {
// 		if stat.ReportStatus == ReportStatusReported {
// 			continue
// 		}
//
// 		if reader.checkProcsRunningWarnLimit > 0 && stat.ProcsRunning > reader.checkProcsRunningWarnLimit {
// 			reader.checkProcsRunningWarnCounter += 1
// 		} else {
// 			reader.checkProcsRunningWarnCounter = 0
// 		}
// 		if reader.checkProcsRunningCritLimit > 0 && stat.ProcsRunning > reader.checkProcsRunningCritLimit {
// 			reader.checkProcsRunningCritCounter += 1
// 		} else {
// 			reader.checkProcsRunningCritCounter = 0
// 		}
//
// 		if reader.checkProcsBlockedWarnLimit > 0 && stat.ProcsBlocked > reader.checkProcsBlockedWarnLimit {
// 			reader.checkProcsBlockedWarnCounter += 1
// 		} else {
// 			reader.checkProcsBlockedWarnCounter = 0
// 		}
// 		if reader.checkProcsBlockedCritLimit > 0 && stat.ProcsBlocked > reader.checkProcsBlockedCritLimit {
// 			reader.checkProcsBlockedCritCounter += 1
// 		} else {
// 			reader.checkProcsBlockedCritCounter = 0
// 		}
//
// 		metrics = append(metrics, spec.ResourceMetric{
// 			Name: "system_cpu",
// 			Time: stat.Timestamp,
// 			Metric: map[string]interface{}{
// 				"intr":          stat.Intr,
// 				"ctx":           stat.Ctx,
// 				"btime":         stat.Btime,
// 				"processes":     stat.Processes,
// 				"procs_running": stat.ProcsRunning,
// 				"procs_blocked": stat.ProcsBlocked,
// 				"softirq":       stat.Softirq,
// 			},
// 		})
// 	}
//
// 	for _, stat := range reader.cpuProcessorStats {
// 		if stat.ReportStatus == ReportStatusReported {
// 			continue
// 		}
//
// 		metrics = append(metrics, spec.ResourceMetric{
// 			Name: "system_processor",
// 			Time: stat.Timestamp,
// 			Tag: map[string]string{
// 				"processor": strconv.FormatInt(stat.Processor, 10),
// 			},
// 			Metric: map[string]interface{}{
// 				"mhz":       stat.Mhz,
// 				"user":      stat.User,
// 				"nice":      stat.Nice,
// 				"system":    stat.System,
// 				"idle":      stat.Idle,
// 				"iowait":    stat.Iowait,
// 				"irq":       stat.Irq,
// 				"softirq":   stat.Softirq,
// 				"steal":     stat.Steal,
// 				"guest":     stat.Guest,
// 				"guestnice": stat.GuestNice,
// 			},
// 		})
// 	}
//
// 	return
// }
//
// func (reader *CpuReader) ReportEvents() (events []spec.ResourceEvent) {
// 	if len(reader.cpuStats) == 0 {
// 		return
// 	}
//
// 	cpuStat := reader.cpuStats[len(reader.cpuStats)-1]
// 	eventProcsRunningLevel := consts.EventLevelSuccess
// 	if reader.checkProcsRunningCritCounter > reader.checkProcsRunningOccurences {
// 		eventProcsRunningLevel = consts.EventLevelCritical
// 	} else if reader.checkProcsRunningWarnCounter > reader.checkProcsRunningOccurences {
// 		eventProcsRunningLevel = consts.EventLevelWarning
// 	}
//
// 	events = append(events, spec.ResourceEvent{
// 		Name:            "CheckCpuProcsRunning",
// 		Time:            cpuStat.Timestamp,
// 		Level:           eventProcsRunningLevel,
// 		Msg:             fmt.Sprintf("ProcsRunning: %d", cpuStat.ProcsRunning),
// 		ReissueDuration: reader.checkProcsRunningReissueDuration,
// 	})
//
// 	eventProcsBlockedLevel := consts.EventLevelSuccess
// 	if reader.checkProcsBlockedCritCounter > reader.checkProcsBlockedOccurences {
// 		eventProcsBlockedLevel = consts.EventLevelCritical
// 	} else if reader.checkProcsBlockedWarnCounter > reader.checkProcsBlockedOccurences {
// 		eventProcsBlockedLevel = consts.EventLevelWarning
// 	}
// 	events = append(events, spec.ResourceEvent{
// 		Name:            "CheckCpuProcsBlocked",
// 		Time:            cpuStat.Timestamp,
// 		Level:           eventProcsBlockedLevel,
// 		Msg:             fmt.Sprintf("ProcsBlocked: %d", cpuStat.ProcsBlocked),
// 		ReissueDuration: reader.checkProcsBlockedReissueDuration,
// 	})
//
// 	return
// }
//
// func (reader *CpuReader) Reported() {
// 	for i := range reader.cpuStats {
// 		reader.cpuStats[i].ReportStatus = ReportStatusReported
// 	}
// 	for i := range reader.cpuProcessorStats {
// 		reader.cpuProcessorStats[i].ReportStatus = ReportStatusReported
// 	}
// 	return
// }
