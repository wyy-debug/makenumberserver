package repository

import (
	"example.com/henna-queue/internal/model"
	"example.com/henna-queue/pkg/db"
	"gorm.io/gorm"
)

type DesignRepository struct {
	db *gorm.DB
}

func NewDesignRepository() *DesignRepository {
	return &DesignRepository{
		db: db.DB,
	}
}

// UpdateCategoryByOldCategory 更新指定店铺中特定分类下的所有图案到新分类
func (r *DesignRepository) UpdateCategoryByOldCategory(shopID uint, oldCategory, newCategory string) error {
	return r.db.Model(&model.TattooDesign{}).
		Where("shop_id = ? AND category = ?", shopID, oldCategory).
		Update("category", newCategory).Error
}
