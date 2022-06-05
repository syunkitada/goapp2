package virt_utils

type Resource struct {
	Kind string      `validate:"required"`
	Spec interface{} `validate:"required"`
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

	Image Image `validate:"required"`
	Ports []Port
}

type Network struct {
}
