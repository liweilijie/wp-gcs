package database

import (
	"errors"
	"gorm.io/gorm"
)

type WpUploadsQuery struct {
	db *gorm.DB
}

func (repo *WpUploadsQuery) SelectByNames(localPath string, bucketPath string) ([]WpUploads, error) {
	// if the length of localPath or bucketPath  is more than 255
	// we should cut short than 255 from prefix and to find
	if len(localPath) > 255 {
		localPath = localPath[len(localPath)-255:]
	}

	if len(bucketPath) > 255 {
		bucketPath = bucketPath[len(bucketPath)-255:]
	}

	dataModel := []WpUploads{}

	tx := repo.db.Where("origin_path = ? or bucket_path = ?", localPath, bucketPath).First(&dataModel)
	if tx.Error != nil {
		return []WpUploads{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return []WpUploads{}, nil
	}
	// 分批次查询速度更快

	//tx := repo.db.Where("origin_path = ? or bucket_path = ?", localPath, bucketPath).First(&dataModel)
	//if tx.Error != nil {
	//	return []WpUploads{}, tx.Error
	//}
	//if tx.RowsAffected == 0 {
	//	return []WpUploads{}, nil
	//}

	return dataModel, nil
}

// Insert implements WpUploads
func (repo *WpUploadsQuery) Insert(data WpUploads) error {
	if len(data.OriginPath) > 255 {
		data.OriginPath = data.OriginPath[len(data.OriginPath)-255:]
	}
	if len(data.BucketPath) > 255 {
		data.BucketPath = data.BucketPath[len(data.BucketPath)-255:]
	}
	tx := repo.db.Create(&data)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("insert error, row affected = 0")
	}
	return nil
}

func New(db *gorm.DB) WpUploadsHandle {
	return &WpUploadsQuery{
		db: db,
	}
}
