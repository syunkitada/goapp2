package virt_utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/syunkitada/goapp2/pkg/lib/errors"
	"github.com/syunkitada/goapp2/pkg/lib/logger"
	"github.com/syunkitada/goapp2/pkg/lib/str_utils"
)

type ImageResources []ImageResource

func (self ImageResources) String() string {
	tableString, table := str_utils.GetTable()
	table.SetHeader([]string{"Kind", "Name"})
	for _, r := range self {
		s := r.Spec
		table.Append([]string{r.Kind, s.Name})
	}
	table.Render()
	return tableString.String()
}

type ImageResource struct {
	Kind string
	Spec Image
}

type ImageDetectSpec struct {
	Name string
}

type ImageSpec struct {
	Name      string      `gorm:"not null;uniqueIndex:udx_name;" validate:"required"`
	Namespace string      `gorm:"not null;uniqueIndex:udx_name;" validate:"required"`
	Kind      string      `gorm:"not null;" validate:"required,oneof=url"`
	Spec      interface{} `gorm:"-"`
}

type Image struct {
	ImageSpec
	Id        uint       `gorm:"not null;primaryKey;autoIncrement;"`
	DeletedAt *time.Time `gorm:"uniqueIndex:udx_name;"`
	SpecStr   string     `gorm:"not null;column:spec" json:"-"`
}

type VmImage struct {
	ImageId        uint   `gorm:"-" json:"-"`
	ImageName      string `gorm:"-" json:"-"`
	ImageNamespace string `gorm:"-" json:"-"`
	ImageKind      string `gorm:"-" json:"-"`
	ImageSpecStr   string `gorm:"-" json:"-"`
	imageUrlSpec   ImageUrlSpec
}

type ImageUrlSpec struct {
	Url        string `gorm:"not null;" validate:"required"`
	PullPolicy string `gorm:"not null;" validate:"required,oneof=IfNotPresent"`
}

const (
	KindImageUrl = "url"
)

func (self *VirtController) BootstrapImage(tctx *logger.TraceContext) (err error) {
	if err = self.sqlClient.DB.AutoMigrate(&Image{}).Error; err != nil {
		return
	}
	return
}

func (self *VirtController) CreateOrUpdateImage(tctx *logger.TraceContext, spec *ImageSpec) (err error) {
	if err = self.validate.Struct(spec); err != nil {
		return
	}

	var specBytes []byte
	if specBytes, err = json.Marshal(spec.Spec); err != nil {
		return
	}

	var imageUrlSpec ImageUrlSpec
	switch spec.Kind {
	case KindImageUrl:
		if err = json.Unmarshal(specBytes, &imageUrlSpec); err != nil {
			return
		}
		if err = self.validate.Struct(imageUrlSpec); err != nil {
			return
		}
	default:
		err = errors.NewBadInputErrorf("invalid image kind: kind=%s", spec.Kind)
		return
	}

	var image *Image
	if image, err = self.GetImage(spec.Name); err != nil {
		if errors.IsNotFoundError(err) {
			err = self.sqlClient.Transact(tctx, func(tx *gorm.DB) (err error) {
				image := Image{
					ImageSpec: *spec,
					SpecStr:   string(specBytes),
				}
				if err = tx.Create(&image).Error; err != nil {
					return
				}
				return
			})
		}
		return
	} else {
		if string(specBytes) != image.Spec {
			if err = self.sqlClient.DB.Table("images").Where("id = ?", image.Id).Updates(map[string]interface{}{
				"spec": string(specBytes),
			}).Error; err != nil {
				return
			}
		}
	}
	return
}

func (self *VirtController) GetImage(name string) (image *Image, err error) {
	var images []Image
	sql := self.sqlClient.DB.Table("images").Select("*").Where("deleted_at IS NULL")
	if err = sql.Scan(&images).Error; err != nil {
		return
	}
	if len(images) > 1 {
		err = errors.NewConflictErrorf("duplicated images are found: name=%s, len=%d", name, len(images))
		return
	} else if len(images) == 0 {
		err = errors.NewNotFoundErrorf("image is not found: name=%s", name)
		return
	}
	image = &images[0]
	return
}

func (self *VirtController) GetImageResources(tctx *logger.TraceContext, names []string) (imageResources ImageResources, err error) {
	var images []Image
	sql := self.sqlClient.DB.Table("images").Select("*").Where("deleted_at IS NULL")
	if len(names) > 0 {
		sql = sql.Where("name in (?)", names)
	}
	if err = sql.Scan(&images).Error; err != nil {
		return
	}

	for _, image := range images {
		imageResources = append(imageResources, ImageResource{
			Kind: KindImage,
			Spec: image,
		})
	}

	return
}

func (self *VirtController) DetectImage(tctx *logger.TraceContext, tx *gorm.DB,
	detectSpec *ImageDetectSpec) (image *Image, err error) {
	var images []Image
	sql := self.sqlClient.DB.Table("images").Select("*").Where("deleted_at IS NULL")
	if detectSpec.Name != "" {
		sql = sql.Where("name = ?", detectSpec.Name)
	}
	if err = sql.Scan(&images).Error; err != nil {
		return
	}

	if len(images) == 0 {
		err = errors.NewNotFoundErrorf("requested image is not found")
		return
	}

	image = &images[0]
	return
}

func (self *VirtController) PrepareImages(tctx *logger.TraceContext, vmResources VmResources) (err error) {
	// namespaceImageMap := map[string]map[string]VmImage{}
	fmt.Println("prepre images")
	for _, vm := range vmResources {
		imagePath := filepath.Join(self.imagesDir, vm.Spec.ImageName)
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			fmt.Println("download", imagePath)
		}
	}
	// TODO prepare image and set vmResources
	return
}
