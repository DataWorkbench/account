package executor

import (
	"context"

	"github.com/DataWorkbench/common/constants"
	"github.com/DataWorkbench/gproto/pkg/accountpb"
)

type User struct {
	UserID        string `gorm:"column:user_id;"`
	UserName      string `gorm:"column:user_name;"`
	Lang          string `gorm:"column:lang;"`
	Email         string `gorm:"column:email;"`
	Phone         string `gorm:"column:phone;"`
	Status        string `gorm:"column:status;"`
	Role          string `gorm:"column:role;"`
	Currency      string `gorm:"column:currency;"`
	GravatarEmail string `gorm:"column:gravatar_email;"`
	Privilege     int32  `gorm:column:privilege;"`
}

func (u User) TableName() string {
	return constants.UserTableName
}

func (u *User) ToUserReply() *accountpb.User {
	return &accountpb.User{
		UserId:        u.UserID,
		UserName:      u.UserName,
		Lang:          u.Lang,
		Email:         u.Email,
		Phone:         u.Phone,
		Status:        u.Status,
		Role:          u.Role,
		Currency:      u.Currency,
		GravatarEmail: u.GravatarEmail,
		Privilege:     u.Privilege,
	}
}

func (dbe *DBExecutor) ListUsers(
	ctx context.Context, user_ids []string, limit int, offset int) (u []*User, err error) {

	query := "user_id in ?"
	var args []interface{}
	args = append(args, user_ids)

	db := dbe.db.WithContext(ctx)
	err = db.Table(constants.UserTableName).
		Select(constants.UserColumns).
		Where(query, args...).Limit(limit).Offset(offset).Scan(&u).Error
	if err != nil {
		return
	}
	return
}

func (dbe *DBExecutor) CountUsers(
	ctx context.Context, user_ids []string) (count int64, err error) {

	query := "user_id in ?"
	var args []interface{}
	args = append(args, user_ids)

	db := dbe.db.WithContext(ctx)
	err = db.Table(constants.UserTableName).
		Select(constants.UserColumns).
		Where(query, args...).Count(&count).Error
	if err != nil {
		return
	}
	return
}
