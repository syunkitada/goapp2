package virt_utils

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/syunkitada/goapp2/pkg/lib/errors"
	"github.com/syunkitada/goapp2/pkg/lib/logger"
)

type NetworkResource struct {
	Kind string
	Spec Network
}

type Network struct {
	NetworkSpec
	Id        uint       `gorm:"not null;primaryKey;autoIncrement;"`
	DeletedAt *time.Time `gorm:"uniqueIndex:udx_name;"`
}

type NetworkSpec struct {
	Name      string            `gorm:"not null;uniqueIndex:udx_name;"`
	Kind      string            `gorm:"not null;"`
	Subnet    string            `gorm:"-"`
	StartIp   string            `gorm:"-"`
	EndIp     string            `gorm:"-"`
	Gateway   string            `gorm:"-"`
	Resolvers []NetworkResolver `gorm:"-"`
	Nat       NetworkNat        `gorm:"-"`
}

type NetworkResolver struct {
	Resolver string
}

type NetworkNat struct {
	Enable bool
	Ports  string
}

type NewworkPort struct {
}

func (self *VirtController) BootstrapNetwork(tctx *logger.TraceContext) (err error) {
	if err = self.sqlClient.DB.AutoMigrate(&Network{}).Error; err != nil {
		return
	}
	return
}

func (self *VirtController) CreateOrUpdateNetwork(tctx *logger.TraceContext, spec *NetworkSpec) (err error) {
	if err = self.validate.Struct(spec); err != nil {
		return
	}

	var network *Network
	if network, err = self.GetNetwork(spec.Name); err != nil {
		if errors.IsNotFoundError(err) {
			err = self.sqlClient.Transact(tctx, func(tx *gorm.DB) (err error) {
				// create
				return
			})
		}
		return
	} else {
		fmt.Println("debug update network", network)
		return
	}
	return
}

func (self *VirtController) GetNetwork(name string) (network *Network, err error) {
	var networks []Network
	sql := self.sqlClient.DB.Table("networks").Select("*").Where("deleted_at IS NULL")
	if err = sql.Scan(&networks).Error; err != nil {
		return
	}
	if len(networks) > 1 {
		err = errors.NewConflictErrorf("duplicated networks are found: name=%s, len=%d", name, len(networks))
		return
	} else if len(networks) == 0 {
		err = errors.NewNotFoundErrorf("network is not found: name=%s", name)
		return
	}
	network = &networks[0]
	return
}
