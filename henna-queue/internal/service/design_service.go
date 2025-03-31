package service

import (
	"errors"
	"time"

	"henna-queue/internal/model"
	"henna-queue/internal/repository/mysql"
	"henna-queue/pkg/db"
)

type DesignService struct {
	designRepo   *mysql.DesignRepository
	favoriteRepo *mysql.FavoriteRepository
}

func NewDesignService() *DesignService {
	return &DesignService{
		designRepo:   mysql.NewDesignRepository(),
		favoriteRepo: mysql.NewFavoriteRepository(),
	}
}

// GetDesigns 获取图案列表
func (s *DesignService) GetDesigns(shopID uint, category, userID string, page, pageSize int) ([]*model.DesignResponse, int, error) {
	designs, total, err := s.designRepo.GetDesigns(shopID, category, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	
	// 构造响应, 并检查是否已收藏
	var result []*model.DesignResponse
	for _, design := range designs {
		isLiked := false
		if userID != "" {
			// 检查是否已收藏
			if liked, _ := s.favoriteRepo.CheckFavorite(userID, design.ID); liked {
				isLiked = true
			}
		}
		
		result = append(result, &model.DesignResponse{
			TattooDesign: *design,
			IsLiked:      isLiked,
		})
	}
	
	return result, total, nil
}

// GetDesign 获取单个图案
func (s *DesignService) GetDesign(designID uint, userID string) (*model.DesignResponse, error) {
	design, err := s.designRepo.GetByID(designID)
	if err != nil {
		return nil, err
	}
	
	isLiked := false
	if userID != "" {
		// 检查是否已收藏
		if liked, _ := s.favoriteRepo.CheckFavorite(userID, design.ID); liked {
			isLiked = true
		}
	}
	
	return &model.DesignResponse{
		TattooDesign: *design,
		IsLiked:      isLiked,
	}, nil
}

// ToggleFavorite 收藏/取消收藏图案
func (s *DesignService) ToggleFavorite(userID string, designID uint) (bool, error) {
	// 检查图案是否存在
	_, err := s.designRepo.GetByID(designID)
	if err != nil {
		return false, errors.New("图案不存在")
	}
	
	// 检查是否已收藏
	isFavorite, err := s.favoriteRepo.CheckFavorite(userID, designID)
	if err != nil && err.Error() != "record not found" {
		return false, err
	}
	
	if isFavorite {
		// 已收藏，取消收藏
		err = s.favoriteRepo.DeleteFavorite(userID, designID)
		if err != nil {
			return true, err
		}
		
		// 减少收藏数
		err = db.DB.Model(&model.TattooDesign{}).Where("id = ?", designID).
			UpdateColumn("likes", db.DB.Raw("likes - 1")).Error
		if err != nil {
			return false, err
		}
		
		return false, nil
	} else {
		// 未收藏，添加收藏
		favorite := &model.Favorite{
			UserID:    userID,
			DesignID:  designID,
			CreatedAt: time.Now(),
		}
		
		err = s.favoriteRepo.Create(favorite)
		if err != nil {
			return false, err
		}
		
		// 增加收藏数
		err = db.DB.Model(&model.TattooDesign{}).Where("id = ?", designID).
			UpdateColumn("likes", db.DB.Raw("likes + 1")).Error
		if err != nil {
			return true, err
		}
		
		return true, nil
	}
}

// GetUserFavorites 获取用户收藏
func (s *DesignService) GetUserFavorites(userID string, page, pageSize int) ([]*model.TattooDesign, int, error) {
	return s.favoriteRepo.GetUserFavorites(userID, page, pageSize)
}

// GetAdminDesigns 获取管理后台图案列表
func (s *DesignService) GetAdminDesigns(shopID uint, category, status string, page, pageSize int) ([]*model.TattooDesign, int, error) {
	return s.designRepo.GetAdminDesigns(shopID, category, status, page, pageSize)
}

// CreateDesign 创建图案
func (s *DesignService) CreateDesign(shopID uint, req interface{}) (*model.TattooDesign, error) {
	reqMap, ok := req.(*struct {
		Title       string `json:"title" binding:"required"`
		Category    string `json:"category" binding:"required"`
		ImageURL    string `json:"image_url" binding:"required"`
		Description string `json:"description"`
	})
	if !ok {
		return nil, errors.New("无效的请求参数")
	}
	
	design := &model.TattooDesign{
		ShopID:      shopID,
		Title:       reqMap.Title,
		Category:    reqMap.Category,
		ImageURL:    reqMap.ImageURL,
		Description: reqMap.Description,
		Status:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	err := s.designRepo.Create(design)
	if err != nil {
		return nil, err
	}
	
	return design, nil
}

// UpdateDesign 更新图案
func (s *DesignService) UpdateDesign(designID, shopID uint, req interface{}) (*model.TattooDesign, error) {
	reqMap, ok := req.(*struct {
		Title       string `json:"title"`
		Category    string `json:"category"`
		ImageURL    string `json:"image_url"`
		Description string `json:"description"`
		Status      *int8  `json:"status"`
	})
	if !ok {
		return nil, errors.New("无效的请求参数")
	}
	
	design, err := s.designRepo.GetByID(designID)
	if err != nil {
		return nil, errors.New("图案不存在")
	}
	
	if design.ShopID != shopID {
		return nil, errors.New("无权操作该图案")
	}
	
	// 更新字段
	if reqMap.Title != "" {
		design.Title = reqMap.Title
	}
	
	if reqMap.Category != "" {
		design.Category = reqMap.Category
	}
	
	if reqMap.ImageURL != "" {
		design.ImageURL = reqMap.ImageURL
	}
	
	if reqMap.Description != "" {
		design.Description = reqMap.Description
	}
	
	if reqMap.Status != nil {
		design.Status = *reqMap.Status
	}
	
	design.UpdatedAt = time.Now()
	
	err = s.designRepo.Update(design)
	if err != nil {
		return nil, err
	}
	
	return design, nil
}

// DeleteDesign 删除图案
func (s *DesignService) DeleteDesign(designID, shopID uint) error {
	design, err := s.designRepo.GetByID(designID)
	if err != nil {
		return errors.New("图案不存在")
	}
	
	if design.ShopID != shopID {
		return errors.New("无权操作该图案")
	}
	
	return s.designRepo.Delete(designID)
} 