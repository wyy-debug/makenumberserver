package db

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitMySQL 初始化MySQL连接
func InitMySQL() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.database"),
		viper.GetString("mysql.charset"),
		viper.GetString("mysql.loc"))

	logLevel := logger.Silent
	if viper.GetString("app.mode") == "development" {
		logLevel = logger.Info
	}

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))
	sqlDB.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}

// CloseMySQL 关闭MySQL连接
func CloseMySQL() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
} 