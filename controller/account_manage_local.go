package controller

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DataWorkbench/account/handler/user"
	"github.com/DataWorkbench/account/options"
	"github.com/DataWorkbench/common/gormwrap"
	"github.com/DataWorkbench/gproto/xgo/service/pbsvcaccount"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
	"gorm.io/gorm"
)

type sessionCache struct {
	UserSet *pbmodel.User
	KeySet  *pbmodel.AccessKey
}

// AccountManagerLocal implements grpc server Interface pbsvcaccount.AccountManagerLocal
type AccountManagerLocal struct {
	pbsvcaccount.UnimplementedAccountManageServer
}

func (x *AccountManagerLocal) ListUsers(ctx context.Context, req *pbrequest.ListUsers) (*pbresponse.ListUsers, error) {
	tx := options.DBConn.WithContext(ctx)
	users, err := user.ListUsers(tx, req)
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (x *AccountManagerLocal) DescribeUser(ctx context.Context, req *pbrequest.DescribeUser) (*pbresponse.DescribeUser, error) {
	tx := options.DBConn.WithContext(ctx)
	info, err := user.DescribeUserById(tx, req.UserId)
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.DescribeUser{UserSet: info}
	return reply, nil
}
func (x *AccountManagerLocal) CreateUser(ctx context.Context, req *pbrequest.CreateUser) (*pbresponse.CreateUser, error) {
	userId, err := options.IdGeneratorUser.Take()
	p := &pbresponse.CreateUser{
		Id: userId,
	}
	if err != nil {
		return nil, err
	}
	err = gormwrap.ExecuteFuncWithTxn(ctx, options.DBConn, func(tx *gorm.DB) error {
		if xErr := user.CreateUser(tx, userId, req.Name, req.Password, req.Email, pbmodel.User_Native); err != nil {
			return xErr
		}
		if xErr := user.InitAccessKey(tx, userId); xErr != nil {
			return xErr
		}
		return nil
	})
	return p, nil
}
func (x *AccountManagerLocal) UpdateUser(ctx context.Context, req *pbrequest.UpdateUser) (*pbmodel.EmptyStruct, error) {
	tx := options.DBConn.WithContext(ctx)
	err := user.UpdateUser(tx, req.UserId, req.Email)
	if err != nil {
		return nil, err
	}
	// 重新给sid赋值
	rdb := options.RedisClient
	userSet, err := user.DescribeUserById(tx, req.UserId)
	if err != nil {
		return nil, err
	}
	keySet, err := user.DescribePitrixAccessKeyByOwner(tx, userSet.UserId)
	if err != nil {
		return nil, err
	}
	sessionValue := &sessionCache{
		UserSet: userSet,
		KeySet:  keySet,
	}
	jsonString, err := json.Marshal(sessionValue)
	if err != nil {
		return nil, err
	}
	rdb.Set(ctx, req.SessionId, jsonString, time.Second*60*60*24)
	return &pbmodel.EmptyStruct{}, nil
}
func (x *AccountManagerLocal) DeleteUsers(ctx context.Context, req *pbrequest.DeleteUsers) (*pbmodel.EmptyStruct, error) {
	err := gormwrap.ExecuteFuncWithTxn(ctx, options.DBConn, func(tx *gorm.DB) error {
		if xErr := user.DeleteUserByIds(tx, req.UserIds); xErr != nil {
			return xErr
		}
		if xErr := user.DeleteAccessKeysByUserIds(tx, req.UserIds); xErr != nil {
			return xErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pbmodel.EmptyStruct{}, nil
}

func (x *AccountManagerLocal) ChangeUserPassword(ctx context.Context, req *pbrequest.ChangeUserPassword) (*pbmodel.EmptyStruct, error) {
	tx := options.DBConn.WithContext(ctx)
	err := user.ChangePassword(tx, req.UserId, req.OldPassword, req.NewPassword)
	if err != nil {
		return &pbmodel.EmptyStruct{}, err
	}
	return &pbmodel.EmptyStruct{}, err
}
func (x *AccountManagerLocal) ResetUserPassword(ctx context.Context, req *pbrequest.ResetUserPassword) (*pbmodel.EmptyStruct, error) {
	tx := options.DBConn.WithContext(ctx)
	err := user.ResetPassword(tx, req.UserId, req.NewPassword)
	if err != nil {
		return &pbmodel.EmptyStruct{}, err
	}
	return &pbmodel.EmptyStruct{}, nil
}
func (x *AccountManagerLocal) DescribeAccessKey(ctx context.Context, req *pbrequest.DescribeAccessKey) (*pbresponse.DescribeAccessKey, error) {
	tx := options.DBConn.WithContext(ctx)
	key, err := user.DescriptAccessKey(tx, req.AccessKeyId)
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.DescribeAccessKey{
		KeySet: key,
	}
	return reply, nil
}
func (x *AccountManagerLocal) CreateSession(ctx context.Context, req *pbrequest.CreateSession) (*pbresponse.CreateSession, error) {
	tx := options.DBConn.WithContext(ctx)
	userSet, sessionId, err := user.CreateSession(ctx, tx, options.RedisClient, req.UserName, req.Password, req.IgnorePassword)
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.CreateSession{
		SessionId: sessionId,
		UserSet:   userSet,
	}

	return reply, nil
}
func (x *AccountManagerLocal) CheckSession(ctx context.Context, req *pbrequest.CheckSession) (*pbresponse.CheckSession, error) {
	userSet, keySet, err := user.CheckSession(ctx, options.RedisClient, req.SessionId)
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.CheckSession{
		UserSet: userSet,
		KeySet:  keySet,
	}
	return reply, nil
}
func (x *AccountManagerLocal) ListNotifications(ctx context.Context, req *pbrequest.ListNotifications) (*pbresponse.ListNotifications, error) {
	tx := options.DBConn.WithContext(ctx)
	notifications, err := user.ListNotifications(tx, req)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}
func (x *AccountManagerLocal) DescribeNotification(ctx context.Context, req *pbrequest.DescribeNotification) (*pbresponse.DescribeNotification, error) {
	tx := options.DBConn.WithContext(ctx)
	notification, err := user.DescribeNotification(tx, req.NfId)
	if err != nil {
		return nil, err
	}
	return notification, nil
}
func (x *AccountManagerLocal) CreateNotification(ctx context.Context, req *pbrequest.CreateNotification) (*pbresponse.CreateNotification, error) {
	take, err := options.IdGeneratorNotification.Take()
	if err != nil {
		return nil, err
	}
	tx := options.DBConn.WithContext(ctx)
	err = user.CreateNotification(tx, req.UserId, take, req.Name, req.Description, req.Email)
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.CreateNotification{
		Id: take,
	}
	return reply, nil
}
func (x *AccountManagerLocal) UpdateNotification(ctx context.Context, req *pbrequest.UpdateNotification) (*pbmodel.EmptyStruct, error) {
	tx := options.DBConn.WithContext(ctx)
	err := user.UpdateNotification(tx, req.NfId, req.Name, req.Description, req.Email)
	if err != nil {
		return nil, err
	}
	return &pbmodel.EmptyStruct{}, nil
}
func (x *AccountManagerLocal) DeleteNotifications(ctx context.Context, req *pbrequest.DeleteNotifications) (*pbmodel.EmptyStruct, error) {
	tx := options.DBConn.WithContext(ctx)
	err := user.DeleteNotifications(tx, req.NfIds)
	if err != nil {
		return nil, err
	}
	return &pbmodel.EmptyStruct{}, nil
}

func CreateAdminUser(ctx context.Context) error {
	tx := options.DBConn.WithContext(ctx)
	adminId, err := options.IdGeneratorUser.Take()
	if err != nil {
		return err
	}
	hash := sha256.New()
	_, err = hash.Write([]byte("zhu88jie"))
	if err != nil {
		return err
	}
	password := hex.EncodeToString(hash.Sum(nil))
	err = user.CreateAdminUser(tx, adminId, "admin", password, "account@yunify.com")
	return err
}

func (x *AccountManagerLocal) CreateSessionAuth(ctx context.Context, req *pbrequest.CreateSession) (reply *pbresponse.CreateSession, err error) {
	tx := options.DBConn.WithContext(ctx)

	ok := user.ExistsUsername(tx, req.UserName)
	// create user
	if !ok {
		id, err := options.IdGeneratorUser.Take()
		if err != nil {
			return nil, err
		}
		hash := sha256.New()
		_, err = hash.Write([]byte(req.UserName))
		if err != nil {
			return nil, err
		}
		password := hex.EncodeToString(hash.Sum(nil))

		err = gormwrap.ExecuteFuncWithTxn(ctx, options.DBConn, func(tx *gorm.DB) error {
			if xErr := user.CreateUser(tx, id, req.UserName, password, fmt.Sprintf("%s@%s.com", req.UserName, req.UserName), pbmodel.User_Native); err != nil {
				return xErr
			}
			if xErr := user.InitAccessKey(tx, id); xErr != nil {
				return xErr
			}
			return nil
		})
	}

	userSet, sessionId, err := user.CreateSessionAuth(ctx, tx, options.RedisClient, req.UserName)
	if err != nil {
		return nil, err
	}
	reply = &pbresponse.CreateSession{
		SessionId: sessionId,
		UserSet:   userSet,
	}

	return reply, nil
}
