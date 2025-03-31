package factory

import (
	"example.com/henna-queue/internal/repository"
	"example.com/henna-queue/internal/repository/mysql"
)

// NewAdminRepository 创建管理员仓库
func NewAdminRepository() repository.AdminRepository {
	return mysql.NewAdminRepository()
}

// NewUserRepository 创建用户仓库
func NewUserRepository() repository.UserRepository {
	return mysql.NewUserRepository()
}

// NewShopRepository 创建店铺仓库
func NewShopRepository() repository.ShopRepository {
	return mysql.NewShopRepository()
}

// NewServiceRepository 创建服务仓库
func NewServiceRepository() repository.ServiceRepository {
	return mysql.NewServiceRepository()
}

// NewQueueRepository 创建排队仓库
func NewQueueRepository() repository.QueueRepository {
	return mysql.NewQueueRepository()
}

// NewStatisticRepository 创建统计仓库
func NewStatisticRepository() repository.StatisticRepository {
	return mysql.NewStatisticRepository()
}
