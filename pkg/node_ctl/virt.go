package node_ctl

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/syunkitada/goapp2/pkg/lib/file_utils"
	"github.com/syunkitada/goapp2/pkg/lib/logger"
	"github.com/syunkitada/goapp2/pkg/lib/virt_utils"
)

var files []string

var virtCmd = &cobra.Command{
	Use:   "virt",
	Short: "control virt",
}

var startCmd = &cobra.Command{
	Use:   "create",
	Short: "create",
	Run: func(cmd *cobra.Command, args []string) {
		var resources []interface{}
		resources, err := file_utils.ReadFilesFromMultiPath(files)
		if err != nil {
			fmt.Println("Failed", err.Error())
			return
		}

		logger.Init(&logger.Config{})
		virtController := virt_utils.NewVirtContoller()
		virtController.Init()
		fmt.Println("DEBUG virt", len(resources))
	},
}

func init() {
	startCmd.PersistentFlags().StringSliceVarP(&files, "files", "f", []string{}, "source file")
	startCmd.MarkPersistentFlagRequired("files")

	virtCmd.AddCommand(startCmd)
	rootCmd.AddCommand(virtCmd)
}
