package mysql

import (
	"henna-queue/internal/model"
	"henna-queue/pkg/db"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	result := db.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return db.DB.Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	return db.DB.Save(user).Error
}

// UpdatePhone 更新用户手机号
func (r *UserRepository) UpdatePhone(userID string, phone string) error {
	return db.DB.Model(&model.User{}).Where("id = ?", userID).Update("phone", phone).Error
} 