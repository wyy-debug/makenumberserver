package mysql

import (
	"errors"
	"example.com/henna-queue/internal/model"
	"example.com/henna-queue/pkg/db"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository() *UserRepository {
	// 不在这里 panic，允许延迟初始化
	var dbConn *gorm.DB
	if db.DB != nil {
		dbConn = db.DB
	}
	
	return &UserRepository{
		DB: dbConn,
	}
}

// 获取安全的数据库连接
func (r *UserRepository) getDBConnection() (*gorm.DB, error) {
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

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id string) (*model.User, error) {
	dbConn, err := r.getDBConnection()
	if err != nil {
		return nil, err
	}
	
	var user model.User
	result := dbConn.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetByPhone 根据手机号获取用户
func (r *UserRepository) GetByPhone(phone string) (*model.User, error) {
	dbConn, err := r.getDBConnection()
	if err != nil {
		return nil, err
	}
	
	var user model.User
	result := dbConn.Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	dbConn, err := r.getDBConnection()
	if err != nil {
		return err
	}
	
	return dbConn.Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	dbConn, err := r.getDBConnection()
	if err != nil {
		return err
	}
	
	return dbConn.Save(user).Error
}

// UpdatePhone 更新用户手机号
func (r *UserRepository) UpdatePhone(userID string, phone string) error {
	dbConn, err := r.getDBConnection()
	if err != nil {
		return err
	}
	
	return dbConn.Model(&model.User{}).Where("id = ?", userID).Update("phone", phone).Error
}
