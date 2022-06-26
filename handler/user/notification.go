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
		err = tx.Table(tableNameNotification).Select("count(id)").Clauses(clause.Where{Exprs: exprs}).Scan(&total).Error
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
	notification := &pbmodel.Notification{}
	err = tx.Table(tableNameNotification).Where("name = ? OR email = ?", name, email).Find(&pbmodel.Notification{}).Error
	if err != nil {
		return err
	}
	if notification.Name != "" || notification.Email != "" {
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

func DeleteNotifications(tx *gorm.DB, nfIds []string) (err error) {
	for _, id := range nfIds {
		de := tx.Table(tableNameNotification).Where("id = ?", id).Delete(&pbmodel.Notification{})
		if de.Error != nil {
			return err
		}
	}
	return nil
}

func UpdateNotification(tx *gorm.DB, id, name, description, email string) (err error) {
	conflict := CheckUpdateNotificationConflict(tx, id, name, email)
	if conflict {
		return qerror.ResourceAlreadyExists
	}
	updates := tx.Table(tableNameNotification).Where("id = ?", id).Updates(&pbmodel.Notification{
		Name:        name,
		Description: description,
		Email:       email,
	})
	err = updates.Error
	if err != nil {
		return err
	}
	return nil
}

func CheckUpdateNotificationConflict(tx *gorm.DB, id, name, email string) bool {
	var count int64
	tx.Table(tableNameNotification).Where("id <> ? and (name = ? or email = ?)", id, name, email).Count(&count)
	return count > 0
}
