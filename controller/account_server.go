package controller

import (
	"context"

	"github.com/DataWorkbench/account/handler"
	"github.com/DataWorkbench/gproto/xgo/service/pbsvcaccount"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
)

// AccountServer implements grpc server Interface pbsvcaccount.AccountServer
type AccountServer struct {
	pbsvcaccount.UnimplementedAccountServer
}

//func (s *AccountServer) ValidateRequestSignature(ctx context.Context,
//	req *pbrequest.ValidateRequestSignature) (*pbresponse.ValidateRequestSignature, error) {
//	secretKey, err := handler.ValidateRequestSignature(ctx, req)
//	if err != nil {
//		return nil, err
//	}
//	return &pbresponse.ValidateRequestSignature{
//		Status:  200,
//		Message: "",
//		UserId:  secretKey.Owner,
//	}, nil
//}

func (s *AccountServer) DescribeUsers(ctx context.Context, req *pbrequest.DescribeUsers) (*pbresponse.DescribeUsers, error) {
	users, totalCount, err := handler.DescribeUsers(ctx, req)
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.DescribeUsers{
		Status:     0,
		Message:    "",
		TotalCount: totalCount,
		UserSet:    users,
	}
	return reply, nil
}

func (s *AccountServer) DescribeAccessKey(ctx context.Context, req *pbrequest.DescribeAccessKey) (*pbresponse.DescribeAccessKey, error) {
	output, err := handler.DescribeAccessKey(ctx, req)
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.DescribeAccessKey{
		Owner:           output.Owner,
		SecretAccessKey: output.SecretAccessKey,
	}
	return reply, nil
}

func (s *AccountServer) CreateUser(ctx context.Context, req *pbrequest.CreateUser) (*pbresponse.CreateUser, error) {
	output, err := handler.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.CreateUser{
		User: output,
	}
	return reply, nil
}

func (s *AccountServer) UpdateUser(ctx context.Context, req *pbrequest.UpdateUser) (*pbresponse.UpdateUser, error) {
	output, err := handler.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.UpdateUser{
		User: output,
	}
	return reply, nil
}

func (s *AccountServer) DeleteUser(ctx context.Context, req *pbrequest.DeleteUser) (*pbresponse.DeleteUser, error) {
	err := handler.DeleteUser(ctx, req)
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.DeleteUser{
		UserId: req.UserId,
	}
	return reply, nil
}

func (s *AccountServer) CheckSession(ctx context.Context, req *pbrequest.CheckSession) (*pbresponse.CheckSession, error) {
	session, err := handler.CheckSession(ctx, req)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *AccountServer) CreateSession(ctx context.Context, req *pbrequest.CreateSession) (*pbresponse.CreateSession, error) {
	session, err := handler.CreateSession(ctx, req)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *AccountServer) CheckUserExists(ctx context.Context, req *pbrequest.CheckUserExists) (*pbresponse.CheckUserExists, error) {
	resp, err := handler.CheckUserExists(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
