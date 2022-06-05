package node_ctl

import (
	"fmt"

	"github.com/spf13/cobra"
)

var file string

var virtCmd = &cobra.Command{
	Use:   "virt",
	Short: "control virt",
}

var startCmd = &cobra.Command{
	Use:   "create",
	Short: "create",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("DEBUG virt", file)
	},
}

func init() {
	startCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "source file")
	startCmd.MarkPersistentFlagRequired("file")

	virtCmd.AddCommand(startCmd)
	rootCmd.AddCommand(virtCmd)
}
