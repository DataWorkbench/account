package controller

import (
	"context"

	"github.com/DataWorkbench/account/handler/user"
	"github.com/DataWorkbench/account/options"
	"github.com/DataWorkbench/gproto/xgo/service/pbsvcaccount"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
)

type AccountProxyLocal struct {
	pbsvcaccount.UnimplementedAccountProxyServer
}

func (x *AccountProxyLocal) ListUsersByProxy(ctx context.Context, req *pbrequest.ListUsersByProxy) (
	*pbresponse.ListUsersByProxy, error) {
	tx := options.DBConn.WithContext(ctx)

	output, err := user.ListUsers(tx, &pbrequest.ListUsers{
		Limit:   req.Limit,
		Offset:  req.Offset,
		UserIds: req.UserIds,
		Name:    req.Name,
	})
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.ListUsersByProxy{
		Infos:   output.Infos,
		Total:   output.Total,
		HasMore: output.HasMore,
	}

	return reply, nil
}
func (x *AccountProxyLocal) DescribeAccessKeyByProxy(ctx context.Context, req *pbrequest.DescribeAccessKeyByProxy) (
	*pbresponse.DescribeAccessKeyByProxy, error) {
	tx := options.DBConn.WithContext(ctx)
	keySet, err := user.DescribeAccessKeyByKeyId(tx, req.AccessKeyId)
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.DescribeAccessKeyByProxy{KeySet: keySet}
	return reply, nil
}
func (x *AccountProxyLocal) ListNotificationsByProxy(ctx context.Context, req *pbrequest.ListNotificationsByProxy) (
	*pbresponse.ListNotificationsByProxy, error) {
	tx := options.DBConn.WithContext(ctx)
	output, err := user.ListNotifications(tx, &pbrequest.ListNotifications{
		UserId: req.UserId,
		Limit:  req.Limit,
		Offset: req.Offset,
		NfIds:  req.NfIds,
	})
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.ListNotificationsByProxy{
		Infos:   output.Infos,
		Total:   output.Total,
		HasMore: output.HasMore,
	}
	return reply, nil
}
