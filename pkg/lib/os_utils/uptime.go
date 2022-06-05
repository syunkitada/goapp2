package os_utils

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type UptimeStat struct {
	Uptime int
}

func GetUptimeStat(rootDir string) (uptimeStat *UptimeStat, err error) {
	var procUptime *os.File
	if procUptime, err = os.Open(rootDir + "proc/uptime"); err != nil {
		return
	}
	defer procUptime.Close()
	tmpReader := bufio.NewReader(procUptime)
	tmpBytes, _, _ := tmpReader.ReadLine()
	uptimeWords := strings.Split(string(tmpBytes), " ")
	uptimeF, _ := strconv.ParseFloat(uptimeWords[0], 64)
	uptime := int(uptimeF)

	uptimeStat = &UptimeStat{
		Uptime: uptime,
	}
	return
}
