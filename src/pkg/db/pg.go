package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPG(opts *Options) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		opts.Host, opts.Port, opts.Username, opts.Password, opts.Database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: opts.Logger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)       //设置数据库连接池最大连接数
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime) //
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)       //如果没有sql任务需要执行的连接数大于20， //超过的连接会被连接池关闭

	return db, nil
}
