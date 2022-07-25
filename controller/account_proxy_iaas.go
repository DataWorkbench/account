package controller

import (
	"context"

	"github.com/DataWorkbench/account/options"
	"github.com/DataWorkbench/common/lib/iaas"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/gproto/xgo/service/pbsvcaccount"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel/pbiaas"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
)

type AccountProxyIaaS struct {
	pbsvcaccount.UnimplementedAccountProxyServer
}

func (x *AccountProxyIaaS) ListUsersByProxy(ctx context.Context, req *pbrequest.ListUsersByProxy) (
	*pbresponse.ListUsersByProxy, error) {
	if len(req.UserIds) == 0 {
		// IaaS will returns all users if id list of id is empty.
		// Thus wo should return empty in proxy.
		return &pbresponse.ListUsersByProxy{Infos: nil, Total: 0, HasMore: false}, nil
	}
	output, err := options.IaaSClient.DescribeUsers(ctx, &iaas.DescribeUsersInput{
		Users:    req.UserIds,
		Limit:    int(req.Limit),
		Offset:   int(req.Offset),
		Status:   "",
		Email:    "",
		Phone:    "",
		RootUser: req.RootUser,
	})
	if err != nil {
		return nil, err
	}

	infos := x.iaasUsersToUsers(output.UserSet)
	hasMore := len(infos) >= int(req.Limit)

	reply := &pbresponse.ListUsersByProxy{
		Infos:   infos,
		Total:   int64(output.TotalCount),
		HasMore: hasMore,
	}
	return reply, nil
}
func (x *AccountProxyIaaS) DescribeAccessKeyByProxy(ctx context.Context, req *pbrequest.DescribeAccessKeyByProxy) (
	*pbresponse.DescribeAccessKeyByProxy, error) {
	keySet, err := options.IaaSClient.DescribeAccessKeysById(ctx, req.AccessKeyId)
	if err != nil {
		if err == iaas.ErrAccessKeyNotExists {
			err = qerror.AccessKeyNotExists.Format(req.AccessKeyId)
		}
		return nil, err
	}

	reply := &pbresponse.DescribeAccessKeyByProxy{
		KeySet: x.iaasAccessKeyToAccessKey(keySet),
	}
	return reply, nil
}

func (x *AccountProxyIaaS) ListNotificationsByProxy(ctx context.Context, req *pbrequest.ListNotificationsByProxy) (
	*pbresponse.ListNotificationsByProxy, error) {
	output, err := options.IaaSClient.DescribeNotificationLists(ctx, req.UserId, req.NfIds, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	infos := x.iaasNotificationListsToNotifications(output.NotificationListSet)
	hasMore := len(infos) >= int(req.Limit)

	reply := &pbresponse.ListNotificationsByProxy{
		Infos:   infos,
		Total:   int64(output.TotalCount),
		HasMore: hasMore,
	}
	return reply, nil
}

func (x *AccountProxyIaaS) iaasUserToUser(iaasUser *pbiaas.User) *pbmodel.User {
	role := pbmodel.User_User
	if iaasUser.Role == "global_admin" {
		role = pbmodel.User_Admin
	}

	u := &pbmodel.User{
		UserId:   iaasUser.UserId,
		Name:     iaasUser.UserName,
		Email:    iaasUser.Email,
		Role:     role,
		Status:   pbmodel.User_active,
		Password: "******",
		Created:  iaas.TimeStringToTimestampSecond(iaasUser.CreateTime),
		Updated:  iaas.TimeStringToTimestampSecond(iaasUser.StatusTime),
	}
	return u
}

func (x *AccountProxyIaaS) iaasUsersToUsers(iaasUsers []*pbiaas.User) []*pbmodel.User {
	if iaasUsers == nil {
		return nil
	}
	users := make([]*pbmodel.User, len(iaasUsers))
	for i := 0; i < len(iaasUsers); i++ {
		users[i] = x.iaasUserToUser(iaasUsers[i])
	}
	return users
}

func (x *AccountProxyIaaS) iaasAccessKeyToAccessKey(iaasKey *pbiaas.AccessKey) *pbmodel.AccessKey {
	if iaasKey == nil {
		return nil
	}
	controller, ok := pbmodel.AccessKey_Controller_value[iaasKey.Controller]
	if !ok {
		controller = int32(pbmodel.AccessKey_self)
	}
	status, ok := pbmodel.AccessKey_Status_value[iaasKey.Status]
	if !ok {
		status = int32(pbmodel.AccessKey_active)
	}

	key := &pbmodel.AccessKey{
		AccessKeyId:     iaasKey.AccessKeyId,
		SecretAccessKey: iaasKey.SecretAccessKey,
		Name:            iaasKey.AccessKeyName,
		Owner:           iaasKey.Owner,
		Controller:      pbmodel.AccessKey_Controller(controller),
		Status:          pbmodel.AccessKey_Status(status),
		Description:     iaasKey.Description,
		IpWhiteList:     iaasKey.IpWhiteList,
		Created:         iaas.TimeStringToTimestampSecond(iaasKey.CreateTime),
		Updated:         iaas.TimeStringToTimestampSecond(iaasKey.StatusTime),
	}
	return key
}

func (x *AccountProxyIaaS) iaasNotificationListToNotification(iaasNFList *pbiaas.NotificationList) *pbmodel.Notification {
	var email string
	for _, item := range iaasNFList.Items {
		switch item.NotificationItemType {
		case "email":
			email = item.Content
		}
	}
	nf := &pbmodel.Notification{
		Owner:       iaasNFList.Owner,
		Id:          iaasNFList.NotificationListId,
		Name:        iaasNFList.NotificationListName,
		Description: "", // No description in iaas.
		Email:       email,
		Created:     iaas.TimeStringToTimestampSecond(iaasNFList.CreateTime),
		Updated:     iaas.TimeStringToTimestampSecond(iaasNFList.CreateTime), // No update time in iaas response.
	}
	return nf
}

func (x *AccountProxyIaaS) iaasNotificationListsToNotifications(iaasNFLists []*pbiaas.NotificationList) []*pbmodel.Notification {
	if iaasNFLists == nil {
		return nil
	}
	nfs := make([]*pbmodel.Notification, len(iaasNFLists))
	for i := 0; i < len(iaasNFLists); i++ {
		nfs[i] = x.iaasNotificationListToNotification(iaasNFLists[i])
	}
	return nfs
}
