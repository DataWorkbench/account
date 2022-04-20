package executor

import (
	"errors"
	"strings"
	"time"

	"github.com/DataWorkbench/common/qerror"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/DataWorkbench/common/constants"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
)

type User struct {
	UserID        string `gorm:"column:user_id;"`
	UserName      string `gorm:"column:user_name;"`
	Password      string `gorm:"column:password"`
	Lang          string `gorm:"column:lang;default:'cn'"`
	Email         string `gorm:"column:email;"`
	Phone         string `gorm:"column:phone;"`
	Status        string `gorm:"column:status;"`
	Role          string `gorm:"column:role;default:'user'"`
	Currency      string `gorm:"column:currency;default:'CNY'"`
	GravatarEmail string `gorm:"column:gravatar_email;"`
	Privilege     int32  `gorm:"column:privilege;default:1"`
	Zones         string `gorm:"column:zones;"`
	Regions       string `gorm:"column:regions;"`
	StatusTime    int64  `gorm:"column:status_time"`
	CreateTime    int64  `gorm:"column:create_time"`
}

func (u User) TableName() string {
	return constants.UserTableName
}

func (u *User) ToUserReply() *pbmodel.User {
	return &pbmodel.User{
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
		Zones:         strings.Split(u.Zones, ","),
		Regions:       strings.Split(u.Regions, ","),
	}
}

func (dbe *DBExecutor) ListUsers(
	db *gorm.DB, user_ids []string, limit int, offset int) (u []*User, err error) {

	query := "user_id in ?"
	var args []interface{}
	args = append(args, user_ids)

	err = db.Table(constants.UserTableName).
		Select(constants.UserColumns).
		Where(query, args...).Limit(limit).Offset(offset).Scan(&u).Error
	if err != nil {
		return
	}
	return
}

func (dbe *DBExecutor) CountUsers(
	db *gorm.DB, user_ids []string) (count int64, err error) {

	query := "user_id in ?"
	var args []interface{}
	args = append(args, user_ids)

	err = db.Table(constants.UserTableName).
		Select(constants.UserColumns).
		Where(query, args...).Count(&count).Error
	if err != nil {
		return
	}
	return
}

func checkUserInfoIsConflict(tx *gorm.DB, user *User) (err error) {
	existUser := new(User)
	var conflictClause []clause.Expression
	if user.UserName != "" {
		conflictClause = append(conflictClause, clause.Eq{Column: "user_name", Value: user.UserName})
	}
	if user.Email != "" {
		conflictClause = append(conflictClause, clause.Eq{Column: "email", Value: user.Email})
	}
	if user.Phone != "" {
		conflictClause = append(conflictClause, clause.Eq{Column: "phone", Value: user.Phone})
	}
	if len(conflictClause) == 0 {
		return
	}
	err = tx.Table(constants.UserTableName).Select(constants.UserColumns).Clauses(clause.Where{
		Exprs: []clause.Expression{
			clause.Or(conflictClause...),
			clause.Neq{Column: "status", Value: constants.UserStatusDelete},
			clause.Neq{Column: "user_id", Value: user.UserID},
		},
	}).Take(existUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	if existUser.UserName == user.UserName {
		err = qerror.ResourceAlreadyExists.Format(user.UserName)
		return
	}
	if existUser.Phone == user.Phone {
		err = qerror.ResourceAlreadyExists.Format(user.Phone)
		return
	}
	if existUser.Email == user.Email {
		err = qerror.ResourceAlreadyExists.Format(user.Email)
		return
	}

	return
}

func (dbe *DBExecutor) CreateUser(
	db *gorm.DB, user *User) (err error) {
	err = checkUserInfoIsConflict(db, user)
	if err != nil {
		return
	}
	if err = db.Table(constants.UserTableName).Create(user).Error; err != nil {
		return
	}
	return
}

func (dbe *DBExecutor) UpdateUser(
	db *gorm.DB, user *User) (err error) {
	err = checkUserInfoIsConflict(db, user)
	if err != nil {
		return
	}
	updateUserInfo := &User{
		UserName:   user.UserName,
		Phone:      user.Phone,
		Email:      user.Email,
		Lang:       user.Lang,
		Currency:   user.Currency,
		StatusTime: time.Now().Unix(),
	}
	err = db.Table(constants.UserTableName).Where("user_id = ? AND status = ?", user.UserID, constants.UserStatusActive).
		Updates(updateUserInfo).Error
	return err
}

func (dbe *DBExecutor) DeleteUser(tx *gorm.DB, userId string) (err error) {
	err = tx.Table(constants.UserTableName).Clauses(clause.Where{
		Exprs: []clause.Expression{
			clause.Neq{Column: "status", Value: constants.UserStatusDelete},
			clause.Eq{Column: "user_id", Value: userId},
		},
	}).Updates(map[string]interface{}{
		"status":      constants.UserStatusDelete,
		"status_time": time.Now().Unix(),
	}).Error
	if err != nil {
		return
	}
	err = tx.Table(constants.AccessKeyTableName).Clauses(clause.Where{
		Exprs: []clause.Expression{
			clause.Eq{Column: "owner", Value: userId},
		},
	}).Updates(map[string]interface{}{
		"status":      constants.AccessKeyStatusDisable,
		"status_time": time.Now().Unix(),
	}).Error
	if err != nil {
		return
	}
	return
}

func (dbe *DBExecutor) GetUserByName(tx *gorm.DB, userName string) (*User, error) {
	var user User

	err := tx.Table(constants.UserTableName).
		Select(constants.UserColumns).
		Where(map[string]interface{}{"user_name": userName, "status": constants.UserStatusActive}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dbe *DBExecutor) GetUserById(tx *gorm.DB, userId string) (*User, error) {
	var user User

	err := tx.Table(constants.UserTableName).
		Select(constants.UserColumns).
		Where(map[string]interface{}{"user_id": userId, "status": constants.UserStatusActive}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
