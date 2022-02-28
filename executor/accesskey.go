package executor

import (
	"gorm.io/gorm"

	"github.com/DataWorkbench/common/constants"
)

type AccessKey struct {
	AccessKeyID     string `gorm:"column:access_key_id;" json:"access_key_id"`
	AccessKeyName   string `gorm:"column:access_key_name;" json:"access_key_name"`
	SecretAccessKey string `gorm:"column:secret_access_key;" json:"secret_access_key"`
	Description     string `gorm:"column:description;" json:"description"`
	Owner           string `gorm:"column:owner;" json:"owner"`
	Status          string `gorm:"column:status;" json:"status"`
	IPWhiteList     string `gorm:"column:ip_white_list;" json:"ip_white_list"`
	CreateTime      int64  `gorm:"column:create_time;" json:"create_time"`
	StatusTime      int64  `gorm:"column:status_time;" json:"status_time"`
}

func (k AccessKey) TableName() string {
	return constants.AccessKeyTableName
}

func (dbe *DBExecutor) ListAccessKeys(
	db *gorm.DB, accessKeyIDs []string, owner string, limit int, offset int) (k []*AccessKey, err error) {

	query := "access_key_id in ?"
	var args []interface{}
	args = append(args, accessKeyIDs)

	if owner != "" {
		query += " and owner = ?"
		args = append(args, owner)
	}

	err = db.Table(constants.AccessKeyTableName).
		Select(constants.AccessKeyColumns).
		Where(query, args...).Limit(limit).Offset(offset).Scan(&k).Error
	if err != nil {
		return
	}
	return
}

func (dbe *DBExecutor) CreateAccessKey(
	db *gorm.DB, key *AccessKey) (err error) {

	if err = db.Table(constants.AccessKeyTableName).Create(key).Error; err != nil {
		return
	}
	return
}

func (dbe *DBExecutor) GetAccessKeyByOwner(tx *gorm.DB, owner string) (*AccessKey, error) {
	var key AccessKey

	err := tx.Table(constants.AccessKeyTableName).
		Select(constants.AccessKeyColumns).
		Where(map[string]interface{}{"owner": owner}).First(&key).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}
