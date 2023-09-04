package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormprom "gorm.io/plugin/prometheus"
)

type News struct {
	gorm.Model
	Title   string
	Content string
}

func connect(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", config.Db.Host, config.Db.Port, config.Db.User, config.Db.Pass, config.Db.Name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&News{})
	db.Use(gormprom.New(gormprom.Config{
		DBName:          "db",
		RefreshInterval: 15,
		StartServer:     false,
	}))
	return db, nil
}
