package user

import (
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
	"gorm.io/gorm"
)

type Notification struct {
	Owner       string `yaml:"owner" gorm:"column:owner;"`
	Id          string `yaml:"id" gorm:"column:id;primaryKey"`
	Name        string `yaml:"name" gorm:"column:name;"`
	Description string `yaml:"description" gorm:"column:description;"`
	Email       string `yaml:"email" gorm:"column:email;"`
	Created     int64  `yaml:"created" gorm:"column:created;autoCreateTime;"`
	Updated     int64  `yaml:"updated" gorm:"column:updated;autoUpdateTime;"`
}

func (n Notification) TableName() string {
	return tableNameNotification
}

func ListNotifications(tx *gorm.DB, input *pbrequest.ListNotifications) (output *pbresponse.ListNotifications, err error) {
	var count int64
	var notificatuions []Notification
	offset := input.Offset
	limit := input.Limit
	userid := input.UserId
	res := tx.Where("owner = ?", userid).Count(&count).Offset(int((offset - 1) * limit)).Limit(int(limit)).Find(&notificatuions)
	err = res.Error
	if err != nil {
		return nil, err
	}
	total := int64(len(notificatuions))
	reply := &pbresponse.ListNotifications{
		Infos:   nil,
		Total:   total,
		HasMore: false,
	}
	var temp []*pbmodel.Notification
	if count > total {
		reply.HasMore = true
	}
	for _, v := range notificatuions {
		n := &pbmodel.Notification{
			Owner:       v.Owner,
			Id:          v.Id,
			Name:        v.Name,
			Description: v.Description,
			Email:       v.Email,
			Created:     v.Created,
			Updated:     v.Updated,
		}
		temp = append(temp, n)
	}
	reply.Infos = temp
	return reply, nil
}

func CreateNotification(tx *gorm.DB, owner, id, name, description, email string) (err error) {
	notification := &Notification{}
	err = tx.Table(tableNameNotification).Where("name = ? OR email = ?", name, email).Find(notification).Error
	if err != nil {
		return err
	}
	if notification.Name != "" || notification.Email != "" {
		return qerror.ResourceAlreadyExists
	}
	err = tx.Create(&Notification{
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
		de := tx.Where("id = ?", id).Delete(&Notification{})
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
	var n = &Notification{
		Id: id,
	}
	updates := tx.Model(n).Updates(&Notification{
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

func notification2PbModel(notification *Notification) *pbmodel.Notification {
	return &pbmodel.Notification{
		Owner:       notification.Owner,
		Id:          notification.Id,
		Name:        notification.Name,
		Description: notification.Description,
		Email:       notification.Email,
		Created:     notification.Created,
		Updated:     notification.Updated,
	}
}
