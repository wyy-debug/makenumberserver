package mysql

import (
	"example.com/henna-queue/internal/model"
	"example.com/henna-queue/pkg/db"
)

type FavoriteRepository struct{}

func NewFavoriteRepository() *FavoriteRepository {
	return &FavoriteRepository{}
}

// CheckFavorite 检查是否已收藏
func (r *FavoriteRepository) CheckFavorite(userID string, designID uint) (bool, error) {
	var favorite model.Favorite
	result := db.DB.Where("user_id = ? AND design_id = ?", userID, designID).First(&favorite)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

// Create 创建收藏
func (r *FavoriteRepository) Create(favorite *model.Favorite) error {
	return db.DB.Create(favorite).Error
}

// DeleteFavorite 删除收藏
func (r *FavoriteRepository) DeleteFavorite(userID string, designID uint) error {
	return db.DB.Where("user_id = ? AND design_id = ?", userID, designID).Delete(&model.Favorite{}).Error
}

// GetUserFavorites 获取用户收藏
func (r *FavoriteRepository) GetUserFavorites(userID string, page, pageSize int) ([]*model.TattooDesign, int, error) {
	offset := (page - 1) * pageSize
	var designs []*model.TattooDesign
	var total int64

	// 获取总数
	db.DB.Model(&model.Favorite{}).Where("user_id = ?", userID).Count(&total)

	// 获取分页数据
	result := db.DB.Model(&model.Favorite{}).
		Select("tattoo_designs.*").
		Joins("JOIN tattoo_designs ON favorites.design_id = tattoo_designs.id").
		Where("favorites.user_id = ? AND tattoo_designs.status = ?", userID, 1).
		Order("favorites.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&designs)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return designs, int(total), nil
}
