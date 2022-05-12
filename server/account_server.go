package server

import (
	"context"
	"github.com/DataWorkbench/account/handler"
	"github.com/DataWorkbench/gproto/xgo/service/pbsvcaccount"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
)

// AccountServer implements grpc server Interface pbsvcaccount.AccountServer
type AccountServer struct {
	pbsvcaccount.UnimplementedAccountServer
}

func (s *AccountServer) ValidateRequestSignature(ctx context.Context,
	req *pbrequest.ValidateRequestSignature) (*pbresponse.ValidateRequestSignature, error) {
	secretKey, err := handler.ValidateRequestSignature(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pbresponse.ValidateRequestSignature{
		Status:  200,
		Message: "",
		UserId:  secretKey.Owner,
	}, nil
}

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

func (s *AccountServer) UpdateUserRole(ctx context.Context, req *pbrequest.UpdateUser) (*pbresponse.UpdateUser, error) {
	output, err := handler.UpdateUserRole(ctx, req)
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.UpdateUser{
		User: output,
	}
	return reply, nil
}

func (s *AccountServer) UpdateUserZones(ctx context.Context, req *pbrequest.UpdateUser) (*pbresponse.UpdateUser, error) {
	output, err := handler.UpdateUserZones(ctx, req)
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.UpdateUser{
		User: output,
	}
	return reply, nil
}

func (s *AccountServer) UpdateUserPassword(ctx context.Context, req *pbrequest.UpdateUserPassword) (*pbmodel.EmptyStruct, error) {
	err := handler.UpdateUserPassword(ctx, req)
	if err != nil {
		return nil,err
	}
	return &pbmodel.EmptyStruct{}, nil
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

func (s *AccountServer) GetProvider(ctx context.Context, req *pbrequest.GetProvider) (*pbresponse.GetProvider, error) {
	provider, err := handler.GetProvider(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pbresponse.GetProvider{
		Provider: provider,
	}, nil
}

func (s *AccountServer) CreateProvider(ctx context.Context, req *pbrequest.CreateProvider) (*pbresponse.CreateProvider, error) {
	provider, err := handler.CreateProvider(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pbresponse.CreateProvider{
		Provider: provider,
	}, nil
}


func (s *AccountServer) DeleteProvider(ctx context.Context, req *pbrequest.DeleteProvider) (*pbresponse.DeleteProvider, error) {
	err := handler.DeleteProvider(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pbresponse.DeleteProvider{
		Name: req.Name,
	}, nil
}
