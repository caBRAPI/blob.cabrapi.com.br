package database

import (
	"os"
	"time"

	"blob/src/functions"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Postgres() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	}

	now := time.Now().Format("Mon Jan 2 15:04:05 2006")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		functions.Error("[POSTGRES ERROR] %v (%s)", err, now)
		return
	}

	sqlDB, err := DB.DB()
	if err != nil {
		functions.Error("[POSTGRES ERROR] %v (%s)", err, now)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		functions.Error("[POSTGRES ERROR] %v (%s)", err, now)
	} else {
		functions.Info("[POSTGRES] Connected successfully. (%s)", now)
	}
}
