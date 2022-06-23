package user

import (
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
	"gorm.io/gorm"
	"math"
)

func ListNotifications(tx *gorm.DB, input *pbrequest.ListNotifications) (output *pbresponse.ListNotifications, err error) {
	var count int64
	var notifications []*pbmodel.Notification
	offset := input.Offset
	limit := input.Limit
	userid := input.UserId
	if offset <= 0 {
		offset = 1
	}
	if limit <= 0 {
		limit = 10
	}
	res := tx.Table(tableNameNotification).Where("owner = ?", userid).Count(&count).Offset(int((offset - 1) * limit)).Limit(int(limit)).Find(&notifications)
	err = res.Error
	if err != nil {
		return nil, err
	}
	totalPage := math.Ceil(float64(count) / float64(limit))
	reply := &pbresponse.ListNotifications{
		Infos:   nil,
		Total:   count,
		HasMore: false,
	}
	reply.Infos = notifications
	if int64(offset) < int64(totalPage) {
		reply.HasMore = true
	}
	return reply, nil
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
	if count > 0 {
		return true
	}
	return false
}
