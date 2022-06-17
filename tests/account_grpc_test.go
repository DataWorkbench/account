package tests

import (
	"context"
	"testing"

	"github.com/DataWorkbench/common/constants"

	"github.com/DataWorkbench/glog"
	"github.com/DataWorkbench/gproto/xgo/service/pbsvcaccount"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
	"github.com/stretchr/testify/require"

	"github.com/DataWorkbench/common/grpcwrap"
	"github.com/DataWorkbench/common/gtrace"
	"github.com/DataWorkbench/common/utils/idgenerator"
)

//func sortReqQueryString(query string) string {
//	queryList := strings.Split(query, "&")
//	sort.Strings(queryList)
//	return strings.Join(queryList, "&")
//}

//func TestSignatureLocalWithNotExistKey(t *testing.T) {
//	address := "127.0.0.1:9110"
//	lp := glog.NewDefault()
//	ctx := glog.WithContext(context.Background(), lp)
//	conn, err := grpcwrap.NewConn(ctx, &grpcwrap.ClientConfig{
//		Address: address,
//	})
//
//	require.Nil(t, err, "%+v", err)
//
//	client := pbsvcaccount.NewAccountClient(conn)
//	logger := glog.NewDefault()
//
//	worker := idgenerator.New("")
//	reqId, _ := worker.Take()
//
//	ln := logger.Clone()
//	ln.WithFields().AddString("rid", reqId)
//	ctx = glog.WithContext(ctx, ln)
//	ctx = gtrace.ContextWithId(ctx, reqId)
//	_, err = client.ValidateRequestSignature(ctx, &pbrequest.ValidateRequestSignature{
//		ReqMethod:      "GET",
//		ReqPath:        "/iaas",
//		ReqQueryString: "access_key_id=NOTEXISTKEY&access_keys.1=NOTEXISTKEY&action=DescribeAccessKeys&limit=20&offset=0&signature_method=HmacSHA256&signature_version=1&time_stamp=2021-07-07T14%3A35%3A05Z&verbose=0",
//		ReqBody:        "",
//		ReqSignature:   "HhscGZcG6xgEME6sGwjbfb9KodcvF4WuJKRP1W3Z16Q%3D",
//		ReqAccessKeyId: "NOTEXISTKEY",
//		ReqSource:      constants.LocalSource,
//	})
//	if err != nil {
//		fmt.Printf("error message: %s", err.Error())
//	}
//	require.NotNilf(t, err, "%+v", err)
//}
//
//func TestSignatureQingcloudWithNotExistKey(t *testing.T) {
//	address := "127.0.0.1:9110"
//	lp := glog.NewDefault()
//	ctx := glog.WithContext(context.Background(), lp)
//	conn, err := grpcwrap.NewConn(ctx, &grpcwrap.ClientConfig{
//		Address: address,
//	})
//
//	require.Nil(t, err, "%+v", err)
//
//	client := pbsvcaccount.NewAccountClient(conn)
//	logger := glog.NewDefault()
//
//	worker := idgenerator.New("")
//	reqId, _ := worker.Take()
//
//	ln := logger.Clone()
//	ln.WithFields().AddString("rid", reqId)
//
//	ctx = glog.WithContext(ctx, ln)
//	ctx = gtrace.ContextWithId(ctx, reqId)
//	_, err = client.ValidateRequestSignature(ctx, &pbrequest.ValidateRequestSignature{
//		ReqMethod:      "GET",
//		ReqPath:        "/staging/v1/workspace/",
//		ReqQueryString: "access_key_id=NOTEXISTKEY&limit=10&offset=0&owner=usr-CVIshpN1&service=bigdata&signature_method=HmacSHA256&signature_version=1&time_stamp=2021-08-18T08%3A55%3A46Z&timestamp=2021-08-18T08%3A55%3A46Z&version=1",
//		ReqBody:        "null",
//		ReqSignature:   "xwd%2BbqeqwQkMrBVaGbQ1yBdm0F6Y198vf1dHnUp0Uj4%3D",
//		ReqAccessKeyId: "NOTEXISTKEY",
//		ReqSource:      constants.QingcloudSource,
//	})
//	if err != nil {
//		fmt.Printf("error message: %s", err.Error())
//	}
//	require.NotNilf(t, err, "%+v", err)
//
//}
//
//func TestSignatureWithDefaultSource(t *testing.T) {
//	address := "127.0.0.1:9110"
//	lp := glog.NewDefault()
//	ctx := glog.WithContext(context.Background(), lp)
//	conn, err := grpcwrap.NewConn(ctx, &grpcwrap.ClientConfig{
//		Address: address,
//	})
//
//	require.Nil(t, err, "%+v", err)
//
//	client := pbsvcaccount.NewAccountClient(conn)
//	logger := glog.NewDefault()
//
//	worker := idgenerator.New("")
//	reqId, _ := worker.Take()
//
//	ln := logger.Clone()
//	ln.WithFields().AddString("rid", reqId)
//	ctx = glog.WithContext(ctx, ln)
//	ctx = gtrace.ContextWithId(ctx, reqId)
//
//	_, err = client.ValidateRequestSignature(ctx, &pbrequest.ValidateRequestSignature{
//		ReqMethod:      "GET",
//		ReqPath:        "/staging/v1/workspace/",
//		ReqQueryString: sortReqQueryString("signature_method=HmacSHA256&reverse=True&service=bigdata&timestamp=2021-09-18T07%3A56%3A50Z&signature_version=1&access_key_id=ZMEVDSBCKAVSLHJECTUU&sort_by=created&version=1&limit=10&offset=0&owner=usr-CVIshpN1&time_stamp=2021-09-18T07%3A56%3A50Z"),
//		ReqBody:        "null",
//		ReqSignature:   "UIXVtLYN2MioTSzxdqranwqB%2BO4VHWZ7yFIiSr%2FGNqQ%3D",
//		ReqAccessKeyId: "ZMEVDSBCKAVSLHJECTUU",
//	})
//	require.Nil(t, err, "%+v", err)
//
//	_, err = client.ValidateRequestSignature(ctx, &pbrequest.ValidateRequestSignature{
//		ReqMethod:      "GET",
//		ReqPath:        "/api/region/",
//		ReqQueryString: sortReqQueryString("signature_method=HmacSHA256&service=bigdata&timestamp=2021-09-18T07%3A56%3A50Z&signature_version=1&version=1&access_key_id=ZMEVDSBCKAVSLHJECTUU&owner=usr-CVIshpN1&time_stamp=2021-09-18T07%3A56%3A50Z"),
//		ReqBody:        "null",
//		ReqSignature:   "5AGQsTinQLwy%2BYqilXVDmFlMjaAAEPsV3f1PP7uGepo%3D",
//		ReqAccessKeyId: "ZMEVDSBCKAVSLHJECTUU",
//	})
//	require.Nil(t, err, "%+v", err)
//}

//func TestSignatureLocal(t *testing.T) {
//	address := "127.0.0.1:9110"
//	lp := glog.NewDefault()
//	ctx := glog.WithContext(context.Background(), lp)
//	conn, err := grpcwrap.NewConn(ctx, &grpcwrap.ClientConfig{
//		Address: address,
//	})
//
//	require.Nil(t, err, "%+v", err)
//
//	client := pbsvcaccount.NewAccountClient(conn)
//	logger := glog.NewDefault()
//
//	worker := idgenerator.New("")
//	reqId, _ := worker.Take()
//
//	ln := logger.Clone()
//	ln.WithFields().AddString("rid", reqId)
//	ctx = glog.WithContext(ctx, ln)
//	ctx = gtrace.ContextWithId(ctx, reqId)
//	_, err = client.ValidateRequestSignature(ctx, &pbrequest.ValidateRequestSignature{
//		ReqMethod:      "GET",
//		ReqPath:        "/iaas",
//		ReqQueryString: "access_key_id=IKSMDBWVIECPIVNDYZAB&access_keys.1=IKSMDBWVIECPIVNDYZAB&action=DescribeAccessKeys&limit=20&offset=0&signature_method=HmacSHA256&signature_version=1&time_stamp=2021-07-07T14%3A35%3A05Z&verbose=0",
//		ReqBody:        "",
//		ReqSignature:   "HhscGZcG6xgEME6sGwjbfb9KodcvF4WuJKRP1W3Z16Q%3D",
//		ReqAccessKeyId: "IKSMDBWVIECPIVNDYZAB",
//		ReqSource:      constants.LocalSource,
//	})
//	if err != nil {
//		fmt.Printf("error message: %s", err.Error())
//	}
//	require.Nil(t, err, "%+v", err)
//}

//func TestSignatureQingcloud(t *testing.T) {
//	address := "127.0.0.1:9110"
//	lp := glog.NewDefault()
//	ctx := glog.WithContext(context.Background(), lp)
//	conn, err := grpcwrap.NewConn(ctx, &grpcwrap.ClientConfig{
//		Address: address,
//	})
//
//	require.Nil(t, err, "%+v", err)
//
//	client := pbsvcaccount.NewAccountClient(conn)
//	logger := glog.NewDefault()
//
//	worker := idgenerator.New("")
//	reqId, _ := worker.Take()
//
//	ln := logger.Clone()
//	ln.WithFields().AddString("rid", reqId)
//
//	_, err = client.ValidateRequestSignature(ctx, &pbrequest.ValidateRequestSignature{
//		ReqMethod:      "GET",
//		ReqPath:        "/staging/v1/workspace/",
//		ReqQueryString: sortReqQueryString("signature_method=HmacSHA256&reverse=True&service=bigdata&timestamp=2021-09-18T07%3A56%3A50Z&signature_version=1&access_key_id=ZMEVDSBCKAVSLHJECTUU&sort_by=created&version=1&limit=10&offset=0&owner=usr-CVIshpN1&time_stamp=2021-09-18T07%3A56%3A50Z"),
//		ReqBody:        "null",
//		ReqSignature:   "UIXVtLYN2MioTSzxdqranwqB%2BO4VHWZ7yFIiSr%2FGNqQ%3D",
//		ReqAccessKeyId: "ZMEVDSBCKAVSLHJECTUU",
//		ReqSource:      constants.QingcloudSource,
//	})
//	require.Nil(t, err, "%+v", err)
//
//	_, err = client.ValidateRequestSignature(ctx, &pbrequest.ValidateRequestSignature{
//		ReqMethod:      "GET",
//		ReqPath:        "/api/region/",
//		ReqQueryString: sortReqQueryString("signature_method=HmacSHA256&service=bigdata&timestamp=2021-09-18T07%3A56%3A50Z&signature_version=1&version=1&access_key_id=ZMEVDSBCKAVSLHJECTUU&owner=usr-CVIshpN1&time_stamp=2021-09-18T07%3A56%3A50Z"),
//		ReqBody:        "null",
//		ReqSignature:   "5AGQsTinQLwy%2BYqilXVDmFlMjaAAEPsV3f1PP7uGepo%3D",
//		ReqAccessKeyId: "ZMEVDSBCKAVSLHJECTUU",
//		ReqSource:      constants.QingcloudSource,
//	})
//	require.Nil(t, err, "%+v", err)
//
//}

func TestUsersWithDefaultSource(t *testing.T) {
	address := "127.0.0.1:9110"
	lp := glog.NewDefault()
	ctx := glog.WithContext(context.Background(), lp)
	conn, err := grpcwrap.NewConn(ctx, &grpcwrap.ClientConfig{
		Address: address,
	})

	require.Nil(t, err, "%+v", err)

	client := pbsvcaccount.NewAccountClient(conn)
	logger := glog.NewDefault()

	worker := idgenerator.New("")
	reqId, _ := worker.Take()

	ln := logger.Clone()
	ln.WithFields().AddString("rid", reqId)
	ctx = glog.WithContext(ctx, ln)
	ctx = gtrace.ContextWithId(ctx, reqId)
	result, err := client.DescribeUsers(ctx, &pbrequest.DescribeUsers{
		Users:  []string{"usr-CVIshpN1", "usr-notexist1", "usr-notexist2"},
		Offset: 0,
		Limit:  20,
	})

	if result != nil {
		ln.Info().String("receive Users ", result.String()).Fire()
	}

	require.Nil(t, err, "%+v", err)
}

func TestUsersLocal(t *testing.T) {
	address := "127.0.0.1:9110"
	lp := glog.NewDefault()
	ctx := glog.WithContext(context.Background(), lp)
	conn, err := grpcwrap.NewConn(ctx, &grpcwrap.ClientConfig{
		Address: address,
	})

	require.Nil(t, err, "%+v", err)

	client := pbsvcaccount.NewAccountClient(conn)
	logger := glog.NewDefault()

	worker := idgenerator.New("")
	reqId, _ := worker.Take()

	ln := logger.Clone()
	ln.WithFields().AddString("rid", reqId)
	ctx = glog.WithContext(ctx, ln)
	ctx = gtrace.ContextWithId(ctx, reqId)
	result, err := client.DescribeUsers(ctx, &pbrequest.DescribeUsers{
		Users:     []string{"usr-iDMTmjGs", "usr-notexist1", "usr-notexist2"},
		Offset:    0,
		Limit:     20,
		ReqSource: constants.LocalSource,
	})

	if result != nil {
		ln.Info().String("receive Users ", result.String()).Fire()
	}

	require.Nil(t, err, "%+v", err)
}

func TestUsersQingcloud(t *testing.T) {
	address := "127.0.0.1:9110"
	lp := glog.NewDefault()
	ctx := glog.WithContext(context.Background(), lp)
	conn, err := grpcwrap.NewConn(ctx, &grpcwrap.ClientConfig{
		Address: address,
	})

	require.Nil(t, err, "%+v", err)

	client := pbsvcaccount.NewAccountClient(conn)
	logger := glog.NewDefault()

	worker := idgenerator.New("")
	reqId, _ := worker.Take()

	ln := logger.Clone()
	ln.WithFields().AddString("rid", reqId)
	ctx = glog.WithContext(ctx, ln)
	ctx = gtrace.ContextWithId(ctx, reqId)
	result, err := client.DescribeUsers(ctx, &pbrequest.DescribeUsers{
		Users:     []string{"usr-CVIshpN1", "usr-notexist1", "usr-notexist2"},
		Offset:    0,
		Limit:     3,
		ReqSource: constants.QingcloudSource,
	})

	if result != nil {
		lp.Info().String("receive Users ", result.String()).Fire()
	}

	require.Nil(t, err, "%+v", err)
}

func TestCreateUser(t *testing.T) {
	address := "127.0.0.1:9110"
	lp := glog.NewDefault()
	ctx := glog.WithContext(context.Background(), lp)
	conn, err := grpcwrap.NewConn(ctx, &grpcwrap.ClientConfig{
		Address: address,
	})

	require.Nil(t, err, "%+v", err)

	client := pbsvcaccount.NewAccountClient(conn)
	logger := glog.NewDefault()

	worker := idgenerator.New("")
	reqId, _ := worker.Take()

	ln := logger.Clone()
	ln.WithFields().AddString("rid", reqId)
	ctx = glog.WithContext(ctx, ln)
	ctx = gtrace.ContextWithId(ctx, reqId)
	userName := "ryansun"
	lang := "cn"
	email := "ryansun@yunify.com"
	phone := "15802789195"
	currency := "CNY"
	createUser, err := client.CreateUser(ctx, &pbrequest.CreateUser{
		UserName: userName,
		Password: "zhu88jie",
		Lang:     lang,
		Email:    email,
		Phone:    phone,
		Currency: currency,
	})
	require.Nil(t, err, "%+v", err)
	require.NotNil(t, createUser)
	require.NotNil(t, createUser.User)
	require.NotEmpty(t, createUser.User.UserId)
	require.Equal(t, userName, createUser.User.UserName)
	require.Equal(t, lang, createUser.User.Lang)
	require.Equal(t, email, createUser.User.Email)
	require.Equal(t, phone, createUser.User.Phone)
	//require.Equal(t, currency, createUser.User.Currency)
	require.Equal(t, constants.UserStatusActive, createUser.User.Status)

	describeUser, err := client.DescribeUsers(ctx, &pbrequest.DescribeUsers{
		Users:     []string{createUser.User.UserId},
		ReqSource: constants.LocalSource})
	require.Nil(t, err, "%+v", err)
	require.NotNil(t, describeUser)
	require.NotEmpty(t, describeUser.UserSet)
	require.NotEmpty(t, describeUser.UserSet[0].UserId)
	require.Equal(t, userName, describeUser.UserSet[0].UserName)
	require.Equal(t, lang, describeUser.UserSet[0].Lang)
	require.Equal(t, email, describeUser.UserSet[0].Email)
	require.Equal(t, phone, describeUser.UserSet[0].Phone)
	//require.Equal(t, currency, describeUser.UserSet[0].Currency)
	require.Equal(t, constants.UserStatusActive, describeUser.UserSet[0].Status)

	deleteUser, err := client.DeleteUser(ctx, &pbrequest.DeleteUser{UserId: createUser.User.UserId})
	require.Nil(t, err, "%+v", err)
	require.NotNil(t, deleteUser)

}

func TestUpdateUser(t *testing.T) {
	address := "127.0.0.1:9110"
	lp := glog.NewDefault()
	ctx := glog.WithContext(context.Background(), lp)
	conn, err := grpcwrap.NewConn(ctx, &grpcwrap.ClientConfig{
		Address: address,
	})

	require.Nil(t, err, "%+v", err)

	client := pbsvcaccount.NewAccountClient(conn)
	logger := glog.NewDefault()

	worker := idgenerator.New("")
	reqId, _ := worker.Take()

	ln := logger.Clone()
	ln.WithFields().AddString("rid", reqId)
	ctx = glog.WithContext(ctx, ln)
	ctx = gtrace.ContextWithId(ctx, reqId)
	userName1 := "ryansun1"
	email1 := "ryansun1@yunify.com"
	phone1 := "100861"
	userName2 := "ryansun2"
	email2 := "ryansun2@yunify.com"
	phone2 := "100862"
	lang := "cn"
	currency := "CNY"

	createUser1, err := client.CreateUser(ctx, &pbrequest.CreateUser{
		UserName: userName1,
		Password: "zhu88jie",
		Lang:     lang,
		Email:    email1,
		Phone:    phone1,
		Currency: currency,
	})
	require.Nil(t, err, "%+v", err)
	require.NotNil(t, createUser1)
	require.NotNil(t, createUser1.User)

	createUser2, err := client.CreateUser(ctx, &pbrequest.CreateUser{
		UserName: userName2,
		Password: "zhu88jie",
		Lang:     lang,
		Email:    email2,
		Phone:    phone2,
		Currency: currency,
	})
	require.Nil(t, err, "%+v", err)
	require.NotNil(t, createUser2)
	require.NotNil(t, createUser2.User)

	updateUser1, err := client.UpdateUser(ctx, &pbrequest.UpdateUser{
		UserId:   createUser1.User.UserId,
		UserName: createUser2.User.UserName,
	})
	require.NotNil(t, err, "%+v", err)
	require.Nil(t, updateUser1)

	newUserName := "newryansun"
	updateUser2, err := client.UpdateUser(ctx, &pbrequest.UpdateUser{
		UserId:   createUser1.User.UserId,
		UserName: newUserName,
	})
	require.Nil(t, err, "%+v", err)
	require.NotNil(t, updateUser2)
	require.NotNil(t, updateUser2.User)
	require.Equal(t, newUserName, updateUser2.User.UserName)

	deleteUser1, err := client.DeleteUser(ctx, &pbrequest.DeleteUser{UserId: createUser1.User.UserId})
	require.Nil(t, err, "%+v", err)
	require.NotNil(t, deleteUser1)
	deleteUser2, err := client.DeleteUser(ctx, &pbrequest.DeleteUser{UserId: createUser2.User.UserId})
	require.Nil(t, err, "%+v", err)
	require.NotNil(t, deleteUser2)
}

func TestDeleteUser(t *testing.T) {
	address := "127.0.0.1:9110"
	lp := glog.NewDefault()
	ctx := glog.WithContext(context.Background(), lp)
	conn, err := grpcwrap.NewConn(ctx, &grpcwrap.ClientConfig{
		Address: address,
	})

	require.Nil(t, err, "%+v", err)

	client := pbsvcaccount.NewAccountClient(conn)
	logger := glog.NewDefault()

	worker := idgenerator.New("")
	reqId, _ := worker.Take()

	ln := logger.Clone()
	ln.WithFields().AddString("rid", reqId)
	ctx = glog.WithContext(ctx, ln)
	ctx = gtrace.ContextWithId(ctx, reqId)
	userName := "ryansun"
	lang := "cn"
	email := "ryansun@yunify.com"
	phone := "15802789195"
	currency := "CNY"
	createUser, err := client.CreateUser(ctx, &pbrequest.CreateUser{
		UserName: userName,
		Password: "zhu88jie",
		Lang:     lang,
		Email:    email,
		Phone:    phone,
		Currency: currency,
	})
	require.Nil(t, err, "%+v", err)
	require.NotNil(t, createUser)
	require.NotNil(t, createUser.User)

	deleteUser, err := client.DeleteUser(ctx, &pbrequest.DeleteUser{UserId: createUser.User.UserId})
	require.Nil(t, err, "%+v", err)
	require.NotNil(t, deleteUser)

	describeUser, err := client.DescribeUsers(ctx, &pbrequest.DescribeUsers{
		Users:     []string{createUser.User.UserId},
		ReqSource: constants.LocalSource})
	require.Nil(t, err, "%+v", err)
	require.NotNil(t, describeUser)
	require.NotEmpty(t, describeUser.UserSet)
	require.Equal(t, constants.UserStatusDelete, describeUser.UserSet[0].Status)
}

func TestCreateSession(t *testing.T) {
	address := "127.0.0.1:9110"
	lp := glog.NewDefault()
	ctx := glog.WithContext(context.Background(), lp)
	conn, err := grpcwrap.NewConn(ctx, &grpcwrap.ClientConfig{
		Address: address,
	})

	require.Nil(t, err, "%+v", err)

	client := pbsvcaccount.NewAccountClient(conn)
	logger := glog.NewDefault()

	worker := idgenerator.New("")
	reqId, _ := worker.Take()

	ln := logger.Clone()
	ln.WithFields().AddString("rid", reqId)
	ctx = glog.WithContext(ctx, ln)
	ctx = gtrace.ContextWithId(ctx, reqId)
	userName := "ryansun"
	lang := "cn"
	email := "ryansun@yunify.com"
	phone := "15802789195"
	currency := "CNY"
	password := "zhu88jie"
	createUser, err := client.CreateUser(ctx, &pbrequest.CreateUser{
		UserName: userName,
		Password: password,
		Lang:     lang,
		Email:    email,
		Phone:    phone,
		Currency: currency,
	})
	require.Nil(t, err, "%+v", err)
	require.NotNil(t, createUser)
	require.NotNil(t, createUser.User)

	createSession, err := client.CreateSession(ctx, &pbrequest.CreateSession{UserName: userName, Password: password})
	require.Nil(t, err, "%+v", err)
	require.NotNil(t, createSession)
	require.NotNil(t, createSession.Session)

	checkSession, err := client.CheckSession(ctx, &pbrequest.CheckSession{Session: createSession.Session})
	require.Nil(t, err, "%+v", err)
	require.NotNil(t, checkSession)
	require.Equal(t, createUser.User.UserId, checkSession.UserId)
	accessKey, err := client.DescribeAccessKey(ctx, &pbrequest.DescribeAccessKey{
		AccessKeyId: checkSession.AccessKeyId})
	require.Nil(t, err, "%+v", err)
	require.NotNil(t, accessKey)
	require.Equal(t, accessKey.Owner, checkSession.UserId)
	require.Equal(t, accessKey.SecretAccessKey, checkSession.SecretAccessKey)

	deleteUser, err := client.DeleteUser(ctx, &pbrequest.DeleteUser{UserId: createUser.User.UserId})
	require.Nil(t, err, "%+v", err)
	require.NotNil(t, deleteUser)

	checkSession, err = client.CheckSession(ctx, &pbrequest.CheckSession{Session: createSession.Session})
	require.NotNil(t, err, "%+v", err)
	require.Nil(t, checkSession)
}
