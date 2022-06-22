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
	err = tx.Table(tableNameUser).Select("id").Clauses(clause.Where{Exprs: []clause.Expression{
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
			err = qerror.ResourceNotExists.Format(userName)
		}
		return
	}
	return
}

func ChangePassword(tx *gorm.DB, userId, oldPassWord, newPassWord string) (err error) {
	panic("unrealized")
}

func RestPassword(tx *gorm.DB, userId, newPassWord string) (err error) {
	panic("unrealized")
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
