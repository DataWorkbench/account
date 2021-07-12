package executor

import (
	"context"

	"github.com/DataWorkbench/common/constants"
)

type AccessKey struct {
	AccessKeyID     string `gorm:"column:access_key_id;"`
	AccessKeyName   string `gorm:"column:access_key_name;"`
	SecretAccessKey string `gorm:"column:secret_access_key;"`
	Description     string `gorm:"column:description;"`
	Owner           string `gorm:"column:owner;"`
	Status          string `gorm:"column:status;"`
	IPWhiteList     string `gorm:"column:ip_white_list;"`
	CreateTime      int64  `gorm:"column:create_time;"`
	StatusTime      int64  `gorm:"column:status_time;"`
}

func (k AccessKey) TableName() string {
	return constants.AccessKeyTableName
}

func (dbe *DBExecutor) ListAccessKeys(
	ctx context.Context, accessKeyIDs []string, owner string, limit int, offset int) (k []*AccessKey, err error) {

	query := "access_key_id in ?"
	var args []interface{}
	args = append(args, accessKeyIDs)

	if owner != "" {
		query += " and owner = ?"
		args = append(args, owner)
	}

	db := dbe.db.WithContext(ctx)
	err = db.Table(constants.AccessKeyTableName).
		Select(constants.AccessKeyColumns).
		Where(query, args...).Limit(limit).Offset(offset).Scan(&k).Error
	if err != nil {
		return
	}
	return
}
