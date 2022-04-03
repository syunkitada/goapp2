package process_utils

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetProcesses(t *testing.T) {
	a := assert.New(t)
	a.Equal(true, true)

	wd, err := os.Getwd()
	a.NoError(err)
	rootDir := wd + "/testdata/root/"

	beforeTime := time.Now().Add(-1 * time.Second)
	processes, pidIndexMap, err := GetProcesses(rootDir)
	a.NoError(err)
	afterTime := time.Now().Add(1 * time.Second)

	expectedPidIntMap := map[int]int{}

	for i, process := range processes {
		// Timestampにはtime.Now()が入るので、別途評価する
		a.Greater(process.Stat.Timestamp.Unix(), beforeTime.Unix())
		a.Greater(afterTime.Unix(), process.Stat.Timestamp.Unix())
		processes[i].Stat.Timestamp = time.Time{}

		expectedPidIntMap[process.Pid] = i
	}

	expectedProcesses := []Process{
		Process{
			Name: "systemd",
			Pid:  1,
			Tgid: 1,
			Ppid: 0,
			Cmd:  "",
			Cmds: []string{
				"/sbin/init",
				"splash",
			},
			State: 1,
			Stat: ProcessStat{
				VmSizeKb:                 225752,
				VmRssKb:                  9412,
				State:                    0,
				SchedCpuTime:             4670045920,
				SchedWaitTime:            622418398,
				SchedTimeSlices:          23374,
				HugetlbPages:             0,
				Threads:                  1,
				VoluntaryCtxtSwitches:    22527,
				NonvoluntaryCtxtSwitches: 435,
				Utime:                    181,
				Stime:                    285,
				Gtime:                    0,
				Cgtime:                   0,
				StartTime:                3,
				Syscr:                    1065760,
				Syscw:                    574769,
				ReadBytes:                1707520000,
				WriteBytes:               1873141760,
			},
		},
		Process{
			Name:  "kthreadd",
			Cmds:  []string{},
			State: 1,
			Pid:   2,
			Tgid:  2,
			Stat: ProcessStat{
				State:                    0,
				SchedCpuTime:             8758561,
				SchedWaitTime:            228198,
				SchedTimeSlices:          385,
				HugetlbPages:             0,
				Threads:                  1,
				VoluntaryCtxtSwitches:    385,
				NonvoluntaryCtxtSwitches: 0,
				Cgtime:                   0,
				StartTime:                3,
			},
		},
		Process{
			Name: "process_utils.t",
			Pid:  21607,
			Tgid: 21607,
			Ppid: 21401,
			Cmd:  "",
			Cmds: []string{
				"/tmp/go-build033075601/b001/process_utils.test",
				"-test.timeout=10m0s",
				"-test.v=true",
				"-test.coverprofile=/tmp/go-build033075601/b001/_cover_.out",
				"-test.run=TestGetProcesses",
			},
			Children: []int{21613},
			State:    1,
			Stat: ProcessStat{
				VmSizeKb:                 1446744,
				VmRssKb:                  13716,
				State:                    0,
				SchedCpuTime:             12035808,
				SchedWaitTime:            15679,
				SchedTimeSlices:          4,
				HugetlbPages:             0,
				Threads:                  6,
				VoluntaryCtxtSwitches:    5,
				NonvoluntaryCtxtSwitches: 2,
				Utime:                    0,
				Stime:                    0,
				Gtime:                    0,
				Cgtime:                   0,
				StartTime:                1408738,
				Syscr:                    14,
				Syscw:                    1,
				ReadBytes:                0,
				WriteBytes:               0,
			},
		},
		Process{
			Name: "sleep",
			Pid:  21613,
			Tgid: 21613,
			Ppid: 21607,
			Cmd:  "",
			Cmds: []string{
				"sleep",
				"1000",
			},
			State: 1,
			Stat: ProcessStat{
				VmSizeKb:                 8496,
				VmRssKb:                  856,
				State:                    0,
				SchedCpuTime:             610364,
				SchedWaitTime:            17503,
				SchedTimeSlices:          2,
				HugetlbPages:             0,
				Threads:                  1,
				VoluntaryCtxtSwitches:    1,
				NonvoluntaryCtxtSwitches: 1,
				Utime:                    0,
				Stime:                    0,
				Gtime:                    0,
				Cgtime:                   0,
				StartTime:                1408739,
				Syscr:                    6,
				Syscw:                    0,
				ReadBytes:                0,
				WriteBytes:               0,
			},
		},
	}
	a.ElementsMatch(expectedProcesses, processes)

	a.Equal(expectedPidIntMap, pidIndexMap)

	{
		// rootがない
		_, _, err := GetProcesses(wd + "/testdata/none/")
		a.Error(err)
	}

	{
		// procがディレクトリではなくファイル
		_, _, err := GetProcesses(wd + "/testdata/exception_file_proc/")
		a.Error(err)
	}

	{
		// cmdlineがない
		_, _, err := GetProcesses(wd + "/testdata/exception_no_cmdline/")
		a.Error(err)
	}

	{
		// statusがない
		_, _, err := GetProcesses(wd + "/testdata/exception_no_status/")
		a.Error(err)
	}

	{
		// schedstatがない
		_, _, err := GetProcesses(wd + "/testdata/exception_no_schedstat/")
		a.Error(err)
	}

	{
		// statがない
		_, _, err := GetProcesses(wd + "/testdata/exception_no_stat/")
		a.Error(err)
	}

	{
		// ioがない
		_, _, err := GetProcesses(wd + "/testdata/exception_no_io/")
		a.Error(err)
	}

	{
		// schedstatのフォーマットが間違ってる
		_, _, err := GetProcesses(wd + "/testdata/exception_invalid_schedstat/")
		a.Error(err)
	}
}
