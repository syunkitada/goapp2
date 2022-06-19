package virt_utils

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/syunkitada/goapp2/pkg/lib/errors"
	"github.com/syunkitada/goapp2/pkg/lib/logger"
	"github.com/syunkitada/goapp2/pkg/lib/str_utils"
)

type NetworkResources []NetworkResource

func (self NetworkResources) String() string {
	tableString, table := str_utils.GetTable()
	table.SetHeader([]string{"Kind", "Name", "Subnet", "StartIp", "EndIp", "Gateway"})
	for _, r := range self {
		s := r.Spec
		table.Append([]string{r.Kind, s.Name, s.Subnet, s.StartIp, s.EndIp, s.Gateway})
	}
	table.Render()
	return tableString.String()
}

type NetworkResource struct {
	Kind string
	Spec Network
}

type Network struct {
	NetworkSpec
	Id           uint       `gorm:"not null;primaryKey;autoIncrement;"`
	DeletedAt    *time.Time `gorm:"uniqueIndex:udx_name;"`
	SpecStr      string     `gorm:"not null;column:spec" json:"-"`
	availableIps uint64     `gorm:"-" json:"-"`
	Priority     int        `gorm:"-" json:"-"`
}

type NetworkDetectSpec struct {
	Name                string
	CandidateNetworkIds []uint `json:"-"`
}

type NetworkSpec struct {
	Name    string      `gorm:"not null;uniqueIndex:udx_name;"`
	Kind    string      `gorm:"not null;"`
	Subnet  string      `gorm:"not null;"`
	StartIp string      `gorm:"not null;"`
	EndIp   string      `gorm:"not null;"`
	Gateway string      `gorm:"not null;"`
	Spec    interface{} `gorm:"-"`
}

type NetworkLocalSpec struct {
	Resolvers []Resolver `gorm:"-"`
	Nat       NetworkNat `gorm:"-"`
}

type Resolver struct {
	Resolver string
}

type NetworkNat struct {
	Enable bool
	Ports  string
}

type NetworkPort struct {
	NetworkId uint   `gorm:"not null;"`
	Ip        string `gorm:"not null;"`
	Mac       string `gorm:"not null;"`
}

const (
	KindNetworkLocal = "local"
)

func (self *VirtController) BootstrapNetwork(tctx *logger.TraceContext) (err error) {
	if err = self.sqlClient.DB.AutoMigrate(&Network{}).Error; err != nil {
		return
	}
	if err = self.sqlClient.DB.AutoMigrate(&NetworkPort{}).Error; err != nil {
		return
	}
	return
}

func (self *VirtController) CreateOrUpdateNetwork(tctx *logger.TraceContext, spec *NetworkSpec) (err error) {
	if err = self.validate.Struct(spec); err != nil {
		return
	}

	var specBytes []byte
	if specBytes, err = json.Marshal(spec.Spec); err != nil {
		return
	}

	var networkLocalSpec NetworkLocalSpec
	switch spec.Kind {
	case KindNetworkLocal:
		if err = json.Unmarshal(specBytes, &networkLocalSpec); err != nil {
			return
		}
		if err = self.validate.Struct(networkLocalSpec); err != nil {
			return
		}
	default:
		err = errors.NewBadInputErrorf("invalid network kind: kind=%s", spec.Kind)
		return
	}

	if _, err = ParseNetwork(spec); err != nil {
		return
	}

	var network *Network
	if network, err = self.GetNetwork(spec.Name); err != nil {
		if errors.IsNotFoundError(err) {
			err = self.sqlClient.Transact(tctx, func(tx *gorm.DB) (err error) {
				network := Network{
					NetworkSpec: *spec,
					SpecStr:     string(specBytes),
				}
				if err = tx.Create(&network).Error; err != nil {
					return
				}
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

func (self *VirtController) GetNetworkResources(tctx *logger.TraceContext, names []string) (networkResources NetworkResources, err error) {
	var networks []Network
	sql := self.sqlClient.DB.Table("networks").Select("*").Where("deleted_at IS NULL")
	if len(names) > 0 {
		sql = sql.Where("name in (?)", names)
	}
	if err = sql.Scan(&networks).Error; err != nil {
		return
	}

	for _, network := range networks {
		networkResources = append(networkResources, NetworkResource{
			Kind: KindNetwork,
			Spec: network,
		})
	}

	return
}

func (self *VirtController) AssignNetworkPort(tctx *logger.TraceContext, tx *gorm.DB,
	detectSpecs []NetworkDetectSpec) (sports []NetworkPort, err error) {

	var networks []Network
	sql := self.sqlClient.DB.Table("networks").Select("*").Where("deleted_at IS NULL")
	if err = sql.Scan(&networks).Error; err != nil {
		return
	}

	netIds := []uint{}
	netPortMap := map[uint]map[string]bool{}
	netMacMap := map[uint]map[string]bool{}

	candidateNetworkMap := map[uint]Network{}

	for i, spec := range detectSpecs {
		for _, network := range networks {
			if network.Name != spec.Name {
				continue
			}

			detectSpecs[i].CandidateNetworkIds = append(detectSpecs[i].CandidateNetworkIds, network.Id)
			candidateNetworkMap[network.Id] = network
		}
		if len(detectSpecs[i].CandidateNetworkIds) == 0 {
			err = errors.NewBadInputErrorf("candidate network is not found")
			return
		}
	}

	fmt.Println("DEBUG assign network port1")

	for id := range candidateNetworkMap {
		netIds = append(netIds, id)
		netPortMap[id] = map[string]bool{}
		netMacMap[id] = map[string]bool{}
	}

	// 使用済みのportをすべて取得し、利用可能なportを洗い出す
	var ports []NetworkPort
	if err = tx.Table("network_ports").
		Select("network_ports.*").
		Joins("INNER JOIN networks ON networks.id = network_ports.network_id").
		Where("networks.id IN (?)", netIds).Find(&ports).Error; err != nil {
		return
	}

	for _, port := range ports {
		netPortMap[port.NetworkId][port.Ip] = true
		netMacMap[port.NetworkId][port.Mac] = true
	}

	for id, network := range candidateNetworkMap {
		parsedStartIp := net.ParseIP(network.StartIp)
		parsedEndIp := net.ParseIP(network.EndIp)
		endIpInt := Ip2Int(parsedEndIp)
		startIpInt := Ip2Int(parsedStartIp)
		tmpAvailableIps := big.NewInt(0).Sub(endIpInt, startIpInt)
		portMap := netPortMap[network.Id]
		network.availableIps = tmpAvailableIps.Uint64() - uint64(len(portMap))
		candidateNetworkMap[id] = network
	}

	fmt.Println("DEBUG assign network port1", ports)

	for _, spec := range detectSpecs {
		fmt.Println("DEBUG spec", spec)
		candidateNetworks := []Network{}
		for _, id := range spec.CandidateNetworkIds {
			network := candidateNetworkMap[id]
			candidateNetworks = append(candidateNetworks, network)
		}

		network := candidateNetworks[0]
		portMap := netPortMap[network.Id]
		candidateIp := net.ParseIP(network.StartIp)
		for {
			if _, ok := portMap[candidateIp.String()]; ok {
				IncrementIp(candidateIp)
				continue
			}
			ipStr := candidateIp.String()
			portMap[ipStr] = true
			netPortMap[network.Id] = portMap
			fmt.Println("TODO assign", ipStr)
			break
		}
	}

	// fmt.Println("DEBUG assign network port1", availableIps)

	return
}
