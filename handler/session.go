package handler

import (
	"context"
	"github.com/DataWorkbench/account/executor"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/common/utils/password"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
)

func CheckSession(ctx context.Context, req *pbrequest.CheckSession) (*pbresponse.CheckSession, error) {
	accessKey, err := cache.GetSession(req.Session)
	if err != nil {
		return nil, err
	}
	if accessKey != nil {

		return &pbresponse.CheckSession{
			UserId:          accessKey.Owner,
			AccessKeyId:     accessKey.AccessKeyID,
			SecretAccessKey: accessKey.SecretAccessKey,
		}, nil
	} else {
		return nil, qerror.InvalidSession.Format(req.Session)
	}
}

func CreateSession(ctx context.Context, req *pbrequest.CreateSession) (*pbresponse.CreateSession, error) {
	db := executor.AccountExecutor.Db.WithContext(ctx)
	user, err := executor.AccountExecutor.GetUserByName(db, req.UserName)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, qerror.UserNotExists.Format(req.UserName)
	}
	if !password.Check(req.Password, user.Password) {
		logger.Warn().String(req.Password, user.Password).Fire()
		return nil, qerror.UserNameOrPasswordError
	}
	accessKey, err := executor.AccountExecutor.GetAccessKeyByOwner(db, user.UserID)
	if err != nil {
		return nil, err
	}
	session := password.GenerateSession()
	if err = cache.CacheSession(accessKey, session, user.UserID); err != nil {
		return nil, err
	}
	return &pbresponse.CreateSession{Session: session}, nil
}
