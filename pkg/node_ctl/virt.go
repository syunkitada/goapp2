package node_ctl

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/syunkitada/goapp2/pkg/lib/file_utils"
	"github.com/syunkitada/goapp2/pkg/lib/logger"
	"github.com/syunkitada/goapp2/pkg/lib/virt_utils"
)

var files []string
var outputFormat string

var virtCmd = &cobra.Command{
	Use:   "virt",
	Short: "control virt",
}

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "bootstrap",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init(&logger.Config{})
		virtController := virt_utils.NewVirtContoller()
		virtController.MustInit()
		virtController.MustBootstrap()
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var resources [][]byte
		if resources, err = file_utils.ReadFilesBytesFromMultiPath(files); err != nil {
			fmt.Println("Failed", err.Error())
			return
		}

		logger.Init(&logger.Config{})
		virtController := virt_utils.NewVirtContoller()
		virtController.MustInit()
		tctx := logger.NewTraceContext()
		if err = virtController.Create(tctx, resources); err != nil {
			fmt.Println("Failed", err.Error())
			return
		}
	},
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
}

func getResource(kind string, args []string) {
	var err error
	logger.Init(&logger.Config{})
	virtController := virt_utils.NewVirtContoller()
	virtController.MustInit()
	tctx := logger.NewTraceContext()
	var result *virt_utils.GetResult
	if result, err = virtController.Get(tctx, kind, args); err != nil {
		fmt.Println("Failed", err.Error())
		return
	}

	result.Output(outputFormat)
}

func init() {
	createCmd.PersistentFlags().StringSliceVarP(&files, "files", "f", []string{}, "source file")
	createCmd.MarkPersistentFlagRequired("files")

	resources := []string{
		"all",
		"vm",
		"image",
		"network",
	}
	for i := range resources {
		resource := resources[i]
		var getResourceCmd = &cobra.Command{
			Use:   resource + " [name]...",
			Short: "get " + resource + " information",
			Run: func(cmd *cobra.Command, args []string) {
				getResource(resource, args)
			},
		}
		getCmd.AddCommand(getResourceCmd)
	}

	virtCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "output format")
	virtCmd.AddCommand(getCmd)
	virtCmd.AddCommand(bootstrapCmd)
	virtCmd.AddCommand(createCmd)
	rootCmd.AddCommand(virtCmd)

}
