package mysql

import (
	"example.com/henna-queue/internal/model"
	"example.com/henna-queue/pkg/db"
)

type ServiceRepository struct {}

func NewServiceRepository() *ServiceRepository {

	return &ServiceRepository{}
}


// GetByID 根据ID获取服务
func (r *ServiceRepository) GetByID(id uint) (*model.Service, error) {
	var service model.Service
	result := db.DB.Where("id = ?", id).First(&service)
	if result.Error != nil {
		return nil, result.Error
	}
	return &service, nil
}

// GetByShopID 获取店铺的所有服务
func (r *ServiceRepository) GetByShopID(shopID uint) ([]*model.Service, error) {	
	var services []*model.Service
	result := db.DB.Where("shop_id = ? AND status = ?", shopID, 1).
		Order("sort_order").
		Find(&services)
	if result.Error != nil {
		return nil, result.Error
	}
	return services, nil
}

// Create 创建服务
func (r *ServiceRepository) Create(service *model.Service) error {
	return db.DB.Create(service).Error
}

// Update 更新服务
func (r *ServiceRepository) Update(service *model.Service) error {
	return db.DB.Save(service).Error
}

// Delete 删除服务
func (r *ServiceRepository) Delete(id uint) error {
	return db.DB.Delete(&model.Service{}, id).Error
}

// GetAll 获取所有可用服务
func (r *ServiceRepository) GetAll() ([]*model.Service, error) {
	var services []*model.Service
	result := db.DB.Where("status = ?", 1).
		Order("sort_order").
		Find(&services)
	if result.Error != nil {
		return nil, result.Error
	}
	return services, nil
}
