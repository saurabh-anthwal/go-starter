package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db     *gorm.DB
	models []interface{}
	zlog   = GetLogger()
)

// GetDB returns an initialized instance of gorm.DB
func GetDB() *gorm.DB {
	if db == nil {
		InitDatabase()
	}
	return db
}

// InitDatabase initializes the database.
func InitDatabase() *gorm.DB {
	var (
		err  error
		gCfg = &gorm.Config{
			Logger: GetGormLogger().LogMode(logger.Silent),
		}
	)

	db, err = gorm.Open(sqlite.Open("crispr.db"), gCfg)
	if err != nil {
		zlog.Fatalf("failed to connect to database, got: %v", err)
	}
	zlog.Infof("database connected successfully")
	return db
}

func RegisterModel(model interface{}) {
	models = append(models, model)
}

// Auto migrations here.
func AutoMigrate() {
	zlog.Info("running database migrations")
	if err := db.AutoMigrate(models...); err != nil {
		zlog.Fatalf("failed running db migrations, Got error: %v", err)
	}
}
