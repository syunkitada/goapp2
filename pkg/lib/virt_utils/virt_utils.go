package virt_utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"

	"github.com/syunkitada/goapp2/pkg/lib/db_utils"
	"github.com/syunkitada/goapp2/pkg/lib/logger"
	"github.com/syunkitada/goapp2/pkg/lib/str_utils"
)

type VirtController struct {
	conf      VirtControllerConfig
	sqlClient *db_utils.SqlClient
	validate  *validator.Validate

	imagesDir string
	vmsDir    string
}

var virtControllerConf = VirtControllerConfig{
	VarDir: "",
}

func init() {
	home := os.Getenv("HOME")
	virtControllerConf.VarDir = filepath.Join(home, ".cache/goapp2/")
}

type VirtControllerConfig struct {
	VarDir   string
	Database db_utils.Config
}

func NewVirtContoller(conf *VirtControllerConfig) (virtController *VirtController) {
	imagesDir := filepath.Join(virtControllerConf.VarDir, "images")
	vmsDir := filepath.Join(virtControllerConf.VarDir, "vms")

	sqlClient := db_utils.NewSqlClient(&virtControllerConf.Database)

	return &VirtController{
		conf:      virtControllerConf,
		sqlClient: sqlClient,
		validate:  validator.New(),
		imagesDir: imagesDir,
		vmsDir:    vmsDir,
	}
}

func (self *VirtController) MustInit() {
	tctx := logger.NewTraceContext()
	self.sqlClient.MustOpen(tctx)
	return
}

func (self *VirtController) MustBootstrap() {
	var err error
	tctx := logger.NewTraceContext()
	defer func() {
		if err != nil {
			logger.Fatalf(tctx, "FailedBootstrap: err=%s", err.Error())
		}
	}()
	if err = self.BootstrapImage(tctx); err != nil {
		return
	}
	if err = self.BootstrapNetwork(tctx); err != nil {
		return
	}
	if err = self.BootstrapVm(tctx); err != nil {
		return
	}
}

const (
	KindAll     = "all"
	KindVm      = "vm"
	KindImage   = "image"
	KindNetwork = "network"
)

type Resource struct {
	Kind string `validate:"required"`
	Name string `validate:"required"`
	Spec interface{}
}

func (self *VirtController) Create(tctx *logger.TraceContext, resourcesBytes [][]byte) (err error) {
	var vms []VmSpec
	var networks []NetworkSpec
	var images []ImageSpec
	for _, resourceBytes := range resourcesBytes {
		var resource Resource
		if err = yaml.Unmarshal(resourceBytes, &resource); err != nil {
			return
		}
		var bytes []byte
		if bytes, err = json.Marshal(resource.Spec); err != nil {
			return
		}
		switch resource.Kind {
		case KindVm:
			var spec VmSpec
			if err = json.Unmarshal(bytes, &spec); err != nil {
				return
			}
			vms = append(vms, spec)
		case KindNetwork:
			var spec NetworkSpec
			if err = json.Unmarshal(bytes, &spec); err != nil {
				return
			}
			networks = append(networks, spec)
		case KindImage:
			var spec ImageSpec
			if err = json.Unmarshal(bytes, &spec); err != nil {
				return
			}
			images = append(images, spec)
		}
	}

	for i := range networks {
		if err = self.CreateOrUpdateNetwork(tctx, &networks[i]); err != nil {
			return
		}
	}

	for i := range images {
		if err = self.CreateOrUpdateImage(tctx, &images[i]); err != nil {
			return
		}
	}

	for i := range vms {
		if err = self.CreateOrUpdateVm(tctx, &vms[i]); err != nil {
			return
		}
	}
	return

}

type GetResult struct {
	Vms      VmResources
	Networks NetworkResources
	Images   ImageResources
}

func (self *GetResult) Output(format string) {
	switch format {
	case "yaml", "json":
		outputs := []string{}
		for _, data := range self.Vms {
			outputs = str_utils.AppendOutputByFormat(outputs, data, format)
		}
		for _, data := range self.Images {
			outputs = str_utils.AppendOutputByFormat(outputs, data, format)
		}
		for _, data := range self.Networks {
			outputs = str_utils.AppendOutputByFormat(outputs, data, format)
		}
		fmt.Println(strings.Join(outputs, "\n---\n\n"))
	default:
		if len(self.Vms) > 0 {
			str_utils.OutputByFormat(self.Vms, format)
		}
		if len(self.Images) > 0 {
			str_utils.OutputByFormat(self.Images, format)
		}
		if len(self.Networks) > 0 {
			str_utils.OutputByFormat(self.Networks, format)
		}
	}

}

func (self *VirtController) Get(tctx *logger.TraceContext, kind string, args []string) (result *GetResult, err error) {
	var vms VmResources
	var networks NetworkResources
	var images ImageResources

	switch kind {
	case KindAll:
		if vms, err = self.GetVmResources(tctx, args); err != nil {
			return
		}
		if networks, err = self.GetNetworkResources(tctx, args); err != nil {
			return
		}
		if images, err = self.GetImageResources(tctx, args); err != nil {
			return
		}
	case KindVm:
		if vms, err = self.GetVmResources(tctx, args); err != nil {
			return
		}
	case KindNetwork:
		if networks, err = self.GetNetworkResources(tctx, args); err != nil {
			return
		}
	case KindImage:
		if images, err = self.GetImageResources(tctx, args); err != nil {
			return
		}
	}

	result = &GetResult{
		Vms:      vms,
		Networks: networks,
		Images:   images,
	}
	return
}

func (self *VirtController) Start(tctx *logger.TraceContext, kind string, args []string) (result *GetResult, err error) {
	var vms VmResources

	switch kind {
	case KindVm:
		if vms, err = self.StartVmResources(tctx, args); err != nil {
			return
		}
	}

	result = &GetResult{
		Vms: vms,
	}
	return
}
