package service

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"henna-queue/internal/model"
	"henna-queue/internal/repository/mysql"
)

type AdminService struct {
	adminRepo    *mysql.AdminRepository
	statisticRepo *mysql.StatisticRepository
}

func NewAdminService() *AdminService {
	return &AdminService{
		adminRepo:    mysql.NewAdminRepository(),
		statisticRepo: mysql.NewStatisticRepository(),
	}
}

// GetAdmin 获取管理员
func (s *AdminService) GetAdmin(adminID uint) (*model.Admin, error) {
	return s.adminRepo.GetByID(adminID)
}

// UpdateAdmin 更新管理员
func (s *AdminService) UpdateAdmin(adminID uint, req interface{}) (*model.Admin, error) {
	reqMap, ok := req.(*struct {
		Username    string `json:"username"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	})
	if !ok {
		return nil, errors.New("无效的请求参数")
	}
	
	admin, err := s.adminRepo.GetByID(adminID)
	if err != nil {
		return nil, errors.New("管理员不存在")
	}
	
	// 更新用户名
	if reqMap.Username != "" && reqMap.Username != admin.Username {
		// 检查用户名是否已存在
		exists, _ := s.adminRepo.CheckUsernameExists(reqMap.Username)
		if exists {
			return nil, errors.New("用户名已存在")
		}
		
		admin.Username = reqMap.Username
	}
	
	// 更新密码
	if reqMap.OldPassword != "" && reqMap.NewPassword != "" {
		// 验证旧密码
		err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(reqMap.OldPassword))
		if err != nil {
			return nil, errors.New("旧密码错误")
		}
		
		// 加密新密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqMap.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		
		admin.PasswordHash = string(hashedPassword)
	}
	
	admin.UpdatedAt = time.Now()
	
	err = s.adminRepo.Update(admin)
	if err != nil {
		return nil, err
	}
	
	return admin, nil
}

// CreateAdmin 创建管理员
func (s *AdminService) CreateAdmin(req interface{}) (*model.Admin, error) {
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
	exists, _ := s.adminRepo.CheckUsernameExists(reqMap.Username)
	if exists {
		return nil, errors.New("用户名已存在")
	}
	
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqMap.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	
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

// GetAllAdmins 获取所有管理员
func (s *AdminService) GetAllAdmins() ([]*model.Admin, error) {
	return s.adminRepo.GetAll()
}

// DeleteAdmin 删除管理员
func (s *AdminService) DeleteAdmin(adminID uint) error {
	admin, err := s.adminRepo.GetByID(adminID)
	if err != nil {
		return errors.New("管理员不存在")
	}
	
	if admin.Role == 2 {
		return errors.New("不能删除超级管理员")
	}
	
	return s.adminRepo.Delete(adminID)
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