package mysql

import (
	"errors"
	"example.com/henna-queue/internal/model"
	"example.com/henna-queue/pkg/db"
	"gorm.io/gorm"
)

type ServiceRepository struct {
	DB *gorm.DB
}

func NewServiceRepository() *ServiceRepository {
	// 不在这里 panic，允许延迟初始化
	var dbConn *gorm.DB
	if db.DB != nil {
		dbConn = db.DB
	}
	
	return &ServiceRepository{
		DB: dbConn,
	}
}

// 获取安全的数据库连接
func (r *ServiceRepository) getDBConnection() (*gorm.DB, error) {
	if r.DB != nil {
		return r.DB, nil
	}
	
	// 尝试使用全局 DB
	if db.DB == nil {
		return nil, errors.New("数据库连接未初始化")
	}
	
	// 更新并返回连接
	r.DB = db.DB
	return r.DB, nil
}

// GetByID 根据ID获取服务
func (r *ServiceRepository) GetByID(id uint) (*model.Service, error) {
	dbConn, err := r.getDBConnection()
	if err != nil {
		return nil, err
	}
	
	var service model.Service
	result := dbConn.Where("id = ?", id).First(&service)
	if result.Error != nil {
		return nil, result.Error
	}
	return &service, nil
}

// GetByShopID 获取店铺的所有服务
func (r *ServiceRepository) GetByShopID(shopID uint) ([]*model.Service, error) {
	dbConn, err := r.getDBConnection()
	if err != nil {
		return nil, err
	}
	
	var services []*model.Service
	result := dbConn.Where("shop_id = ? AND status = ?", shopID, 1).
		Order("sort_order").
		Find(&services)
	if result.Error != nil {
		return nil, result.Error
	}
	return services, nil
}

// Create 创建服务
func (r *ServiceRepository) Create(service *model.Service) error {
	dbConn, err := r.getDBConnection()
	if err != nil {
		return err
	}
	
	return dbConn.Create(service).Error
}

// Update 更新服务
func (r *ServiceRepository) Update(service *model.Service) error {
	dbConn, err := r.getDBConnection()
	if err != nil {
		return err
	}
	
	return dbConn.Save(service).Error
}

// Delete 删除服务
func (r *ServiceRepository) Delete(id uint) error {
	dbConn, err := r.getDBConnection()
	if err != nil {
		return err
	}
	
	return dbConn.Delete(&model.Service{}, id).Error
}

// GetAll 获取所有可用服务
func (r *ServiceRepository) GetAll() ([]*model.Service, error) {
	dbConn, err := r.getDBConnection()
	if err != nil {
		return nil, err
	}
	
	var services []*model.Service
	result := dbConn.Where("status = ?", 1).
		Order("sort_order").
		Find(&services)
	if result.Error != nil {
		return nil, result.Error
	}
	return services, nil
}
