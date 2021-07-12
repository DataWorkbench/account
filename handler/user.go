package handler

import (
	"context"

	"github.com/DataWorkbench/account/internal/source"
	"github.com/DataWorkbench/gproto/pkg/accountpb"
)

func DescribeUsers(ctx context.Context, req *accountpb.DescribeUsersRequest) ([]*accountpb.User, int64, error) {
	users, totalCount, err := source.SelectSource(req.ReqSource, cfg, ctx).DescribeUsers(req.Users, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, 0, err
	}
	userSet := make([]*accountpb.User, len(users))
	for index, User := range users {
		userSet[index] = User.ToUserReply()
	}
	return userSet, totalCount, err

}
