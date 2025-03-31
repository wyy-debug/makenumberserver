package service

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"example.com/henna-queue/internal/model"
	"example.com/henna-queue/internal/repository"
	"example.com/henna-queue/internal/repository/factory"
	"example.com/henna-queue/internal/repository/redis"
	"example.com/henna-queue/pkg/db"
)

type AuthService struct {
	userRepo  repository.UserRepository
	adminRepo repository.AdminRepository
	cacheRepo *redis.CacheRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:  factory.NewUserRepository(),
		adminRepo: factory.NewAdminRepository(),
		cacheRepo: redis.NewCacheRepository(),
	}
}

// GetUser 获取用户
func (s *AuthService) GetUser(userID string) (*model.User, error) {
	return s.userRepo.GetByID(userID)
}

// CreateUser 创建用户
func (s *AuthService) CreateUser(user *model.User) error {
	return s.userRepo.Create(user)
}

// UpdateUser 更新用户
func (s *AuthService) UpdateUser(user *model.User) error {
	return s.userRepo.Update(user)
}

// GetAdminByUsername 根据用户名获取管理员
func (s *AuthService) GetAdminByUsername(username string) (*model.Admin, error) {
	return s.adminRepo.GetByUsername(username)
}

// GetOrCreateUser 获取或创建用户
func (s *AuthService) GetOrCreateUser(openID, unionID string) (*model.User, error) {
	// 尝试获取用户
	user, err := s.userRepo.GetByID(openID)
	if err != nil {
		// 用户不存在, 创建新用户
		user = &model.User{
			ID:        openID,
			UnionID:   unionID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = s.userRepo.Create(user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

// CacheUserToken 缓存用户token
func (s *AuthService) CacheUserToken(userID, token string) error {
	expiration := time.Hour * 24 * 7 // 7天
	return s.cacheRepo.CacheToken(userID, token, expiration)
}

// VerifyAdmin 验证管理员用户名和密码
func (s *AuthService) VerifyAdmin(username, password string) (*model.Admin, error) {
	// 获取管理员
	admin, err := s.adminRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	return admin, nil
}

// UpdateAdminLastLogin 更新管理员最后登录时间
func (s *AuthService) UpdateAdminLastLogin(adminID uint) error {
	now := time.Now()
	return db.DB.Model(&model.Admin{}).Where("id = ?", adminID).Update("last_login", now).Error
}
