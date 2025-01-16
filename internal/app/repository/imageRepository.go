package repository

import "gorm.io/gorm"

type Image struct {
	Sol  int    `gorm:"column:sol;primaryKey;autoIncrement:false"`
	Url  string `gorm:"column:url;type:text"`
	Size int64  `gorm:"column:image_size"`
}

type Handler struct {
	DB *gorm.DB
}

func NewImageRepository(db *gorm.DB) *Handler {
	return &Handler{
		DB: db,
	}
}

func (h *Handler) GetImageById(sol int) (*Image, error) {
	var image Image
	if err := h.DB.Where("sol = ?", sol).First(&image).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

func (h *Handler) SaveImage(image *Image) error {
	return h.DB.Create(image).Error
}
