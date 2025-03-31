package service

import (
	"errors"
	"regexp"

	"example.com/henna-queue/internal/model"
	"example.com/henna-queue/internal/repository"
)

type CategoryService struct {
	repo       *repository.CategoryRepository
	designRepo *repository.DesignRepository
}

func NewCategoryService() *CategoryService {
	return &CategoryService{
		repo:       repository.NewCategoryRepository(),
		designRepo: repository.NewDesignRepository(),
	}
}

// GetCategories 获取分类列表
func (s *CategoryService) GetCategories(shopID uint) ([]model.Category, error) {
	return s.repo.FindByShopID(shopID)
}

// CreateCategory 创建分类
func (s *CategoryService) CreateCategory(shopID uint, name, code string, sortOrder int) (*model.Category, error) {
	// 验证code格式
	matched, _ := regexp.MatchString("^[a-zA-Z_]+$", code)
	if !matched {
		return nil, errors.New("分类代码只能包含字母和下划线")
	}

	// 检查code是否已存在
	exists, err := s.repo.ExistsByCode(shopID, code)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("分类代码已存在")
	}

	// 创建分类
	category := &model.Category{
		ShopID:    shopID,
		Name:      name,
		Code:      code,
		SortOrder: sortOrder,
	}

	err = s.repo.Create(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory 更新分类
func (s *CategoryService) UpdateCategory(categoryID, shopID uint, name string, sortOrder int) (*model.Category, error) {
	// 获取分类
	category, err := s.repo.FindByID(categoryID)
	if err != nil {
		return nil, err
	}

	// 验证归属
	if category.ShopID != shopID {
		return nil, errors.New("您没有权限修改此分类")
	}

	// 更新分类
	category.Name = name
	category.SortOrder = sortOrder

	err = s.repo.Update(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory 删除分类
func (s *CategoryService) DeleteCategory(categoryID, shopID uint) error {
	// 获取分类
	category, err := s.repo.FindByID(categoryID)
	if err != nil {
		return err
	}

	// 验证归属
	if category.ShopID != shopID {
		return errors.New("您没有权限删除此分类")
	}

	// 将此分类下的图案重新分类为"other"
	err = s.designRepo.UpdateCategoryByOldCategory(shopID, category.Code, "other")
	if err != nil {
		return err
	}

	// 删除分类
	return s.repo.Delete(category)
}
