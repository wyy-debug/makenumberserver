package service

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"example.com/henna-queue/internal/model"
	"example.com/henna-queue/internal/repository"
	"example.com/henna-queue/internal/repository/factory"
)

type AdminService struct {
	adminRepo     repository.AdminRepository
	statisticRepo repository.StatisticRepository
}

func NewAdminService() *AdminService {
	return &AdminService{
		adminRepo:     factory.NewAdminRepository(),
		statisticRepo: factory.NewStatisticRepository(),
	}
}

// GetAdmin 获取管理员
func (s *AdminService) GetAdmin(id uint) (*model.Admin, error) {
	return s.adminRepo.GetByID(id)
}

// GetAdminByUsername 根据用户名获取管理员
func (s *AdminService) GetAdminByUsername(username string) (*model.Admin, error) {
	return s.adminRepo.GetByUsername(username)
}

// CheckUsernameExists 检查用户名是否存在
func (s *AdminService) CheckUsernameExists(username string) (bool, error) {
	return s.adminRepo.CheckUsernameExists(username)
}

// UpdateAdmin 更新管理员
func (s *AdminService) UpdateAdmin(id uint, req interface{}) (*model.Admin, error) {
	admin, err := s.adminRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 根据请求参数更新管理员信息
	// 这里需要根据具体的请求结构进行转换
	// 目前根据错误信息推断，请求结构包含 Username, OldPassword, NewPassword 字段
	reqMap, ok := req.(*struct {
		Username    string `json:"username"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	})
	if !ok {
		return nil, errors.New("无效的请求参数")
	}

	// 更新用户名
	if reqMap.Username != "" {
		admin.Username = reqMap.Username
	}

	// 如果提供了旧密码和新密码，则更新密码
	if reqMap.OldPassword != "" && reqMap.NewPassword != "" {
		// 验证旧密码
		err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(reqMap.OldPassword))
		if err != nil {
			return nil, errors.New("旧密码不正确")
		}

		// 加密新密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqMap.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		admin.PasswordHash = string(hashedPassword)
	}

	admin.UpdatedAt = time.Now()

	// 保存更新
	err = s.adminRepo.Update(admin)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

// CreateAdmin 创建管理员
func (s *AdminService) CreateAdmin(req interface{}) (*model.Admin, error) {
	// 根据请求参数创建管理员
	// 这里需要根据具体的请求结构进行转换
	reqMap, ok := req.(*struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		ShopID   *uint  `json:"shop_id"`
		Role     int8   `json:"role" binding:"required"`
	})
	if !ok {
		return nil, errors.New("无效的请求参数")
	}

	// 检查用户名是否已存在
	exists, err := s.CheckUsernameExists(reqMap.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqMap.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建管理员
	admin := &model.Admin{
		Username:     reqMap.Username,
		PasswordHash: string(hashedPassword),
		ShopID:       reqMap.ShopID,
		Role:         reqMap.Role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.adminRepo.Create(admin)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

// CheckSuperAdminExists 检查是否存在超级管理员
func (s *AdminService) CheckSuperAdminExists() (bool, error) {
	admins, err := s.adminRepo.GetAll()
	if err != nil {
		return false, err
	}

	for _, admin := range admins {
		if admin.Role == 1 { // 假设角色1表示超级管理员
			return true, nil
		}
	}

	return false, nil
}

// CreateSuperAdmin 创建超级管理员
func (s *AdminService) CreateSuperAdmin(username string, password string, role int8) (*model.Admin, error) {
	// 检查用户名是否已存在
	exists, err := s.CheckUsernameExists(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建超级管理员
	admin := &model.Admin{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.adminRepo.Create(admin)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

// GetAllAdmins 获取所有管理员
func (s *AdminService) GetAllAdmins() ([]*model.Admin, error) {
	return s.adminRepo.GetAll()
}

// DeleteAdmin 删除管理员
func (s *AdminService) DeleteAdmin(id uint) error {
	return s.adminRepo.Delete(id)
}

// GetStatistics 获取统计数据
func (s *AdminService) GetStatistics(shopID uint, days int) (*model.StatisticsResponse, error) {
	// 获取当天统计数据
	today := time.Now().Format("20060102")
	todayStats, err := s.statisticRepo.GetDailyStatistics(shopID, today)
	if err != nil {
		return nil, err
	}

	// 获取昨天统计数据
	yesterday := time.Now().AddDate(0, 0, -1).Format("20060102")
	yesterdayStats, err := s.statisticRepo.GetDailyStatistics(shopID, yesterday)
	if err != nil {
		return nil, err
	}

	// 获取趋势数据
	trend, err := s.statisticRepo.GetStatisticsByDateRange(shopID, days)
	if err != nil {
		return nil, err
	}

	// 计算环比数据
	var comparison struct {
		ServedPercentage float64 `json:"served_percentage"`
		QueuePercentage  float64 `json:"queue_percentage"`
		WaitPercentage   float64 `json:"wait_percentage"`
	}

	if yesterdayStats.ServedCount > 0 {
		comparison.ServedPercentage = float64(todayStats.ServedCount-yesterdayStats.ServedCount) / float64(yesterdayStats.ServedCount) * 100
	}

	if yesterdayStats.QueueCount > 0 {
		comparison.QueuePercentage = float64(todayStats.QueueCount-yesterdayStats.QueueCount) / float64(yesterdayStats.QueueCount) * 100
	}

	if yesterdayStats.AvgWaitTime > 0 {
		comparison.WaitPercentage = float64(todayStats.AvgWaitTime-yesterdayStats.AvgWaitTime) / float64(yesterdayStats.AvgWaitTime) * 100
	}

	return &model.StatisticsResponse{
		Today:      todayStats,
		Yesterday:  yesterdayStats,
		Trend:      trend,
		Comparison: comparison,
	}, nil
}
