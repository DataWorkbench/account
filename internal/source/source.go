package source

import (
	"context"

	"github.com/DataWorkbench/account/config"
	"github.com/DataWorkbench/account/executor"
	"github.com/DataWorkbench/common/constants"
	"github.com/DataWorkbench/gproto/pkg/accountpb"
)

type User interface {
	ToUserReply() *accountpb.User
}

type Source interface {
	GetSecretAccessKey(accessKeyID string) (*executor.AccessKey, error)
	DescribeUsers(userIDs []string, limit int, offset int) ([]User, int64, error)
}

func SelectSource(source string, conf *config.Config, ctx context.Context) Source {
	switch source {
	case constants.QingcloudSource:
		return &Qingcloud{
			qingcloudConfig: executor.AccountExecutor.GetConf().Qingcloud,
		}
	case constants.LocalSource:
		return &Local{
			ctx: ctx,
		}
	default:
		return &Local{
			ctx: ctx,
		}
	}
}
