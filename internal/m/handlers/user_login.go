package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/iegad/hydra/micro"
	"github.com/iegad/hydra/mod/basic"
	"github.com/iegad/hydra/pb"
	"github.com/iegad/kraken/log"
	"github.com/iegad/sphinx/internal/com"
	"google.golang.org/protobuf/proto"
)

type UserLogin struct {
}

func (this_ *UserLogin) MID() pb.MessageID {
	return pb.MessageID_MID_UserLoginReq
}

func (this_ *UserLogin) Do(user *micro.User, in *pb.Package) error {
	var (
		rsp      *pb.UserLoginRsp
		uinfo    *pb.UserInfo
		where    string
		dataList []*pb.UserInfo
		req      = pb.NewUserLoginReq()
		err      = proto.Unmarshal(in.Data, req)
	)

	for dwf := true; dwf; dwf = false {
		if err != nil {
			user.Kick(true)
			break
		}

		rsp = pb.NewUserLoginRsp()

		// 入参检查
		if len(req.Email) == 0 && len(req.PhoneNum) == 0 {
			rsp.Code = -1
			break
		}

		if len(req.VCode) == 0 {
			rsp.Code = -2
			break
		}

		// 查询数据库
		if len(req.Email) > 0 {
			where = fmt.Sprintf("F_EMAIL='%s'", req.Email)
		} else {
			where = fmt.Sprintf("F_PHONE_NUM='%s'", req.PhoneNum)
		}

		dataList, err = basic.QueryUserInfo(where, 0, 1, "", true, com.Mysql)
		if err != nil {
			rsp.Code = -103
			break
		}

		if len(dataList) == 0 {
			uinfo = pb.NewUserInfo()
			uinfo.Email = req.Email
			uinfo.PhoneNum = req.PhoneNum

			err = basic.AddUserInfo(uinfo, com.Mysql)
			if err != nil {
				rsp.Code = -100
				break
			}

			dataList = append(dataList, uinfo)
		}

		// 在REDIS中记录会话信息
		// TODO

		// 返回信息
		dataList[0].Ver = ""
		dataList[0].VerCode = 0
		rsp.UserLoginInfo = pb.NewUserLoginInfo()
		rsp.UserLoginInfo.UserInfo = dataList[0]
		rsp.UserLoginInfo.UserID = dataList[0].UserID
		rsp.UserLoginInfo.Token = uuid.New().String()
		user.UserID = rsp.UserLoginInfo.UserID
	}

	if user.Valid() {
		err = user.Response(in.Seq, pb.MessageID_MID_UserLoginRsp, pb.ToBytes(rsp))
	}

	if err != nil {
		log.Error(err)
	}

	pb.DeleteUserLoginInfo(rsp.UserLoginInfo)
	pb.DeleteUserLoginRsp(rsp)
	pb.DeleteUserLoginReq(req)
	pb.DeleteUserInfo(uinfo)

	return nil
}
