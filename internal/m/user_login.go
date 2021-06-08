package m

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iegad/hydra/micro"
	"github.com/iegad/hydra/mod/home"
	"github.com/iegad/hydra/pb"
	"github.com/iegad/kraken/log"
	"github.com/iegad/kraken/utils"
	"github.com/iegad/sphinx/internal/com"
	"google.golang.org/protobuf/proto"
)

var (
	ErrAccount = errors.New("account is invalid")
	ErrVCode   = errors.New("vcode is invalid")
)

type UserLogin struct {
}

func (this_ *UserLogin) MID() pb.MessageID {
	return pb.MessageID_MID_UserLoginReq
}

func (this_ *UserLogin) Do(c *micro.User, in *pb.Package) error {
	utils.Assert(c != nil && in != nil, "userLogin.do params are invalid")

	var (
		req      = pb.NewUserLoginReq()
		err      = proto.Unmarshal(in.Data, req)
		user     *home.UserInfo
		dataList []*home.UserInfo
		where    = ""
		rsp      = pb.NewUserLoginRsp()
	)

	if err != nil {
		c.Kick()
		goto DO_EXIT
	}

	// 入参检查
	if len(req.Email) == 0 && len(req.PhoneNum) == 0 {
		rsp.Code = -1
		goto DO_EXIT
	}

	if len(req.VCode) == 0 {
		rsp.Code = -2
		goto DO_EXIT
	}

	// 查询数据库
	if len(req.Email) > 0 {
		where = fmt.Sprintf("F_EMAIL='%s'", req.Email)
	} else {
		where = fmt.Sprintf("F_PHONE_NUM='%s'", req.PhoneNum)
	}

	dataList, err = home.QueryUserInfo(where, 0, 1, "", true, com.Mysql)
	if err != nil {
		return err
	}

	if len(dataList) == 0 {
		user = home.NewUserInfo()
		user.Email = req.Email
		user.PhoneNum = req.PhoneNum

		err = home.AddUserInfo(user, com.Mysql)
		if err != nil {
			log.Error(err)
			return err
		}

		dataList = append(dataList, user)
	}

	// 在REDIS中记录会话信息
	// TODO

	// 返回信息
	rsp.UserLoginInfo = pb.NewUserLoginInfo()
	rsp.UserLoginInfo.UserID = dataList[0].UserID
	rsp.UserLoginInfo.Token = uuid.New().String()
	c.UserID = rsp.UserLoginInfo.UserID

DO_EXIT:
	err = c.Response(pb.MessageID_MID_UserLoginRsp, pb.ToBytes(rsp))

	if rsp.UserLoginInfo != nil {
		pb.DeleteUserLoginInfo(rsp.UserLoginInfo)
	}

	pb.DeleteUserLoginRsp(rsp)
	pb.DeleteUserLoginReq(req)

	if user != nil {
		home.DeleteUserInfo(user)
	}

	return err
}
