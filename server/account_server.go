package server

import (
	"context"

	"github.com/DataWorkbench/account/handler"
	"github.com/DataWorkbench/gproto/pkg/accountpb"
)

// WorkspaceServer implements grpc server Interface accountpb.AccountServer
type AccountServer struct {
	accountpb.UnimplementedAccountServer
}

func (s *AccountServer) ValidateRequestSignature(ctx context.Context, req *accountpb.ValidateRequestSignatureRequest) (*accountpb.ValidateRequestSignatureReply, error) {
	secretKey, err := handler.ValidateRequestSignature(ctx, req)
	if err != nil {
		return nil, err
	}
	return &accountpb.ValidateRequestSignatureReply{
		Status:  200,
		Message: "",
		UserId:  secretKey.Owner,
	}, nil
}

func (s *AccountServer) DescribeUsers(ctx context.Context, req *accountpb.DescribeUsersRequest) (*accountpb.DescribeUsersReply, error) {
	users, totalCount, err := handler.DescribeUsers(ctx, req)
	if err != nil {
		return nil, err
	}
	reply := &accountpb.DescribeUsersReply{
		Status:     0,
		Message:    "",
		TotalCount: totalCount,
		UserSet:    users,
	}
	return reply, nil
}
