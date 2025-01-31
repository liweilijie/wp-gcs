package database

import (
	"fmt"
	logging "github.com/ipfs/go-log/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"wp-gcs/config"
)

var log = logging.Logger("database")

func InitDB(cfg config.AppConfig) *gorm.DB {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB_USERNAME, cfg.DB_PASSWORD, cfg.DB_HOSTNAME, cfg.DB_PORT, cfg.DB_NAME)
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		log.Fatal("error connect to DB", err.Error())
		return nil
	}

	return db
}

func InitialMigration(db *gorm.DB) {
	err := db.AutoMigrate(WpUploads{})
	if err != nil {
		log.Error("auto migrate error:", err)
	}
}

type WpUploadsHandle interface {
	Insert(WpUploads) error
	SelectByNames(localPath string, bucketPath string) ([]WpUploads, error)
	SelectByLocalPath(localPath string) ([]WpUploads, error)
}
