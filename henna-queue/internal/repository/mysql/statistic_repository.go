package mysql

import (
	"example.com/henna-queue/internal/model"
	"example.com/henna-queue/internal/repository"
	"example.com/henna-queue/pkg/db"
	"time"
)

type statisticRepository struct{}

func NewStatisticRepository() repository.StatisticRepository {
	return &statisticRepository{}
}

func (r *statisticRepository) GetDailyStatistics(shopID uint, date string) (*model.DailyStatistics, error) {
	var stats model.DailyStatistics
	err := db.DB.Where("shop_id = ? AND date = ?", shopID, date).First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *statisticRepository) GetStatisticsByDateRange(shopID uint, days int) ([]*model.DailyStatistics, error) {
	var stats []*model.DailyStatistics
	endDate := time.Now().Format("20060102")
	startDate := time.Now().AddDate(0, 0, -days+1).Format("20060102")
	err := db.DB.Where("shop_id = ? AND date BETWEEN ? AND ?", shopID, startDate, endDate).Order("date").Find(&stats).Error
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *statisticRepository) IncrementServedCount(shopID uint) error {
	today := time.Now().Format("20060102")
	return db.DB.Exec("INSERT INTO daily_statistics (shop_id, date, served_count) VALUES (?, ?, 1) ON DUPLICATE KEY UPDATE served_count = served_count + 1", shopID, today).Error
}

func (r *statisticRepository) IncrementCancelCount(shopID uint) error {
	today := time.Now().Format("20060102")
	return db.DB.Exec("INSERT INTO daily_statistics (shop_id, date, cancel_count) VALUES (?, ?, 1) ON DUPLICATE KEY UPDATE cancel_count = cancel_count + 1", shopID, today).Error
}
