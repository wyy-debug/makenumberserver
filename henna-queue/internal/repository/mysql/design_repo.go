package mysql

import (
	"henna-queue/internal/model"
	"henna-queue/pkg/db"
)

type DesignRepository struct{}

func NewDesignRepository() *DesignRepository {
	return &DesignRepository{}
}

// GetByID 根据ID获取图案
func (r *DesignRepository) GetByID(id uint) (*model.TattooDesign, error) {
	var design model.TattooDesign
	result := db.DB.Where("id = ?", id).First(&design)
	if result.Error != nil {
		return nil, result.Error
	}
	return &design, nil
}

// GetDesigns 获取图案列表
func (r *DesignRepository) GetDesigns(shopID uint, category string, page, pageSize int) ([]*model.TattooDesign, int, error) {
	offset := (page - 1) * pageSize
	var designs []*model.TattooDesign
	var total int64
	
	query := db.DB.Model(&model.TattooDesign{}).Where("shop_id = ? AND status = ?", shopID, 1)
	
	if category != "" {
		query = query.Where("category = ?", category)
	}
	
	// 获取总数
	query.Count(&total)
	
	// 获取分页数据
	result := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&designs)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	
	return designs, int(total), nil
}

// GetAdminDesigns 获取管理后台图案列表
func (r *DesignRepository) GetAdminDesigns(shopID uint, category, status string, page, pageSize int) ([]*model.TattooDesign, int, error) {
	offset := (page - 1) * pageSize
	var designs []*model.TattooDesign
	var total int64
	
	query := db.DB.Model(&model.TattooDesign{}).Where("shop_id = ?", shopID)
	
	if category != "" {
		query = query.Where("category = ?", category)
	}
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	// 获取总数
	query.Count(&total)
	
	// 获取分页数据
	result := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&designs)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	
	return designs, int(total), nil
}

// Create 创建图案
func (r *DesignRepository) Create(design *model.TattooDesign) error {
	return db.DB.Create(design).Error
}

// Update 更新图案
func (r *DesignRepository) Update(design *model.TattooDesign) error {
	return db.DB.Save(design).Error
}

// Delete 删除图案
func (r *DesignRepository) Delete(id uint) error {
	return db.DB.Delete(&model.TattooDesign{}, id).Error
} 