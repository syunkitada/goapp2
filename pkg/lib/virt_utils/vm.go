package virt_utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/syunkitada/goapp2/pkg/lib/errors"
	"github.com/syunkitada/goapp2/pkg/lib/logger"
	"github.com/syunkitada/goapp2/pkg/lib/str_utils"
)

type VmResources []VmResource

func (self VmResources) String() string {
	tableString, table := str_utils.GetTable()
	table.SetHeader([]string{"Kind", "Name", "Vcpus", "MemoryMb", "DiskGb"})
	for _, r := range self {
		s := r.Spec
		table.Append([]string{r.Kind, s.Name, strconv.Itoa(int(s.Vcpus)), strconv.Itoa(int(s.MemoryMb)), strconv.Itoa(int(s.DiskGb))})
	}
	table.Render()
	return tableString.String()
}

type VmResource struct {
	Kind string
	Spec Vm
}

type VmSpec struct {
	Name         string              `gorm:"not null;uniqueIndex:udx_name;" validate:"required"`
	Kind         string              `gorm:"not null;" validate:"required"`
	Vcpus        uint                `gorm:"not null;" validate:"required"`
	MemoryMb     uint                `gorm:"not null;" validate:"required"`
	DiskGb       uint                `gorm:"not null;" validate:"required"`
	Image        ImageDetectSpec     `gorm:"-"`
	Networks     []NetworkDetectSpec `gorm:"-"`
	NetworkPorts []NetworkPort       `gorm:"-"`
	Spec         interface{}         `gorm:"-"`
}

type Vm struct {
	VmSpec
	Id        uint          `gorm:"not null;primaryKey;autoIncrement;"`
	DeletedAt *time.Time    `gorm:"uniqueIndex:udx_name;"`
	SpecStr   string        `gorm:"not null;column:spec" json:"-"`
	Image     ImageSpec     `gorm:"-"`
	Networks  []NetworkSpec `gorm:"-"`
}

type VmQemuSpec struct {
	Service SystemdService `gorm:"-"`
}

type SystemdService struct {
	Restart string `oneof=always`
}

const (
	KindVmQemu = "qemu"
)

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

	var specBytes []byte
	if specBytes, err = json.Marshal(spec.Spec); err != nil {
		return
	}

	var vmQemuSpec VmQemuSpec
	switch spec.Kind {
	case KindVmQemu:
		if err = json.Unmarshal(specBytes, &vmQemuSpec); err != nil {
			return
		}
		if err = self.validate.Struct(vmQemuSpec); err != nil {
			return
		}
	default:
		err = errors.NewBadInputErrorf("invalid image kind: kind=%s", spec.Kind)
		return
	}

	var vm *Vm
	if vm, err = self.GetVm(spec.Name); err != nil {
		if errors.IsNotFoundError(err) {
			err = self.sqlClient.Transact(tctx, func(tx *gorm.DB) (err error) {
				fmt.Println("DEBUG create vm")
				var ports []NetworkPort
				if ports, err = self.AssignNetworkPort(tctx, tx, spec.Networks); err != nil {
					return
				}
				fmt.Println("DEBUG ports", ports)

				vm := Vm{
					VmSpec:  *spec,
					SpecStr: string(specBytes),
				}
				if err = tx.Create(&vm).Error; err != nil {
					return
				}

				err = fmt.Errorf("DEBUG error")
				return
			})
		}
		return
	} else {
		if string(specBytes) != vm.Spec {
			if err = self.sqlClient.DB.Table("vms").Where("id = ?", vm.Id).Updates(map[string]interface{}{
				"spec": string(specBytes),
			}).Error; err != nil {
				return
			}
		}
	}
	return
}

func (self *VirtController) GetVm(name string) (vm *Vm, err error) {
	var vms []Vm
	sql := self.sqlClient.DB.Table("vms").Select("*").Where("deleted_at IS NULL")
	if err = sql.Scan(&vms).Error; err != nil {
		return
	}
	if len(vms) > 1 {
		err = errors.NewConflictErrorf("duplicated vms are found: name=%s, len=%d", name, len(vms))
		return
	} else if len(vms) == 0 {
		err = errors.NewNotFoundErrorf("vm is not found: name=%s", name)
		return
	}
	vm = &vms[0]
	return
}

func (self *VirtController) GetVmResources(tctx *logger.TraceContext, names []string) (vmResources VmResources, err error) {
	var vms []Vm
	sql := self.sqlClient.DB.Table("vms").Select("*").Where("deleted_at IS NULL")
	if len(names) > 0 {
		sql = sql.Where("name in (?)", names)
	}
	if err = sql.Scan(&vms).Error; err != nil {
		return
	}

	for _, vm := range vms {
		vmResources = append(vmResources, VmResource{
			Kind: KindVm,
			Spec: vm,
		})
	}

	return
}
