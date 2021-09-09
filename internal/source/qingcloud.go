package source

import (
	"github.com/DataWorkbench/account/config"
	"github.com/DataWorkbench/account/executor"
	"github.com/DataWorkbench/account/internal/qingcloud"
	"github.com/DataWorkbench/common/constants"
	"github.com/DataWorkbench/common/qerror"
	qConfig "github.com/yunify/qingcloud-sdk-go/config"
	qService "github.com/yunify/qingcloud-sdk-go/service"
)

type Qingcloud struct {
	qingcloudConfig *config.QingcloudConfig
}

func (q *Qingcloud) GetSecretAccessKey(accessKeyID string) (*executor.AccessKey, error) {
	qCfg, err := qConfig.New(q.qingcloudConfig.AccessKeyID, q.qingcloudConfig.SecretAccessKey)
	if err != nil {
		return nil, err
	}
	qCfg.Host = q.qingcloudConfig.Host
	qCfg.Port = q.qingcloudConfig.Port
	qCfg.Protocol = q.qingcloudConfig.Protocol
	qCfg.URI = q.qingcloudConfig.Uri

	qSvc, err := qService.Init(qCfg)
	if err != nil {
		return nil, err
	}
	akSvc, err := qSvc.Accesskey("")
	if err != nil {
		return nil, err
	}
	resp, err := akSvc.DescribeAccessKeys(&qService.DescribeAccessKeysInput{
		AccessKeys: []*string{&accessKeyID},
	})
	if err != nil {
		return nil, err
	}
	if len(resp.AccessKeySet) != 1 || *(resp.AccessKeySet[0].Status) != constants.QingcloudAccessKeyStatusActive {
		return nil, qerror.ResourceNotExists
	}
	return &executor.AccessKey{
		AccessKeyID:     *resp.AccessKeySet[0].AccessKeyID,
		SecretAccessKey: *resp.AccessKeySet[0].SecretAccessKey,
		Owner:           *resp.AccessKeySet[0].Owner,
	}, nil
}

func (q *Qingcloud) DescribeUsers(userIDs []string, limit int, offset int) ([]User, int64, error) {
	qCfg, err := qConfig.New(q.qingcloudConfig.AccessKeyID, q.qingcloudConfig.SecretAccessKey)
	if err != nil {
		return nil, 0, err
	}
	qCfg.Host = q.qingcloudConfig.Host
	qCfg.Port = q.qingcloudConfig.Port
	qCfg.Protocol = q.qingcloudConfig.Protocol
	qCfg.URI = q.qingcloudConfig.Uri

	qSvc, err := qService.Init(qCfg)
	if err != nil {
		return nil, 0, err
	}
	var users []*string
	for i := 0; i < len(userIDs); i++ {
		users = append(users, &userIDs[i])
	}
	userSvc := &qingcloud.UserService{Config: qSvc.Config, Properties: &qingcloud.UserServiceProperties{Zone: &qCfg.Zone}}
	resp, err := userSvc.DescribeUsers(&qingcloud.DescribeUsersInput{
		Users:  users,
		Limit:  &limit,
		Offset: &offset,
	})
	if err != nil {
		return nil, 0, err
	}
	iUsers := make([]User, len(resp.UserSet))
	for i, u := range resp.UserSet {
		iUsers[i] = u
	}
	return iUsers, int64(*resp.TotalCount), nil
}
