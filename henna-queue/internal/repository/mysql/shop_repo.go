package mysql

import (
	"example.com/henna-queue/internal/model"
	"example.com/henna-queue/pkg/db"
)

type ShopRepository struct{}

func NewShopRepository() *ShopRepository {
	return &ShopRepository{}
}

// GetByID 根据ID获取店铺
func (r *ShopRepository) GetByID(id uint) (*model.Shop, error) {
	var shop model.Shop
	result := db.DB.Where("id = ?", id).First(&shop)
	if result.Error != nil {
		return nil, result.Error
	}
	return &shop, nil
}

// GetAll 获取所有店铺
func (r *ShopRepository) GetAll() ([]*model.Shop, error) {
	var shops []*model.Shop
	result := db.DB.Find(&shops)
	if result.Error != nil {
		return nil, result.Error
	}
	return shops, nil
}

// Create 创建店铺
func (r *ShopRepository) Create(shop *model.Shop) error {
	return db.DB.Create(shop).Error
}

// Update 更新店铺
func (r *ShopRepository) Update(shop *model.Shop) error {
	return db.DB.Save(shop).Error
}

// Delete 删除店铺
func (r *ShopRepository) Delete(id uint) error {
	return db.DB.Delete(&model.Shop{}, id).Error
}
