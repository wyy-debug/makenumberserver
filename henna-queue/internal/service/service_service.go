package service

import (
	"errors"
	"time"

	"example.com/henna-queue/internal/model"
	"example.com/henna-queue/internal/repository"
	"example.com/henna-queue/internal/repository/factory"
)

type ServiceService struct {
	serviceRepo repository.ServiceRepository
	shopRepo    repository.ShopRepository
}

func NewServiceService() *ServiceService {
	return &ServiceService{
		serviceRepo: factory.NewServiceRepository(),
		shopRepo:    factory.NewShopRepository(),
	}
}

// GetService 获取单个服务
func (s *ServiceService) GetService(serviceID uint) (*model.Service, error) {
	return s.serviceRepo.GetByID(serviceID)
}

// GetShopServices 获取店铺的所有服务
func (s *ServiceService) GetShopServices(shopID uint) ([]*model.Service, error) {
	return s.serviceRepo.GetByShopID(shopID)
}

// CreateService 创建服务
func (s *ServiceService) CreateService(shopID uint, name string, duration int, description string, status int8, sortOrder int) (*model.Service, error) {
	// 验证店铺是否存在
	_, err := s.shopRepo.GetByID(shopID)
	if err != nil {
		return nil, errors.New("店铺不存在")
	}

	// 创建服务
	now := time.Now()
	service := &model.Service{
		ShopID:      shopID,
		Name:        name,
		Duration:    duration,
		Description: description,
		Status:      status,
		SortOrder:   sortOrder,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = s.serviceRepo.Create(service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

// UpdateService 更新服务
func (s *ServiceService) UpdateService(serviceID uint, shopID uint, name string, duration int, description string, status int8, sortOrder int) (*model.Service, error) {
	// 获取服务
	service, err := s.serviceRepo.GetByID(serviceID)
	if err != nil {
		return nil, err
	}

	// 验证服务是否属于该店铺
	if service.ShopID != shopID {
		return nil, errors.New("无权操作该服务")
	}

	// 更新服务
	service.Name = name
	service.Duration = duration
	service.Description = description
	service.Status = status
	service.SortOrder = sortOrder
	service.UpdatedAt = time.Now()

	err = s.serviceRepo.Update(service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

// DeleteService 删除服务
func (s *ServiceService) DeleteService(serviceID uint, shopID uint) error {
	// 获取服务
	service, err := s.serviceRepo.GetByID(serviceID)
	if err != nil {
		return err
	}

	// 验证服务是否属于该店铺
	if service.ShopID != shopID {
		return errors.New("无权操作该服务")
	}

	// 删除服务
	return s.serviceRepo.Delete(serviceID)
}

// GetAllServices 获取所有可用服务
func (s *ServiceService) GetAllServices() ([]*model.Service, error) {
	return s.serviceRepo.GetAll()
} 