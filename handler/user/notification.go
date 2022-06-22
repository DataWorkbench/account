package user

import (
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
	"gorm.io/gorm"
)

func ListNotifications(tx *gorm.DB, input *pbrequest.ListNotifications) (output *pbresponse.ListNotifications, err error) {
	panic("unrealized")
}
