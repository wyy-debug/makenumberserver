package service

import (
	"errors"
	"time"

	"henna-queue/internal/model"
	"henna-queue/internal/repository/mysql"
)

type ShopService struct {
	shopRepo    *mysql.ShopRepository
	serviceRepo *mysql.ServiceRepository
}

func NewShopService() *ShopService {
	return &ShopService{
		shopRepo:    mysql.NewShopRepository(),
		serviceRepo: mysql.NewServiceRepository(),
	}
}

// GetShop 获取店铺
func (s *ShopService) GetShop(shopID uint) (*model.Shop, error) {
	return s.shopRepo.GetByID(shopID)
}

// GetShopServices 获取店铺服务
func (s *ShopService) GetShopServices(shopID uint) ([]*model.Service, error) {
	return s.serviceRepo.GetByShopID(shopID)
}

// GetAllServices 获取店铺所有服务(包括禁用的)
func (s *ShopService) GetAllServices(shopID uint) ([]*model.Service, error) {
	return s.serviceRepo.GetAllByShopID(shopID)
}

// UpdateShop 更新店铺
func (s *ShopService) UpdateShop(shopID uint, req interface{}) (*model.Shop, error) {
	reqMap, ok := req.(*struct {
		Name          string  `json:"name"`
		Address       string  `json:"address"`
		Latitude      float64 `json:"latitude"`
		Longitude     float64 `json:"longitude"`
		Phone         string  `json:"phone"`
		BusinessHours string  `json:"business_hours"`
		Description   string  `json:"description"`
		CoverImage    string  `json:"cover_image"`
	})
	if !ok {
		return nil, errors.New("无效的请求参数")
	}
	
	shop, err := s.shopRepo.GetByID(shopID)
	if err != nil {
		return nil, errors.New("店铺不存在")
	}
	
	// 更新字段
	if reqMap.Name != "" {
		shop.Name = reqMap.Name
	}
	
	if reqMap.Address != "" {
		shop.Address = reqMap.Address
	}
	
	if reqMap.Latitude != 0 {
		shop.Latitude = reqMap.Latitude
	}
	
	if reqMap.Longitude != 0 {
		shop.Longitude = reqMap.Longitude
	}
	
	if reqMap.Phone != "" {
		shop.Phone = reqMap.Phone
	}
	
	if reqMap.BusinessHours != "" {
		shop.BusinessHours = reqMap.BusinessHours
	}
	
	if reqMap.Description != "" {
		shop.Description = reqMap.Description
	}
	
	if reqMap.CoverImage != "" {
		shop.CoverImage = reqMap.CoverImage
	}
	
	shop.UpdatedAt = time.Now()
	
	err = s.shopRepo.Update(shop)
	if err != nil {
		return nil, err
	}
	
	return shop, nil
}

// CreateService 创建服务
func (s *ShopService) CreateService(shopID uint, req interface{}) (*model.Service, error) {
	reqMap, ok := req.(*struct {
		Name        string `json:"name" binding:"required"`
		Duration    int    `json:"duration" binding:"required"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
	})
	if !ok {
		return nil, errors.New("无效的请求参数")
	}
	
	service := &model.Service{
		ShopID:      shopID,
		Name:        reqMap.Name,
		Duration:    reqMap.Duration,
		Description: reqMap.Description,
		SortOrder:   reqMap.SortOrder,
		Status:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	err := s.serviceRepo.Create(service)
	if err != nil {
		return nil, err
	}
	
	return service, nil
}

// UpdateService 更新服务
func (s *ShopService) UpdateService(serviceID, shopID uint, req interface{}) (*model.Service, error) {
	reqMap, ok := req.(*struct {
		Name        string `json:"name"`
		Duration    int    `json:"duration"`
		Description string `json:"description"`
		Status      *int8  `json:"status"`
		SortOrder   int    `json:"sort_order"`
	})
	if !ok {
		return nil, errors.New("无效的请求参数")
	}
	
	service, err := s.serviceRepo.GetByID(serviceID)
	if err != nil {
		return nil, errors.New("服务不存在")
	}
	
	if service.ShopID != shopID {
		return nil, errors.New("无权操作该服务")
	}
	
	// 更新字段
	if reqMap.Name != "" {
		service.Name = reqMap.Name
	}
	
	if reqMap.Duration != 0 {
		service.Duration = reqMap.Duration
	}
	
	if reqMap.Description != "" {
		service.Description = reqMap.Description
	}
	
	if reqMap.Status != nil {
		service.Status = *reqMap.Status
	}
	
	if reqMap.SortOrder != 0 {
		service.SortOrder = reqMap.SortOrder
	}
	
	service.UpdatedAt = time.Now()
	
	err = s.serviceRepo.Update(service)
	if err != nil {
		return nil, err
	}
	
	return service, nil
}

// DeleteService 删除服务
func (s *ShopService) DeleteService(serviceID, shopID uint) error {
	service, err := s.serviceRepo.GetByID(serviceID)
	if err != nil {
		return errors.New("服务不存在")
	}
	
	if service.ShopID != shopID {
		return errors.New("无权操作该服务")
	}
	
	return s.serviceRepo.Delete(serviceID)
} 