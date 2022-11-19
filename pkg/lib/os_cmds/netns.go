package os_cmds

import (
	"fmt"
	"strings"

	"github.com/syunkitada/goapp2/pkg/lib/cmd_runner"
	"github.com/syunkitada/goapp2/pkg/lib/logger"
)

func GetNetnsSet(tctx *logger.TraceContext) (netnsSet map[string]bool, err error) {
	var result *cmd_runner.Result
	if result, err = cmd_runner.Run(&cmd_runner.Config{Cmd: "ip netns"}); err != nil {
		return
	}
	if result.Status != 0 {
		err = fmt.Errorf("Failed ip netns: output=%s, status=%s", result.Output, result.Status)
		return
	}

	netnsSet = map[string]bool{}
	for _, line := range strings.Split(result.Output, "\n") {
		if line != "" {
			netnsSet[strings.Split(line, " ")[0]] = true
		}
	}

	return
}

func AddNetns(tctx *logger.TraceContext, netns string) (err error) {
	_, err = cmd_runner.Run(&cmd_runner.Config{Cmd: fmt.Sprintf("ip netns add %s", netns)})
	return
}

func ExecInIpNetns(tctx *logger.TraceContext, netns string, cmd string) (out string, err error) {
	var result *cmd_runner.Result
	if result, err = cmd_runner.Run(&cmd_runner.Config{Cmd: fmt.Sprintf("ip netns exec %s %s", netns, cmd)}); err != nil {
		return
	}
	out = result.Output
	return
}
