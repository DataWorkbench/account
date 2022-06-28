package user

import (
	"time"

	"github.com/DataWorkbench/account/handler/user/internal/secret"
	"github.com/DataWorkbench/common/gormwrap"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// InitAccessKey for init a pitrix access key.
func InitAccessKey(tx *gorm.DB, owner string) (err error) {
	ak, sk := secret.GenerateAccessKey()
	now := time.Now().Unix()
	info := &pbmodel.AccessKey{
		AccessKeyId:     ak,
		SecretAccessKey: sk,
		Owner:           owner,
		Name:            "pitrix access key for mgmt",
		Controller:      pbmodel.AccessKey_pitrix,
		Status:          pbmodel.AccessKey_active,
		Description:     sk,
		IpWhiteList:     "",
		Created:         now,
		Updated:         now,
	}

	if err = tx.Table(tableNameAccessKey).Create(info).Error; err != nil {
		return
	}
	return
}

func CreateAccessKey(tx *gorm.DB, keySet *pbmodel.AccessKey) (err error) {
	panic("unrealized")
}

func ListAccessKeys(tx *gorm.DB) {
	panic("unrealized")
}

func DescribeAccessKeyByKeyId(tx *gorm.DB, accessKeyId string) (info *pbmodel.AccessKey, err error) {
	info = new(pbmodel.AccessKey)
	err = tx.Table(tableNameAccessKey).Select("*").Clauses(clause.Where{Exprs: []clause.Expression{
		clause.Eq{Column: "access_key_id", Value: accessKeyId},
		clause.Neq{Column: "status", Value: pbmodel.AccessKey_deleted.Number()},
	}}).Take(&info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = qerror.ResourceNotExists.Format(accessKeyId)
		}
		return
	}
	return
}

func DescribePitrixAccessKeyByOwner(tx *gorm.DB, owner string) (info *pbmodel.AccessKey, err error) {
	info = new(pbmodel.AccessKey)
	err = tx.Table(tableNameAccessKey).Select("*").Clauses(clause.Where{Exprs: []clause.Expression{
		clause.Eq{Column: "owner", Value: owner},
		clause.Eq{Column: "controller", Value: pbmodel.AccessKey_pitrix},
		clause.Neq{Column: "status", Value: pbmodel.AccessKey_deleted.Number()},
	}}).Take(&info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = qerror.ResourceNotExists.Format(owner)
		}
		return
	}
	return
}

func UpdateAccessKey(tx *gorm.DB) {
	panic("unrealized")
}

func DeleteAccessKeys(tx *gorm.DB, accessKeyIds []string) (err error) {
	if len(accessKeyIds) == 0 {
		return
	}
	expr := gormwrap.BuildConditionClauseInWithString("access_key_id", accessKeyIds)
	return deleteAccessKeyByExpr(tx, expr)
}

func DeleteAccessKeysByUserIds(tx *gorm.DB, userIds []string) (err error) {
	if len(userIds) == 0 {
		return
	}
	expr := gormwrap.BuildConditionClauseInWithString("owner", userIds)
	return deleteAccessKeyByExpr(tx, expr)
}

func deleteAccessKeyByExpr(tx *gorm.DB, expr clause.Expression) (err error) {
	err = tx.Table(tableNameAccessKey).Clauses(clause.Where{
		Exprs: []clause.Expression{
			clause.Neq{Column: "status", Value: pbmodel.AccessKey_deleted.Number()},
			expr,
		},
	}).Updates(map[string]interface{}{
		"status":  pbmodel.AccessKey_deleted.Number(),
		"updated": time.Now().Unix(),
	}).Error
	if err != nil {
		return
	}
	return
}

func DescriptAccessKey(tx *gorm.DB, accesskeyId string) (key *pbmodel.AccessKey, err error) {
	key = &pbmodel.AccessKey{}
	err = tx.Table(tableNameAccessKey).Where("access_key_id = ?", accesskeyId).Find(key).Error
	if err != nil {
		return
	}
	return
}
