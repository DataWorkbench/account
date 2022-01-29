package server

import (
	"context"

	"github.com/DataWorkbench/account/handler"
	"github.com/DataWorkbench/gproto/pkg/service/pbsvcaccount"
	"github.com/DataWorkbench/gproto/pkg/types/pbrequest"
	"github.com/DataWorkbench/gproto/pkg/types/pbresponse"
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
