package source

import (
	"context"

	"github.com/DataWorkbench/account/executor"
	"github.com/DataWorkbench/common/qerror"
)

type Local struct {
	ctx context.Context
}

func (l *Local) GetSecretAccessKey(accessKeyID string) (*executor.AccessKey, error) {
	accessKeyIDs := []string{accessKeyID}
	db := executor.AccountExecutor.Db.WithContext(l.ctx)
	accessKeys, err := executor.AccountExecutor.ListAccessKeys(db, accessKeyIDs, "", 1, 0)
	if err != nil {
		return nil, err
	}
	if len(accessKeys) != 1 {
		return nil, qerror.ResourceNotExists
	}
	return accessKeys[0], nil
}

func (l *Local) DescribeUsers(userIDs []string, limit int, offset int) ([]User, int64, error) {
	db := executor.AccountExecutor.Db.WithContext(l.ctx)
	users, err := executor.AccountExecutor.ListUsers(db, userIDs, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	totalCount, err := executor.AccountExecutor.CountUsers(db, userIDs)
	if err != nil {
		return nil, 0, err
	}
	iUsers := make([]User, len(users))
	for i, u := range users {
		iUsers[i] = u
	}
	return iUsers, totalCount, nil
}
