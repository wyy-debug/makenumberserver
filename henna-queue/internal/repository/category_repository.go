package repository

import (
	"errors"

	"gorm.io/gorm"

	"example.com/henna-queue/internal/model"
	"example.com/henna-queue/pkg/db"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{
		db: db.DB,
	}
}

// FindByID 根据ID查找分类
func (r *CategoryRepository) FindByID(id uint) (*model.Category, error) {
	var category model.Category
	result := r.db.First(&category, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, result.Error
	}
	return &category, nil
}

// FindByShopID 根据店铺ID查找分类列表
func (r *CategoryRepository) FindByShopID(shopID uint) ([]model.Category, error) {
	var categories []model.Category
	result := r.db.Where("shop_id = ?", shopID).Order("sort_order asc, id asc").Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}
	return categories, nil
}

// Create 创建分类
func (r *CategoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

// Update 更新分类
func (r *CategoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

// Delete 删除分类
func (r *CategoryRepository) Delete(category *model.Category) error {
	return r.db.Delete(category).Error
}

// ExistsByCode 检查分类代码是否存在
func (r *CategoryRepository) ExistsByCode(shopID uint, code string) (bool, error) {
	var count int64
	result := r.db.Model(&model.Category{}).Where("shop_id = ? AND code = ?", shopID, code).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
