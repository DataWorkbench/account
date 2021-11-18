package handler

import (
	"context"
	"strings"

	"github.com/DataWorkbench/account/executor"
	"github.com/DataWorkbench/account/internal/source"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/gproto/pkg/accountpb"
)

func getUsers(userIds []string, source string) ([]*accountpb.User, []string, []string, error) {
	uncachedUsers := []string{}
	cachedUsers := []*accountpb.User{}
	notExistsUsers := []string{}
	for i := 0; i < len(userIds); i++ {
		user, err := cache.GetUser(userIds[i], source)
		if err != nil {
			if err == qerror.ResourceNotExists {
				notExistsUsers = append(notExistsUsers, userIds[i])
				continue
			}
			return nil, []string{}, []string{}, err
		}
		if user != nil {
			cachedUsers = append(cachedUsers, user)
		} else {
			uncachedUsers = append(uncachedUsers, userIds[i])
		}
	}
	cachedUserIds := []string{}
	for i := 0; i < len(cachedUsers); i++ {
		cachedUserIds = append(cachedUserIds, cachedUsers[i].UserId)
	}
	logger.Debug().String("cachedUsers", strings.Join(cachedUserIds, ",")).Fire()
	logger.Debug().String("uncachedUsers", strings.Join(uncachedUsers, ",")).Fire()
	logger.Debug().String("notExistsUsers", strings.Join(notExistsUsers, ",")).Fire()
	return cachedUsers, uncachedUsers, notExistsUsers, nil
}

func DescribeUsers(ctx context.Context, req *accountpb.DescribeUsersRequest) ([]*accountpb.User, int64, error) {
	if req.ReqSource == "" {
		req.ReqSource = executor.AccountExecutor.GetConf().Source
	}
	cachedUsers, uncachedUsers, notExistsUsers, err := getUsers(req.Users, req.ReqSource)
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

func DescribeAccessKey(ctx context.Context, input *accountpb.DescribeAccessKeyRequest) (output *executor.AccessKey, err error) {
	output, err = getAccessKey(ctx, &accountpb.ValidateRequestSignatureRequest{
		ReqAccessKeyId: input.AccessKeyId,
		ReqSource:  cfg.Source,

	})
	return
}