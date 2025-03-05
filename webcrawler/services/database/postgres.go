package database

import (
	"book-search/webcrawler/config"
	"book-search/webcrawler/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDBClient() (*gorm.DB, error) {
	env, err := config.GetEnv()
	if err != nil {
		return nil, err
	}

	dsn := env.PostgresDSN

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	err = migrateTables(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func migrateTables(db *gorm.DB) error {
	return db.AutoMigrate(&models.Book{}, &models.Author{})
}

func CloseDBClient(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
