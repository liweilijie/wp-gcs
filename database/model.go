package database

import "gorm.io/gorm"

type WpUploads struct {
	gorm.Model
	Id         uint   `gorm:"primary_key;auto_increment"`
	OriginPath string `gorm:"index:idx_path,unique"`
	BucketPath string `gorm:"index:idx_path,unique"`
}
