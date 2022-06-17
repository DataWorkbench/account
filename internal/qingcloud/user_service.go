package qingcloud

import (
	"fmt"
	"time"

	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/yunify/qingcloud-sdk-go/config"
	"github.com/yunify/qingcloud-sdk-go/request"
	"github.com/yunify/qingcloud-sdk-go/request/data"
)

var _ fmt.State
var _ time.Time

type User struct {
	UserID        string   `json:"user_id" name:"user_id"`
	UserName      string   `json:"user_name" name:"user_name"`
	Lang          string   `json:"lang" name:"lang"`
	Email         string   `json:"email" name:"email"`
	Phone         string   `json:"phone" name:"phone"`
	Status        string   `json:"status" name:"status"`
	Role          string   `json:"role" name:"role"`
	Currency      string   `json:"currency" name:"currency"`
	GravatarEmail string   `json:"gravatar_email" name:"gravatar_email"`
	Privilege     int32    `json:"privilege" name:"privilege"`
	Zones         []string `json:"zones" name:"zones"`
	Regions       []string `json:"regions" name:"regions"`
}

func (u *User) ToUserReply() *pbmodel.User {
	return &pbmodel.User{
		UserId:        u.UserID,
		UserName:      u.UserName,
		Lang:          u.Lang,
		Email:         u.Email,
		Phone:         u.Phone,
		Status:        u.Status,
		Role:          u.Role,
		GravatarEmail: u.GravatarEmail,
		Privilege:     u.Privilege,
		Zones:         u.Zones,
		Regions:       u.Regions,
	}
}

type UserService struct {
	Config     *config.Config
	Properties *UserServiceProperties
}

type UserServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *UserService) DescribeUsers(i *DescribeUsersInput) (*DescribeUsersOutput, error) {
	if i == nil {
		i = &DescribeUsersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeUsers",
		RequestMethod: "GET",
	}

	x := &DescribeUsersOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type DescribeUsersInput struct {
	Users  []*string `json:"users" name:"users" location:"params"`
	Limit  *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset *int      `json:"offset" name:"offset" default:"0" location:"params"`
}

func (v *DescribeUsersInput) Validate() error {

	return nil
}

type DescribeUsersOutput struct {
	Message    *string `json:"message" name:"message"`
	UserSet    []*User `json:"user_set" name:"user_set" location:"elements"`
	Action     *string `json:"action" name:"action" location:"elements"`
	RetCode    *int    `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int    `json:"total_count" name:"total_count" location:"elements"`
}
