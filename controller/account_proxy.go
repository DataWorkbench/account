package controller

import (
	"context"

	"github.com/DataWorkbench/account/options"
	"github.com/DataWorkbench/common/lib/iaas"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/gproto/xgo/service/pbsvcaccount"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
)

type AccountProxy struct {
	pbsvcaccount.UnimplementedAccountProxyServer
}

func (x *AccountProxy) ListUsersByProxy(ctx context.Context, req *pbrequest.ListUsersByProxy) (
	*pbresponse.ListUsersByProxy, error) {
	output, err := options.IaaSClient.DescribeUsers(ctx, &iaas.DescribeUsersInput{
		Users:  req.Users,
		Limit:  int(req.Limit),
		Offset: int(req.Offset),
		Status: "",
		Email:  "",
		Phone:  "",
	})
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.ListUsersByProxy{
		UserSet:    output.UserSet,
		TotalCount: int64(output.TotalCount),
	}
	return reply, nil
}
func (x *AccountProxy) DescribeAccessKeyByProxy(ctx context.Context, req *pbrequest.DescribeAccessKeyByProxy) (
	*pbresponse.DescribeAccessKeyByProxy, error) {
	keySet, err := options.IaaSClient.DescribeAccessKeysById(ctx, req.AccessKeyId)
	if err != nil {
		if err == iaas.ErrAccessKeyNotExists {
			err = qerror.AccessKeyNotExists.Format(req.AccessKeyId)
		}
		return nil, err
	}

	reply := &pbresponse.DescribeAccessKeyByProxy{
		KeySet: keySet,
	}
	return reply, nil
}

func (x *AccountProxy) ListNotificationsByProxy(ctx context.Context, req *pbrequest.ListNotificationsByProxy) (
	*pbresponse.ListNotificationsByProxy, error) {
	output, err := options.IaaSClient.DescribeNotificationLists(ctx, req.Owner, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.ListNotificationsByProxy{
		NotificationLists: nil,
		Total:             int64(output.TotalCount),
	}
	return reply, nil
}
