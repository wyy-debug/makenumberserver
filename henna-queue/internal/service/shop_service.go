package service

import (
	"errors"
	"example.com/henna-queue/internal/model"
	"example.com/henna-queue/internal/repository"
	"example.com/henna-queue/internal/repository/factory"
	"time"
)

type ShopService struct {
	shopRepo    repository.ShopRepository
	serviceRepo repository.ServiceRepository
}

func NewShopService() *ShopService {
	return &ShopService{
		shopRepo:    factory.NewShopRepository(),
		serviceRepo: factory.NewServiceRepository(),
	}
}

// GetShop 获取店铺信息
func (s *ShopService) GetShop(id uint) (*model.Shop, error) {
	return s.shopRepo.GetByID(id)
}

// GetShopServices 获取店铺服务列表
func (s *ShopService) GetShopServices(shopID uint) ([]*model.Service, error) {
	return s.serviceRepo.GetByShopID(shopID)
}

// GetAllServices 获取店铺所有服务(包括禁用的)
func (s *ShopService) GetAllServices(shopID uint) ([]*model.Service, error) {
	return s.serviceRepo.GetByShopID(shopID)
}

// CreateShop 创建店铺
func (s *ShopService) CreateShop(shop *model.Shop) error {
	return s.shopRepo.Create(shop)
}

// UpdateShop 更新店铺信息
func (s *ShopService) UpdateShop(shopID uint, req *struct {
	Name          string  `json:"name"`
	Address       string  `json:"address"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Phone         string  `json:"phone"`
	BusinessHours string  `json:"business_hours"`
	Description   string  `json:"description"`
	CoverImage    string  `json:"cover_image"`
}) (*model.Shop, error) {
	if req == nil {
		return nil, errors.New("无效的请求参数")
	}

	shop, err := s.shopRepo.GetByID(shopID)
	if err != nil {
		return nil, errors.New("店铺不存在")
	}

	// 更新字段
	if req.Name != "" {
		shop.Name = req.Name
	}

	if req.Address != "" {
		shop.Address = req.Address
	}

	if req.Latitude != 0 {
		shop.Latitude = req.Latitude
	}

	if req.Longitude != 0 {
		shop.Longitude = req.Longitude
	}

	if req.Phone != "" {
		shop.Phone = req.Phone
	}

	if req.BusinessHours != "" {
		shop.BusinessHours = req.BusinessHours
	}

	if req.Description != "" {
		shop.Description = req.Description
	}

	if req.CoverImage != "" {
		shop.CoverImage = req.CoverImage
	}

	shop.UpdatedAt = time.Now()

	err = s.shopRepo.Update(shop)
	if err != nil {
		return nil, err
	}

	return shop, nil
}

// DeleteShop 删除店铺
func (s *ShopService) DeleteShop(id uint) error {
	return s.shopRepo.Delete(id)
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

// UpdateService 更新服务信息
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
