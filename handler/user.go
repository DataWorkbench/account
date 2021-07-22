package handler

import (
	"context"
	"strings"

	"github.com/DataWorkbench/account/internal/source"
	"github.com/DataWorkbench/gproto/pkg/accountpb"
)

func DescribeUsers(ctx context.Context, req *accountpb.DescribeUsersRequest) ([]*accountpb.User, int64, error) {
	cachedUsers, uncachedUsers, notExistsUsers, err := cache.GetCachedUsers(req.Users, req.ReqSource)
	logger.Debug().String("uncachedUsers", strings.Join(uncachedUsers, ",")).Fire()
	logger.Debug().String("notExistsUsers", strings.Join(notExistsUsers, ",")).Fire()
	if err != nil {
		return nil, 0, err
	}
	var users []source.User
	var totalCount int64
	var userSet []*accountpb.User
	if len(uncachedUsers) != 0 {
		users, totalCount, err = source.SelectSource(req.ReqSource, cfg, ctx).DescribeUsers(uncachedUsers, int(req.Limit)-len(cachedUsers), int(req.Offset))
		if err != nil {
			return nil, 0, err
		}
		userMap := make(map[string]bool)
		for index, user := range users {
			userSet = append(userSet, user.ToUserReply())
			userMap[userSet[index].UserId] = true
		}
		if totalCount != int64(len(uncachedUsers)) {
			for _, userID := range uncachedUsers {
				if _, ok := userMap[userID]; !ok {
					notExistsUsers = append(notExistsUsers, userID)
				}
			}
		}
	}
	if len(notExistsUsers) != 0 {
		logger.Warn().String("Users not Exists", strings.Join(notExistsUsers, ",")).Fire()
		for _, userID := range notExistsUsers {
			cache.CacheNotExistUser(userID, req.ReqSource)
		}
	}

	totalCount += int64(len(cachedUsers))
	userSet = append(userSet, cachedUsers...)
	if err := cache.CacheUsers(userSet, req.ReqSource); err != nil {
		return nil, 0, err
	}

	return userSet, totalCount, err

}
