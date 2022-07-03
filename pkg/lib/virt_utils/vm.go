package virt_utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/syunkitada/goapp2/pkg/lib/errors"
	"github.com/syunkitada/goapp2/pkg/lib/logger"
	"github.com/syunkitada/goapp2/pkg/lib/str_utils"
)

type VmResources []VmResource

func (self VmResources) String() string {
	tableString, table := str_utils.GetTable()
	table.SetHeader([]string{"Kind", "Name", "Status", "Vcpus", "MemoryMb", "DiskGb", "Image", "Ports"})
	for _, r := range self {
		s := r.Spec

		portsStrs := []string{}
		for _, port := range s.NetworkPorts {
			portsStrs = append(portsStrs, port.Ip)
		}

		table.Append([]string{
			r.Kind, s.Name, s.Status, strconv.Itoa(int(s.Vcpus)), strconv.Itoa(int(s.MemoryMb)), strconv.Itoa(int(s.DiskGb)),
			s.ImageName, strings.Join(portsStrs, ","),
		})
	}
	table.Render()
	return tableString.String()
}

type VmResource struct {
	Kind string
	Spec Vm
}

type VmSpec struct {
	Name     string              `gorm:"not null;uniqueIndex:udx_name;" validate:"required"`
	Kind     string              `gorm:"not null;" validate:"required"`
	Vcpus    uint                `gorm:"not null;" validate:"required"`
	MemoryMb uint                `gorm:"not null;" validate:"required"`
	DiskGb   uint                `gorm:"not null;" validate:"required"`
	Image    ImageDetectSpec     `gorm:"-"`
	Networks []NetworkDetectSpec `gorm:"-"`
	Spec     interface{}         `gorm:"-"`
}

type Vm struct {
	VmSpec
	VmImage
	Id           uint            `gorm:"not null;primaryKey;autoIncrement;"`
	DeletedAt    *time.Time      `gorm:"uniqueIndex:udx_name;"`
	ImageId      uint            `gorm:"not null;`
	SpecStr      string          `gorm:"not null;column:spec" json:"-"`
	Status       string          `gorm:"not null;"`
	NetworkPorts []VmNetworkPort `gorm:"-"`
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
				var image *Image
				if image, err = self.DetectImage(tctx, tx, &spec.Image); err != nil {
					return
				}

				vm := &Vm{
					VmSpec:  *spec,
					ImageId: image.Id,
					Status:  StatusCreated,
				}
				if err = tx.Create(vm).Error; err != nil {
					return
				}

				if _, err = self.AssignNetworkPorts(tctx, tx, vm, spec.Networks); err != nil {
					return
				}

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
	sql := self.sqlClient.DB.Table("vms AS v").Select("v.*, i.name as image_name, i.kind as image_kind, i.spec as image_spec_str").
		Joins("INNER JOIN images AS i ON v.image_id == i.id").
		Where("v.deleted_at IS NULL")
	if len(names) > 0 {
		sql = sql.Where("name in (?)", names)
	}
	if err = sql.Scan(&vms).Error; err != nil {
		return
	}

	var ports []VmNetworkPort
	sql = self.sqlClient.DB.Table("network_ports").Select("*")
	if err = sql.Scan(&ports).Error; err != nil {
		return
	}

	var networks []Network
	sql = self.sqlClient.DB.Table("networks").Select("*")
	if err = sql.Scan(&networks).Error; err != nil {
		return
	}

	vmMap := map[uint]*Vm{}
	for i := range vms {
		vm := vms[i]
		var imageUrlSpec ImageUrlSpec
		switch vm.ImageKind {
		case KindImageUrl:
			if err = json.Unmarshal([]byte(vm.ImageSpecStr), &imageUrlSpec); err != nil {
				return
			}
		}
		vm.imageUrlSpec = imageUrlSpec
		vmMap[vm.Id] = &vm
	}

	networkMap := map[uint]*VmNetwork{}
	for i := range networks {
		network := networks[i]
		var networkLocalSpec NetworkLocalSpec
		switch network.Kind {
		case KindNetworkLocal:
			if err = json.Unmarshal([]byte(network.SpecStr), &networkLocalSpec); err != nil {
				return
			}
		}
		networkMap[network.Id] = &VmNetwork{
			Network:          network,
			networkLocalSpec: networkLocalSpec,
		}
	}

	for i := range ports {
		port := ports[i]
		if vm, ok := vmMap[port.VmId]; ok {
			network, ok := networkMap[port.NetworkId]
			if !ok {
				err = fmt.Errorf("port's network is not found")
				return
			}
			port.VmNetwork = network
			vm.NetworkPorts = append(vm.NetworkPorts, port)
		}
	}

	for _, vm := range vmMap {
		vmResources = append(vmResources, VmResource{
			Kind: KindVm,
			Spec: *vm,
		})
	}

	return
}

type netnsPort struct {
	Id           uint
	Name         string
	NetnsGateway string
	NetnsIp      string
	VmIp         string
	VmMac        string
	VmSubnet     string
	Kind         string
}

func (self *VirtController) StartVmResources(tctx *logger.TraceContext, names []string) (vmResources VmResources, err error) {
	if vmResources, err = self.GetVmResources(tctx, names); err != nil {
		return
	}

	if err = self.PrepareImages(tctx, vmResources); err != nil {
		return
	}

	if err = self.PrepareNetworks(tctx, vmResources); err != nil {
		return
	}

	// TODO start vm
	return
}
