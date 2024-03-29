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

var virtController *virt_utils.VirtController

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "bootstrap",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init(&logger.Config{})
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
	tctx := logger.NewTraceContext()
	var result *virt_utils.GetResult
	if result, err = virtController.Get(tctx, kind, args); err != nil {
		fmt.Println("Failed", err.Error())
		return
	}

	result.Output(outputFormat)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start",
}

func startResource(kind string, args []string) {
	var err error
	logger.Init(&logger.Config{})
	tctx := logger.NewTraceContext()
	var result *virt_utils.GetResult
	if result, err = virtController.Start(tctx, kind, args); err != nil {
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

	ctlResources := []string{
		"vm",
	}
	for i := range ctlResources {
		resource := ctlResources[i]
		var startResourceCmd = &cobra.Command{
			Use:   resource + " [name]...",
			Short: "start " + resource + " information",
			Run: func(cmd *cobra.Command, args []string) {
				startResource(resource, args)
			},
		}
		startCmd.AddCommand(startResourceCmd)
	}

	virtCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "output format")
	virtCmd.AddCommand(getCmd)
	virtCmd.AddCommand(startCmd)
	virtCmd.AddCommand(bootstrapCmd)
	virtCmd.AddCommand(createCmd)
	rootCmd.AddCommand(virtCmd)

	conf := virt_utils.VirtControllerConfig{}
	virtController = virt_utils.NewVirtContoller(&conf)
	virtController.MustInit()
}
