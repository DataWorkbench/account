package user

import (
	"time"

	secret2 "github.com/DataWorkbench/account/handler/user/internal/secret"
	"github.com/DataWorkbench/common/gormwrap"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ListUsers(tx *gorm.DB, input *pbrequest.ListUsers) (output *pbresponse.ListUsers, err error) {
	// Build where exprs.
	exprs := []clause.Expression{
		clause.Neq{Column: "status", Value: pbmodel.User_deleted.Number()},
	}
	if input.Name != "" {
		exprs = append(exprs, clause.Eq{Column: "name", Value: input.Name})
	}
	if len(input.UserIds) != 0 {
		exprs = append(exprs, gormwrap.BuildConditionClauseInWithString("user_id", input.UserIds))
	}

	var infos []*pbmodel.User
	var total int64

	err = tx.Table(tableNameUser).Select("*").Clauses(clause.Where{Exprs: exprs}).
		Limit(int(input.Limit)).Offset(int(input.Offset)).
		Scan(&infos).Error
	if err != nil {
		return
	}

	// Length of result(infos) less than the limit means no more records with the query conditions.
	if input.Offset == 0 && len(infos) < int(input.Limit) {
		total = int64(len(infos))
	} else {
		err = tx.Table(tableNameUser).Select("count(user_id)").Clauses(clause.Where{Exprs: exprs}).Count(&total).Error
		if err != nil {
			return
		}
	}

	output = &pbresponse.ListUsers{
		Infos:   infos,
		Total:   total,
		HasMore: len(infos) >= int(input.Limit),
	}
	return
}

func CreateUser(tx *gorm.DB, userId, name, password, email string) (err error) {
	// Check user name is conflict
	var x string
	err = tx.Table(tableNameUser).Select("user_id").Clauses(clause.Where{Exprs: []clause.Expression{
		clause.Eq{Column: "name", Value: name},
		clause.Neq{Column: "status", Value: pbmodel.User_deleted.Number()},
	}}).Take(&x).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if err != nil {
		return
	}
	if x != "" {
		err = qerror.ResourceAlreadyExists.Format(name)
		return
	}
	password, err = secret2.EncodePassword(password)
	if err != nil {
		return
	}
	now := time.Now().Unix()
	info := &pbmodel.User{
		UserId:   userId,
		Name:     name,
		Email:    email,
		Role:     pbmodel.User_User,
		Status:   pbmodel.User_active,
		Password: password,
		Created:  now,
		Updated:  now,
	}

	err = tx.Table(tableNameUser).Create(&info).Error
	if err != nil {
		return
	}
	return
}

func DescribeUserById(tx *gorm.DB, userId string) (info *pbmodel.User, err error) {
	info = new(pbmodel.User)
	err = tx.Table(tableNameUser).Select("*").Clauses(clause.Where{Exprs: []clause.Expression{
		clause.Eq{Column: "user_id", Value: userId},
		clause.Neq{Column: "status", Value: pbmodel.User_deleted.Number()},
	}}).Take(&info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = qerror.ResourceNotExists.Format(userId)
		}
		return
	}
	return
}

func DescribeUserByName(tx *gorm.DB, userName string) (info *pbmodel.User, err error) {
	info = new(pbmodel.User)
	err = tx.Table(tableNameUser).Select("*").Clauses(clause.Where{Exprs: []clause.Expression{
		clause.Eq{Column: "name", Value: userName},
		clause.Neq{Column: "status", Value: pbmodel.User_deleted.Number()},
	}}).Take(&info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = qerror.UserNotExists.Format(userName)
		}
		return
	}
	return
}

func ChangePassword(tx *gorm.DB, userId, oldPassWord, newPassWord string) (err error) {
	user := &pbmodel.User{}
	res := tx.Table(tableNameUser).Where("user_id = ?", userId).Find(&user)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return qerror.ResourceNotExists.Format(userId)
	}
	verify := secret2.CheckPassword(oldPassWord, user.Password)
	if !verify {
		return qerror.UserNameOrPasswordError
	}
	encodePassword, err := secret2.EncodePassword(newPassWord)
	if err != nil {
		return err
	}
	updates := tx.Table(tableNameUser).Where("user_id = ?", userId).Update("password", encodePassword)
	if updates.Error != nil {
		return updates.Error
	}
	return nil
}

func ResetPassword(tx *gorm.DB, userId, newPassWord string) (err error) {
	var count int64
	tx.Table(tableNameUser).Where("user_id = ? AND status != ?", userId, pbmodel.User_deleted).Count(&count)
	if err != nil {
		return
	}
	if count == 0 {
		return qerror.ResourceNotExists
	}
	encodePassword, err := secret2.EncodePassword(newPassWord)
	if err != nil {
		return err
	}
	err = tx.Table(tableNameUser).Where("user_id = ?", userId).Update("password", encodePassword).Error
	return
}

func ResetPasswordByName(tx *gorm.DB, userName, newPassWord string) (err error) {
	var count int64
	err = tx.Table(tableNameUser).Where("name = ? AND status != ?", userName, pbmodel.User_deleted).Count(&count).Error
	if err != nil {
		return
	}
	if count == 0 {
		return qerror.ResourceNotExists
	}
	encodePassword, err := secret2.EncodePassword(newPassWord)
	if err != nil {
		return err
	}
	err = tx.Table(tableNameUser).Where("name = ? AND status != ?", userName, pbmodel.User_deleted).Update("password", encodePassword).Error
	return
}

func DeleteUserByIds(tx *gorm.DB, userIds []string) (err error) {
	expr := gormwrap.BuildConditionClauseInWithString("user_id", userIds)

	err = tx.Table(tableNameUser).Clauses(clause.Where{
		Exprs: []clause.Expression{
			clause.Neq{Column: "status", Value: pbmodel.User_deleted.Number()},
			expr,
		},
	}).Updates(map[string]interface{}{
		"status":  pbmodel.User_deleted.Number(),
		"name":    gorm.Expr("concat(name,'.', user_id)"),
		"updated": time.Now().Unix(),
	}).Error
	if err != nil {
		return
	}
	return
}

func DeleteUserByNames(tx *gorm.DB, userNames string) (err error) {

	err = tx.Table(tableNameUser).Clauses(clause.Where{
		Exprs: []clause.Expression{
			clause.Neq{Column: "status", Value: pbmodel.User_deleted.Number()},
			clause.Eq{Column: "name", Value: userNames},
		},
	}).Updates(map[string]interface{}{
		"status":  pbmodel.User_deleted.Number(),
		"name":    gorm.Expr("concat(name,'.', user_id)"),
		"updated": time.Now().Unix(),
	}).Error
	if err != nil {
		return
	}
	return
}

func UpdateUser(tx *gorm.DB, userid, email string) error {
	update := tx.Table(tableNameUser).Where("user_id = ?", userid).Update("email", email)
	if update.Error != nil {
		return update.Error
	}
	return nil
}

func ExistsUsername(tx *gorm.DB, name string) bool {
	var count int64
	tx.Table(tableNameUser).Where("name = ? and status != ?", name, pbmodel.User_deleted).Count(&count)
	return count > 0
}

func CreateAdminUser(tx *gorm.DB, userId, username, password, email string) error {
	var err error
	ok := ExistsUsername(tx, username)

	encodePassword, err := secret2.EncodePassword(password)
	if !ok {
		err = tx.Transaction(func(tx *gorm.DB) error {
			err = tx.Table(tableNameUser).Create(&pbmodel.User{
				UserId:   userId,
				Name:     username,
				Email:    email,
				Role:     pbmodel.User_Admin,
				Status:   pbmodel.User_active,
				Password: encodePassword,
				Created:  time.Now().Unix(),
				Updated:  time.Now().Unix(),
			}).Error
			if err != nil {
				return err
			}
			err = InitAccessKey(tx, userId)
			if err != nil {
				return err
			}
			return nil
		})
	}
	return err
}

func AddUser(tx *gorm.DB, userId, username, password, email string) error {
	var err error
	ok := ExistsUsername(tx, username)

	encodePassword, err := secret2.EncodePassword(password)
	if !ok {
		err = tx.Transaction(func(tx *gorm.DB) error {
			err = tx.Table(tableNameUser).Create(&pbmodel.User{
				UserId:   userId,
				Name:     username,
				Email:    email,
				Role:     pbmodel.User_User,
				Status:   pbmodel.User_active,
				Password: encodePassword,
				Created:  time.Now().Unix(),
				Updated:  time.Now().Unix(),
			}).Error
			if err != nil {
				return err
			}
			err = InitAccessKey(tx, userId)
			if err != nil {
				return err
			}
			return nil
		})
	}
	return err
}

func UpdateUserByUsername(tx *gorm.DB, userName, email string) error {
	var err error
	err = tx.Table(tableNameUser).Where("name = ?", userName).Update("email", email).Error
	return err
}
