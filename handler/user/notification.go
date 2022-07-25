package user

import (
	"github.com/DataWorkbench/common/gormwrap"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ListNotifications(tx *gorm.DB, input *pbrequest.ListNotifications) (output *pbresponse.ListNotifications, err error) {
	var exprs []clause.Expression
	if input.UserId != "" {
		exprs = append(exprs, clause.Eq{Column: "owner", Value: input.UserId})
	}
	if input.Search != "" {
		like1 := clause.Like{Column: "name", Value: "%" + input.Search + "%"}
		like2 := clause.Like{Column: "email", Value: "%" + input.Search + "%"}
		or := clause.Or(like1, like2)
		exprs = append(exprs, or)
	}
	if len(input.NfIds) != 0 {
		exprs = append(exprs, gormwrap.BuildConditionClauseInWithString("id", input.NfIds))
	}
	var infos []*pbmodel.Notification
	var total int64

	err = tx.Table(tableNameNotification).Select("*").Clauses(clause.Where{Exprs: exprs}).
		Limit(int(input.Limit)).Offset(int(input.Offset)).
		Scan(&infos).Error
	if err != nil {
		return nil, err
	}
	if input.Offset == 0 && len(infos) < int(input.Limit) {
		total = int64(len(infos))
	} else {
		err = tx.Table(tableNameNotification).Select("count(id)").Clauses(clause.Where{Exprs: exprs}).Count(&total).Error
		if err != nil {
			return nil, err
		}
	}
	output = &pbresponse.ListNotifications{
		Infos:   infos,
		Total:   total,
		HasMore: len(infos) >= int(input.Limit),
	}
	return output, nil
}

func CreateNotification(tx *gorm.DB, owner, id, name, description, email string) (err error) {
	conflict := CheckUpdateNotificationConflict(tx, owner, email)
	if conflict {
		return qerror.ResourceAlreadyExists
	}
	err = tx.Table(tableNameNotification).Create(&pbmodel.Notification{
		Owner:       owner,
		Id:          id,
		Name:        name,
		Description: description,
		Email:       email,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func DescribeNotification(tx *gorm.DB, nfId string) (resp *pbresponse.DescribeNotification, err error) {
	var info pbmodel.Notification
	err = tx.Table(tableNameNotification).Where("id = ?", nfId).Scan(&info).Error
	if err != nil {
		return nil, err
	}
	resp = &pbresponse.DescribeNotification{
		NfSet: &info,
	}
	return resp, nil
}

func DeleteNotifications(tx *gorm.DB, nfIds []string) (err error) {
	var exprs []clause.Expression
	if len(nfIds) != 0 {
		exprs = append(exprs, gormwrap.BuildConditionClauseInWithString("id", nfIds))
	}
	_ = exprs
	for _, id := range nfIds {
		de := tx.Table(tableNameNotification).Where("id = ?", id).Delete(&pbmodel.Notification{})
		if de.Error != nil {
			return err
		}
	}
	return nil
}

func UpdateNotification(tx *gorm.DB, nfId, name, description, email string) (err error) {
	err = tx.Table(tableNameNotification).Clauses(clause.Where{Exprs: []clause.Expression{
		clause.Eq{Column: "id", Value: nfId},
	}}).Updates(&pbmodel.Notification{
		Name:        name,
		Description: description,
		Email:       email,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func CheckUpdateNotificationConflict(tx *gorm.DB, owner, email string) bool {
	var count int64
	tx.Table(tableNameNotification).Where("owner = ? AND email = ?", owner, email).Count(&count)
	return count > 0
}
