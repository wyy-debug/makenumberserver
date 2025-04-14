package repository

import (
	"example.com/henna-queue/internal/model"
)

// AdminRepository 管理员仓库接口
type AdminRepository interface {
	GetByID(id uint) (*model.Admin, error)
	GetByUsername(username string) (*model.Admin, error)
	GetAll() ([]*model.Admin, error)
	Create(admin *model.Admin) error
	Update(admin *model.Admin) error
	Delete(id uint) error
	CheckUsernameExists(username string) (bool, error)
}

// UserRepository 用户仓库接口
type UserRepository interface {
	GetByID(id string) (*model.User, error)
	GetByPhone(phone string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
}

// ShopRepository 店铺仓库接口
type ShopRepository interface {
	GetByID(id uint) (*model.Shop, error)
	Create(shop *model.Shop) error
	Update(shop *model.Shop) error
	Delete(id uint) error
}

// ServiceRepository 服务仓库接口
type ServiceRepository interface {
	GetByID(id uint) (*model.Service, error)
	GetByShopID(shopID uint) ([]*model.Service, error)
	GetAll() ([]*model.Service, error)
	Create(service *model.Service) error
	Update(service *model.Service) error
	Delete(id uint) error
}

// QueueRepository 排队仓库接口
type QueueRepository interface {
	GetByID(id uint) (*model.Queue, error)
	GetByStatus(shopID uint, status int8) (*model.Queue, error)
	GetActiveByUserID(userID string, shopID uint) (*model.Queue, error)
	GetNextWaiting(shopID uint) (*model.Queue, error)
	GetActiveQueues(shopID uint) ([]*model.Queue, error)
	GetDailyCount(shopID uint, date string) (int, error)
	GetTodayServedCount(shopID uint) (int, error)
	GetPeopleAhead(shopID uint, queueNumber string) (int, error)
	Create(queue *model.Queue) error
	Update(queue *model.Queue) error
}

// StatisticRepository 统计仓库接口
type StatisticRepository interface {
	GetDailyStatistics(shopID uint, date string) (*model.DailyStatistics, error)
	GetStatisticsByDateRange(shopID uint, days int) ([]*model.DailyStatistics, error)
	IncrementServedCount(shopID uint) error
	IncrementCancelCount(shopID uint) error
}
