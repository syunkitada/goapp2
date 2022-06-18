package virt_utils

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/syunkitada/goapp2/pkg/lib/errors"
	"github.com/syunkitada/goapp2/pkg/lib/logger"
)

type ImageResource struct {
	Kind string
	Spec Image
}

type ImageSpec struct {
	Name string      `gorm:"not null;uniqueIndex:udx_name;" validate:"required"`
	Kind string      `gorm:"not null;" validate:"required,oneof=url"`
	Spec interface{} `gorm:"-"`
}

type Image struct {
	ImageSpec
	Id        uint       `gorm:"not null;primaryKey;autoIncrement;"`
	DeletedAt *time.Time `gorm:"uniqueIndex:udx_name;"`
}

type ImageUrlSpec struct {
	Url        string `gorm:"not null;" validate:"required"`
	PullPolicy string `gorm:"not null;" validate:"required,oneof=IfNotPresent"`
}

type ImageUrl struct {
	ImageUrlSpec
	ImageId uint `gorm:"not null;primaryKey;"`
}

const (
	KindImageUrl = "url"
)

func (self *VirtController) BootstrapImage(tctx *logger.TraceContext) (err error) {
	if err = self.sqlClient.DB.AutoMigrate(&Image{}).Error; err != nil {
		return
	}
	if err = self.sqlClient.DB.AutoMigrate(&ImageUrl{}).Error; err != nil {
		return
	}
	return
}

func (self *VirtController) CreateOrUpdateImage(tctx *logger.TraceContext, spec *ImageSpec) (err error) {
	if err = self.validate.Struct(spec); err != nil {
		return
	}

	var bytes []byte
	if bytes, err = json.Marshal(spec.Spec); err != nil {
		return
	}

	var imageUrlSpec ImageUrlSpec
	switch spec.Kind {
	case KindImageUrl:
		if err = json.Unmarshal(bytes, &imageUrlSpec); err != nil {
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
				}
				if err = self.sqlClient.DB.Create(&image).Error; err != nil {
					return
				}
				switch spec.Kind {
				case KindImageUrl:
					imageUrl := ImageUrl{
						ImageUrlSpec: imageUrlSpec,
						ImageId:      image.Id,
					}
					if err = self.sqlClient.DB.Create(&imageUrl).Error; err != nil {
						return
					}
				}
				return
			})
		}
		return
	} else {
		switch spec.Kind {
		case KindImageUrl:
			if err = self.sqlClient.DB.Table("image_urls").Where("image_id = ?", image.Id).Updates(map[string]interface{}{
				"url": imageUrlSpec.Url,
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

func (self *VirtController) GetImageResources(tctx *logger.TraceContext, names []string) (imageResources []ImageResource, err error) {
	var images []Image
	sql := self.sqlClient.DB.Table("images").Select("*").Where("deleted_at IS NULL")
	if err = sql.Scan(&images).Error; err != nil {
		return
	}

	var imageUrls []ImageUrl
	sql = self.sqlClient.DB.Table("image_urls").Select("*")
	if err = sql.Scan(&imageUrls).Error; err != nil {
		return
	}

	imageUrlMap := map[uint]ImageUrl{}
	for _, imageUrl := range imageUrls {
		imageUrlMap[imageUrl.ImageId] = imageUrl
	}

	for _, image := range images {
		switch image.Kind {
		case KindImageUrl:
			image.Spec = imageUrlMap[image.Id]
		}
		imageResources = append(imageResources, ImageResource{
			Kind: KindImage,
			Spec: image,
		})
	}

	return
}
