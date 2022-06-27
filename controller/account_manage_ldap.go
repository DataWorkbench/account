package controller

import (
	"context"
	"github.com/DataWorkbench/account/handler/user"
	"github.com/DataWorkbench/account/options"
	"github.com/DataWorkbench/common/gormwrap"
	"gorm.io/gorm"

	"github.com/DataWorkbench/gproto/xgo/service/pbsvcaccount"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AccountManagerLdap implements grpc server Interface pbsvcaccount.AccountManagerLocal
type AccountManagerLdap struct {
	pbsvcaccount.UnimplementedAccountManageServer
}

func (x *AccountManagerLdap) ListUsers(ctx context.Context, req *pbrequest.ListUsers) (*pbresponse.ListUsers, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUsers not implemented")
}
func (x *AccountManagerLdap) DeleteUsers(ctx context.Context, req *pbrequest.DeleteUsers) (*pbmodel.EmptyStruct, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUsers not implemented")
}
func (x *AccountManagerLdap) DescribeUser(ctx context.Context, req *pbrequest.DescribeUser) (*pbresponse.DescribeUser, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeUser not implemented")
}
func (x *AccountManagerLdap) CreateUser(ctx context.Context, req *pbrequest.CreateUser) (*pbresponse.CreateUser, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (x *AccountManagerLdap) UpdateUser(ctx context.Context, req *pbrequest.UpdateUser) (*pbmodel.EmptyStruct, error) {
	tx := options.DBConn.WithContext(ctx)
	err := user.UpdateUser(tx, req.UserId, req.Email)
	if err != nil {
		return nil, err
	}
	return &pbmodel.EmptyStruct{}, nil
}
func (x *AccountManagerLdap) ChangeUserPassword(ctx context.Context, req *pbrequest.ChangeUserPassword) (*pbmodel.EmptyStruct, error) {
	//tx := options.DBConn.WithContext(ctx)
	//err := user.ChangePassword(tx, req.UserId, req.OldPassword, req.NewPassword)
	//if err != nil {
	//	return nil, err
	//}
	//return &pbmodel.EmptyStruct{}, err
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (x *AccountManagerLdap) ResetUserPassword(ctx context.Context, req *pbrequest.ResetUserPassword) (*pbmodel.EmptyStruct, error) {
	//tx := options.DBConn.WithContext(ctx)
	//err := user.ResetPassword(tx, req.UserId, req.NewPassword)
	//if err != nil {
	//	return nil, err
	//}
	//return &pbmodel.EmptyStruct{}, err
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (x *AccountManagerLdap) ListAccessKeys(ctx context.Context, req *pbrequest.ListAccessKeys) (*pbresponse.ListAccessKeys, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAccessKeys not implemented")
}
func (x *AccountManagerLdap) DeleteAccessKeys(ctx context.Context, req *pbrequest.DeleteAccessKeys) (*pbmodel.EmptyStruct, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAccessKeys not implemented")
}
func (x *AccountManagerLdap) DescribeAccessKey(ctx context.Context, req *pbrequest.DescribeAccessKey) (*pbresponse.DescribeAccessKey, error) {
	// TODO
	return nil, status.Errorf(codes.Unimplemented, "method DescribeAccessKey not implemented")
}
func (x *AccountManagerLdap) CreateAccessKey(ctx context.Context, req *pbrequest.CreateAccessKey) (*pbresponse.CreateAccessKey, error) {
	// TODO
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccessKey not implemented")
}
func (x *AccountManagerLdap) UpdatedAccessKey(ctx context.Context, req *pbrequest.UpdatedAccessKey) (*pbmodel.EmptyStruct, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatedAccessKey not implemented")
}
func (x *AccountManagerLdap) CreateSession(ctx context.Context, req *pbrequest.CreateSession) (*pbresponse.CreateSession, error) {
	tx := options.DBConn.WithContext(ctx)
	userName := req.UserName
	password := req.Password
	res, err := user.LdapProvider.Authentication(userName, password)
	if err != nil {
		return nil, err
	}
	mail := res["mail"].(string)
	exist := user.ExistsUsername(tx, userName)
	if !exist {
		userId, err := options.IdGeneratorUser.Take()
		if err != nil {
			return nil, err
		}
		err = gormwrap.ExecuteFuncWithTxn(ctx, options.DBConn, func(tx *gorm.DB) error {
			if xErr := user.CreateUser(tx, userId, userName, password, mail); err != nil {
				return xErr
			}
			if xErr := user.InitAccessKey(tx, userId); xErr != nil {
				return xErr
			}
			return nil
		})
	}
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
func (x *AccountManagerLdap) CheckSession(ctx context.Context, req *pbrequest.CheckSession) (*pbresponse.CheckSession, error) {
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
func (x *AccountManagerLdap) ListNotifications(ctx context.Context, req *pbrequest.ListNotifications) (*pbresponse.ListNotifications, error) {
	tx := options.DBConn.WithContext(ctx)
	notifications, err := user.ListNotifications(tx, req)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}
func (x *AccountManagerLdap) DescribeNotification(ctx context.Context, req *pbrequest.DescribeNotification) (*pbresponse.DescribeNotification, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeNotification not implemented")
}
func (x *AccountManagerLdap) CreateNotification(ctx context.Context, req *pbrequest.CreateNotification) (*pbresponse.CreateNotification, error) {
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
func (x *AccountManagerLdap) UpdateNotification(ctx context.Context, req *pbrequest.UpdateNotification) (*pbmodel.EmptyStruct, error) {
	tx := options.DBConn.WithContext(ctx)
	err := user.UpdateNotification(tx, req.NfId, req.Name, req.Description, req.Email)
	if err != nil {
		return nil, err
	}
	return &pbmodel.EmptyStruct{}, nil
}
func (x *AccountManagerLdap) DeleteNotifications(ctx context.Context, req *pbrequest.DeleteNotifications) (*pbmodel.EmptyStruct, error) {
	tx := options.DBConn.WithContext(ctx)
	err := user.DeleteNotifications(tx, req.NfIds)
	if err != nil {
		return nil, err
	}
	return &pbmodel.EmptyStruct{}, nil
}
