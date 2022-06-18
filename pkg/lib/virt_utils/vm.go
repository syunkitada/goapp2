package virt_utils

import (
	"time"

	"github.com/syunkitada/goapp2/pkg/lib/logger"
)

type VmResource struct {
	Kind string
	Spec Vm
}

type VmSpec struct {
	Name     string         `gorm:"not null;uniqueIndex:udx_name;" validate:"required"`
	Kind     string         `gorm:"not null;" validate:"required"`
	Vcpus    uint           `gorm:"not null;" validate:"required"`
	MemoryMb uint           `gorm:"not null;" validate:"required"`
	DiskGb   uint           `gorm:"not null;" validate:"required"`
	Image    ImageSpec      `gorm:"-"`
	Networks []Network      `gorm:"-"`
	Service  SystemdService `gorm:"-"`
}

type Vm struct {
	VmSpec
	Id        uint       `gorm:"not null;primaryKey;autoIncrement;"`
	DeletedAt *time.Time `gorm:"uniqueIndex:udx_name;"`
}

type SystemdService struct {
	Restart string `oneof=always`
}

func (self *VirtController) BootstrapVm(tctx *logger.TraceContext) (err error) {
	if err = self.sqlClient.DB.AutoMigrate(&Vm{}).Error; err != nil {
		return
	}
	return
}

func (self *VirtController) CreateOrUpdateVm(tctx *logger.TraceContext, spec *VmSpec) (err error) {
	if err = self.validate.Struct(spec); err != nil {
		return
	}
	return
}
