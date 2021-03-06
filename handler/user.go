package handler

import (
	"context"
	"errors"
	"github.com/DataWorkbench/account/executor"
	"github.com/DataWorkbench/account/internal/source"
	"github.com/DataWorkbench/common/constants"
	"github.com/DataWorkbench/common/gormwrap"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/common/utils/password"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
	"gorm.io/gorm"
	"strings"
	"time"
)

func getUsers(userIds []string, source string) ([]*pbmodel.User, []string, []string, error) {
	uncachedUsers := []string{}
	cachedUsers := []*pbmodel.User{}
	notExistsUsers := []string{}
	for i := 0; i < len(userIds); i++ {
		user, err := cache.GetUser(userIds[i], source)
		if err != nil {
			if err == qerror.ResourceNotExists {
				notExistsUsers = append(notExistsUsers, userIds[i])
				continue
			}
			logger.Error().String("Get user from cache error", err.Error())
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

func DescribeUsers(ctx context.Context, req *pbrequest.DescribeUsers) ([]*pbmodel.User, int64, error) {
	if req.ReqSource == "" {
		req.ReqSource = executor.AccountExecutor.GetConf().Source
	}
	cachedUsers, uncachedUsers, notExistsUsers, err := getUsers(req.Users, req.ReqSource)
	if err != nil {
		return nil, 0, err
	}
	var users []source.User
	var totalCount int64
	var userSet []*pbmodel.User
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

func CreateUser(ctx context.Context, req *pbrequest.CreateUser) (*pbmodel.User, error) {
	id, err := IdGeneratorUser.Take()
	if err != nil {
		return nil, err
	}
	req.Password, err = password.Encode(req.Password)
	if err != nil {
		return nil, err
	}
	var userInfo *executor.User
	err = gormwrap.ExecuteFuncWithTxn(ctx, executor.AccountExecutor.Db, func(tx *gorm.DB) error {
		var xErr error

		userInfo = &executor.User{
			UserID:     id,
			UserName:   req.UserName,
			Password:   req.Password,
			Lang:       req.Lang,
			Email:      req.Email,
			Phone:      req.Phone,
			Status:     constants.UserStatusActive,
			Currency:   req.Currency,
			CreateTime: time.Now().Unix(),
			StatusTime: time.Now().Unix(),
		}

		if xErr = executor.AccountExecutor.CreateUser(tx, userInfo); xErr != nil {
			return xErr
		}
		// create access key
		accessKeyId, secretKey := password.RandomGenerateAccessKey()
		accessKey := &executor.AccessKey{
			AccessKeyID:     accessKeyId,
			SecretAccessKey: secretKey,
			Owner:           id,
			Status:          constants.AccessKeyStatusEnable,
			CreateTime:      time.Now().Unix(),
			StatusTime:      time.Now().Unix(),
		}
		if xErr = executor.AccountExecutor.CreateAccessKey(tx, accessKey); xErr != nil {
			return xErr
		}
		return xErr
	})

	if err != nil {
		return nil, err
	}
	return userInfo.ToUserReply(), nil
}

func UpdateUser(ctx context.Context, req *pbrequest.UpdateUser) (*pbmodel.User, error) {
	var userInfo *executor.User
	var err error
	req.Password, err = password.Encode(req.Password)
	if err != nil {
		return nil, err
	}
	err = gormwrap.ExecuteFuncWithTxn(ctx, executor.AccountExecutor.Db, func(tx *gorm.DB) error {
		var xErr error

		userInfo = &executor.User{
			UserID: req.UserId,
			//UserName: req.UserName,
			Email:    req.Email,
			Phone:    req.Phone,
			Lang:     req.Lang,
			Currency: req.Currency,
		}

		if xErr = executor.AccountExecutor.UpdateUser(tx, userInfo); xErr != nil {
			return xErr
		}
		return xErr
	})

	if err != nil {
		return nil, err
	}
	if ignoreError := cache.DelUser(req.UserId, false); ignoreError != nil {
		logger.Warn().String("delete user cache error", ignoreError.Error()).Fire()
	}
	return userInfo.ToUserReply(), nil
}

func DeleteUser(ctx context.Context, req *pbrequest.DeleteUser) error {
	err := gormwrap.ExecuteFuncWithTxn(ctx, executor.AccountExecutor.Db, func(tx *gorm.DB) error {
		var xErr error

		if xErr = executor.AccountExecutor.DeleteUser(tx, req.UserId); xErr != nil {
			return xErr
		}
		return xErr
	})
	if ignoreError := cache.DelUser(req.UserId, true); ignoreError != nil {
		logger.Warn().String("delete user cache error", ignoreError.Error()).Fire()
	}
	return err
}

func CheckUserExists(ctx context.Context, req *pbrequest.CheckUserExists) (*pbresponse.CheckUserExists, error) {
	resp := &pbresponse.CheckUserExists{}
	var xErr error
	var user *executor.User
	err := gormwrap.ExecuteFuncWithTxn(ctx, executor.AccountExecutor.Db, func(tx *gorm.DB) error {
		if user, xErr = executor.AccountExecutor.GetUserByName(tx, req.UserName); xErr != nil {
			if errors.Is(xErr, gorm.ErrRecordNotFound) {
				xErr = nil
				user = &executor.User{}
			}
		}
		return xErr
	})
	if err != nil {
		return resp, err
	}
	if user.UserID != "" && user.Status == constants.UserStatusActive {
		resp.Exists = true
	}
	return resp, nil
}
