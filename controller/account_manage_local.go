package controller

import (
	"context"

	"github.com/DataWorkbench/account/handler/user"
	"github.com/DataWorkbench/account/options"
	"github.com/DataWorkbench/common/gormwrap"
	"github.com/DataWorkbench/gproto/xgo/service/pbsvcaccount"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

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
		if xErr := user.CreateUser(tx, userId, req.Name, req.Password, req.Email); err != nil {
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

// TODO
func (x *AccountManagerLocal) ChangeUserPassword(ctx context.Context, req *pbrequest.ChangeUserPassword) (*pbmodel.EmptyStruct, error) {
	tx := options.DBConn.WithContext(ctx)
	err := user.ChangePassword(tx, req.UserId, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, err
	}
	return &pbmodel.EmptyStruct{}, err
}
func (x *AccountManagerLocal) ResetUserPassword(context.Context, *pbrequest.ResetUserPassword) (*pbmodel.EmptyStruct, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetUserPassword not implemented")
}
func (x *AccountManagerLocal) DescribeAccessKey(ctx context.Context, req *pbrequest.DescribeAccessKey) (*pbresponse.DescribeAccessKey, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeAccessKey not implemented")
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
	return nil, status.Errorf(codes.Unimplemented, "method DescribeNotification not implemented")
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
