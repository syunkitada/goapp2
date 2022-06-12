package virt_utils

import "time"

type Resource struct {
	Kind string `validate:"required"`
}

type Vm struct {
	Kind string
	Spec VmSpec
}

type ImageSpec struct {
}

type VmSpec struct {
	Name   string
	Vcpus  uint `validate:"required"`
	Memory uint `validate:"required"`
	Disk   uint `validate:"required"`

	Image ImageSpec
	Ports []NewworkPort
}

type Network struct {
	Name      string `gorm:"primaryKey"`
	DeletedAt *time.Time
}

type NewworkPort struct {
}
