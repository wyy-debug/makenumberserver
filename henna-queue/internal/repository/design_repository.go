// Add this method to the existing DesignRepository

// UpdateCategoryByOldCategory 更新指定店铺中特定分类下的所有图案到新分类
func (r *DesignRepository) UpdateCategoryByOldCategory(shopID uint, oldCategory, newCategory string) error {
	return r.db.Model(&model.Design{}).
		Where("shop_id = ? AND category = ?", shopID, oldCategory).
		Update("category", newCategory).Error
} 