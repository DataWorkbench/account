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
	accessKeys, err := executor.AccountExecutor.ListAccessKeys(l.ctx, accessKeyIDs, "", 1, 0)
	if err != nil {
		return nil, err
	}
	if len(accessKeys) != 1 {
		return nil, qerror.AccessKeyNotExists.Format(accessKeyID)
	}
	return accessKeys[0], nil
}

func (l *Local) DescribeUsers(userIDs []string, limit int, offset int) ([]User, int64, error) {
	users, err := executor.AccountExecutor.ListUsers(l.ctx, userIDs, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	totalCount, err := executor.AccountExecutor.CountUsers(l.ctx, userIDs)
	if err != nil {
		return nil, 0, err
	}
	iUsers := make([]User, len(users))
	for i, u := range users {
		iUsers[i] = u
	}
	return iUsers, totalCount, nil
}
