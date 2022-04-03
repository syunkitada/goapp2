package node_ctl

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/syunkitada/goapp2/pkg/lib/process_utils"
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "process",
	Run: func(cmd *cobra.Command, args []string) {
		processes, _, err := process_utils.GetProcesses("/", true)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Pid", "Name", "Cmd"})
		for _, process := range processes {
			table.Append([]string{strconv.Itoa(process.Pid), process.Name, strings.Join(process.Cmds, " ")})
		}
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}
