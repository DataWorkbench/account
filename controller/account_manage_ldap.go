package controller

import (
	"context"

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
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (x *AccountManagerLdap) ChangeUserPassword(ctx context.Context, req *pbrequest.ChangeUserPassword) (*pbmodel.EmptyStruct, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeUserPassword not implemented")
}
func (x *AccountManagerLdap) ResetUserPassword(ctx context.Context, req *pbrequest.ResetUserPassword) (*pbmodel.EmptyStruct, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetUserPassword not implemented")
}
func (x *AccountManagerLdap) ListAccessKeys(ctx context.Context, req *pbrequest.ListAccessKeys) (*pbresponse.ListAccessKeys, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAccessKeys not implemented")
}
func (x *AccountManagerLdap) DeleteAccessKeys(ctx context.Context, req *pbrequest.DeleteAccessKeys) (*pbmodel.EmptyStruct, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAccessKeys not implemented")
}
func (x *AccountManagerLdap) DescribeAccessKey(ctx context.Context, req *pbrequest.DescribeAccessKey) (*pbresponse.DescribeAccessKey, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeAccessKey not implemented")
}
func (x *AccountManagerLdap) CreateAccessKey(ctx context.Context, req *pbrequest.CreateAccessKey) (*pbresponse.CreateAccessKey, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccessKey not implemented")
}
func (x *AccountManagerLdap) UpdatedAccessKey(ctx context.Context, req *pbrequest.UpdatedAccessKey) (*pbmodel.EmptyStruct, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatedAccessKey not implemented")
}
func (x *AccountManagerLdap) CreateSession(ctx context.Context, req *pbrequest.CreateSession) (*pbresponse.CreateSession, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSession not implemented")
}
func (x *AccountManagerLdap) CheckSession(ctx context.Context, req *pbrequest.CheckSession) (*pbresponse.CheckSession, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckSession not implemented")
}
func (x *AccountManagerLdap) ListNotifications(ctx context.Context, req *pbrequest.ListNotifications) (*pbresponse.ListNotifications, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListNotifications not implemented")
}
func (x *AccountManagerLdap) DescribeNotification(ctx context.Context, req *pbrequest.DescribeNotification) (*pbresponse.DescribeNotification, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeNotification not implemented")
}
func (x *AccountManagerLdap) CreateNotification(ctx context.Context, req *pbrequest.CreateNotification) (*pbresponse.CreateNotification, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNotification not implemented")
}
func (x *AccountManagerLdap) UpdateNotification(ctx context.Context, req *pbrequest.UpdateNotification) (*pbmodel.EmptyStruct, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateNotification not implemented")
}
func (x *AccountManagerLdap) DeleteNotifications(ctx context.Context, req *pbrequest.DeleteNotifications) (*pbmodel.EmptyStruct, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteNotifications not implemented")
}
