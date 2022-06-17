package handler

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/DataWorkbench/account/executor"
	"github.com/DataWorkbench/account/internal/source"
	"github.com/DataWorkbench/account/options"
	"github.com/DataWorkbench/common/constants"
	"github.com/DataWorkbench/common/gormwrap"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/common/utils/password"
	"github.com/DataWorkbench/glog"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
	"gorm.io/gorm"
)

func getUsers(ctx context.Context, userIds []string, source string) ([]*pbmodel.User, []string, []string, error) {
	lg := glog.FromContext(ctx)
	var uncachedUsers []string
	var cachedUsers []*pbmodel.User
	var notExistsUsers []string
	for i := 0; i < len(userIds); i++ {
		user, err := cache.GetUser(ctx, userIds[i], source)
		if err != nil {
			if err == qerror.ResourceNotExists {
				notExistsUsers = append(notExistsUsers, userIds[i])
				continue
			}
			lg.Error().String("Get user from cache error", err.Error())
		}
		if user != nil {
			cachedUsers = append(cachedUsers, user)
		} else {
			uncachedUsers = append(uncachedUsers, userIds[i])
		}
	}
	var cachedUserIds []string
	for i := 0; i < len(cachedUsers); i++ {
		cachedUserIds = append(cachedUserIds, cachedUsers[i].UserId)
	}
	lg.Debug().String("cachedUsers", strings.Join(cachedUserIds, ",")).Fire()
	lg.Debug().String("uncachedUsers", strings.Join(uncachedUsers, ",")).Fire()
	lg.Debug().String("notExistsUsers", strings.Join(notExistsUsers, ",")).Fire()
	return cachedUsers, uncachedUsers, notExistsUsers, nil
}

func DescribeUsers(ctx context.Context, req *pbrequest.DescribeUsers) ([]*pbmodel.User, int64, error) {
	lg := glog.FromContext(ctx)
	if req.ReqSource == "" {
		req.ReqSource = executor.AccountExecutor.GetConf().Source
	}
	cachedUsers, uncachedUsers, notExistsUsers, err := getUsers(ctx, req.Users, req.ReqSource)
	if err != nil {
		return nil, 0, err
	}
	var users []source.User
	var totalCount int64
	var userSet []*pbmodel.User
	if len(uncachedUsers) != 0 {
		users, totalCount, err = source.SelectSource(req.ReqSource, options.Config, ctx).DescribeUsers(uncachedUsers, int(req.Limit)-len(cachedUsers), int(req.Offset))
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
		lg.Warn().String("Users not Exists", strings.Join(notExistsUsers, ",")).Fire()
		for _, userID := range notExistsUsers {
			cache.CacheNotExistUser(ctx, userID, req.ReqSource)
		}
	}

	totalCount += int64(len(cachedUsers))
	userSet = append(userSet, cachedUsers...)
	if err := cache.CacheUsers(ctx, userSet, req.ReqSource); err != nil {
		return nil, 0, err
	}

	return userSet, totalCount, err

}

func CreateUser(ctx context.Context, req *pbrequest.CreateUser) (*pbmodel.User, error) {
	id, err := options.IdGeneratorUser.Take()
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
	lg := glog.FromContext(ctx)
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
	if ignoreError := cache.DelUser(ctx, req.UserId, false); ignoreError != nil {
		lg.Warn().String("delete user cache error", ignoreError.Error()).Fire()
	}
	return userInfo.ToUserReply(), nil
}

func DeleteUser(ctx context.Context, req *pbrequest.DeleteUser) error {
	lg := glog.FromContext(ctx)
	err := gormwrap.ExecuteFuncWithTxn(ctx, executor.AccountExecutor.Db, func(tx *gorm.DB) error {
		var xErr error

		if xErr = executor.AccountExecutor.DeleteUser(tx, req.UserId); xErr != nil {
			return xErr
		}
		return xErr
	})
	if ignoreError := cache.DelUser(ctx, req.UserId, true); ignoreError != nil {
		lg.Warn().String("delete user cache error", ignoreError.Error()).Fire()
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
