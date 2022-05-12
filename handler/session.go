package handler

import (
	"context"
	"github.com/DataWorkbench/account/executor"
	"github.com/DataWorkbench/common/constants"
	"github.com/DataWorkbench/common/gormwrap"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/common/utils/password"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
	"gorm.io/gorm"
	"time"
)

func CheckSession(ctx context.Context, req *pbrequest.CheckSession) (*pbresponse.CheckSession, error) {
	accessKey, err := cache.GetSession(req.Session)
	if err != nil {
		return nil, err
	}
	if accessKey != nil {
		var users []*pbmodel.User
		users, _, err = DescribeUsers(ctx, &pbrequest.DescribeUsers{
			Users: []string{accessKey.Owner},
			ReqSource: constants.LocalSource,
		})
		if err != nil {
			return nil, err
		}

		return &pbresponse.CheckSession{
			UserId:          accessKey.Owner,
			AccessKeyId:     accessKey.AccessKeyID,
			SecretAccessKey: accessKey.SecretAccessKey,
			Role: users[0].Role,
			Privilege: users[0].Privilege,
		}, nil
	} else {
		return nil, qerror.InvalidSession.Format(req.Session)
	}
}

func CreateSession(ctx context.Context, req *pbrequest.CreateSession) (*pbresponse.CreateSession, error) {
	var user *executor.User
	var accessKey *executor.AccessKey
	err := gormwrap.ExecuteFuncWithTxn(ctx, executor.AccountExecutor.Db, func(tx *gorm.DB) error {
		var xErr error
		if user, xErr = executor.AccountExecutor.GetUserByName(tx, req.UserName); xErr != nil {
			return xErr
		}
		if user == nil {
			xErr = qerror.UserNotExists.Format(req.UserName)
			return xErr
		}
		if !req.IgnorePassword && !password.Check(req.Password, user.Password) {
			logger.Warn().String(req.Password, user.Password).Fire()
			xErr = qerror.UserNameOrPasswordError
			return xErr
		}
		if accessKey, xErr = executor.AccountExecutor.GetAccessKeyByOwner(tx, user.UserID); xErr != nil {
			if xErr == gorm.ErrRecordNotFound {
				accessKeyId, secretKey := password.RandomGenerateAccessKey()
				accessKey = &executor.AccessKey{
					AccessKeyID:     accessKeyId,
					SecretAccessKey: secretKey,
					Owner:           user.UserID,
					Status:          constants.AccessKeyStatusEnable,
					CreateTime:      time.Now().Unix(),
					StatusTime:      time.Now().Unix(),
				}
				if xErr = executor.AccountExecutor.CreateAccessKey(tx, accessKey); xErr != nil {
					return xErr
				}
			} else {
				return xErr
			}
		}
		return xErr
	})
	if err != nil {
		return nil, err
	}
	session := password.GenerateSession()
	if err = cache.CacheSession(accessKey, session, user.UserID); err != nil {
		return nil, err
	}
	return &pbresponse.CreateSession{Session: session}, nil
}
