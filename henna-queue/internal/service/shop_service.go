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

// GetShopStats 获取店铺统计数据
func (s *ShopService) GetShopStats(shopID uint, startDate, endDate string) (map[string]interface{}, error) {
	// 验证店铺ID
	if _, err := s.GetShop(shopID); err != nil {
		return nil, err
	}

	// 这里应该从数据库获取真实的统计数据
	// 这是一个示例实现
	return map[string]interface{}{
		"total_customers": 120,
		"average_wait_time": 25, // 分钟
		"peak_hours": []string{"14:00", "16:00"},
		"services_breakdown": []map[string]interface{}{
			{
				"name": "普通海娜",
				"count": 45,
				"percentage": 37.5,
			},
			{
				"name": "定制海娜",
				"count": 35,
				"percentage": 29.2,
			},
			{
				"name": "手臂海娜",
				"count": 25,
				"percentage": 20.8,
			},
			{
				"name": "脚部海娜",
				"count": 15,
				"percentage": 12.5,
			},
		},
		"daily_customers": []map[string]interface{}{
			{"date": "2023-09-25", "count": 18},
			{"date": "2023-09-26", "count": 22},
			{"date": "2023-09-27", "count": 15},
			{"date": "2023-09-28", "count": 20},
			{"date": "2023-09-29", "count": 25},
			{"date": "2023-09-30", "count": 12},
			{"date": "2023-10-01", "count": 8},
		},
	}, nil
}

// GetPublicServices 获取公开的服务列表
func (s *ShopService) GetPublicServices(shopID uint) ([]map[string]interface{}, error) {
	// 如果提供了店铺ID，则获取特定店铺的服务
	if shopID > 0 {
		// 验证店铺ID
		if _, err := s.GetShop(shopID); err != nil {
			return nil, err
		}

		// 从数据库获取该店铺下状态为启用的服务
		// 这里应该有真实的数据库查询
		// 这是一个示例实现
		return []map[string]interface{}{
			{
				"id": 1,
				"name": "普通海娜",
				"duration": 30,
				"description": "简单的手部海娜纹身",
				"shop_id": shopID,
			},
			{
				"id": 2,
				"name": "定制海娜",
				"duration": 60,
				"description": "根据客户需求定制的海娜设计",
				"shop_id": shopID,
			},
		}, nil
	}

	// 否则返回所有公开的服务
	return []map[string]interface{}{
		{
			"id": 1,
			"name": "普通海娜",
			"duration": 30,
			"description": "简单的手部海娜纹身",
			"shop_id": 1,
		},
		{
			"id": 2,
			"name": "定制海娜",
			"duration": 60,
			"description": "根据客户需求定制的海娜设计",
			"shop_id": 1,
		},
		{
			"id": 3,
			"name": "脚部海娜",
			"duration": 45,
			"description": "适用于脚部的海娜设计",
			"shop_id": 2,
		},
	}, nil
}
