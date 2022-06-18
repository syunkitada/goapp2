package virt_utils

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"

	"github.com/syunkitada/goapp2/pkg/lib/db_utils"
	"github.com/syunkitada/goapp2/pkg/lib/logger"
)

type VirtController struct {
	sqlClient *db_utils.SqlClient
	validate  *validator.Validate
}

func NewVirtContoller() (virtController *VirtController) {
	sqlClient := db_utils.NewSqlClient(&db_utils.Config{})
	return &VirtController{
		sqlClient: sqlClient,
		validate:  validator.New(),
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
			var network VmSpec
			if err = json.Unmarshal(bytes, &network); err != nil {
				return
			}
			if err = self.CreateOrUpdateVm(tctx, &network); err != nil {
				return
			}
		case KindNetwork:
			var network NetworkSpec
			if err = json.Unmarshal(bytes, &network); err != nil {
				return
			}
			if err = self.CreateOrUpdateNetwork(tctx, &network); err != nil {
				return
			}
		case KindImage:
			var spec ImageSpec
			if err = json.Unmarshal(bytes, &spec); err != nil {
				return
			}
			if err = self.CreateOrUpdateImage(tctx, &spec); err != nil {
				return
			}
		}
	}
	return
}

type GetResult struct {
	Vms      []VmResource
	Networks []NetworkResource
	Images   []ImageResource
}

func (self *VirtController) Get(tctx *logger.TraceContext, kind string, args []string) (result *GetResult, err error) {
	var vms []VmResource
	var networks []NetworkResource
	var images []ImageResource

	switch kind {
	case KindVm:
	case KindNetwork:
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
