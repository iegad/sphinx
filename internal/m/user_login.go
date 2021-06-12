package m

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/iegad/hydra/micro"
	"github.com/iegad/hydra/mod/basic"
	"github.com/iegad/hydra/pb"
	"github.com/iegad/kraken/log"
	"github.com/iegad/kraken/utils"
	"github.com/iegad/sphinx/internal/com"
	"google.golang.org/protobuf/proto"
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
		rsp      = pb.NewUserLoginRsp()
		err      = proto.Unmarshal(in.Data, req)
		user     *pb.UserInfo
		where    = ""
		dataList []*pb.UserInfo
	)

	for dwf := true; dwf; dwf = false {
		if err != nil {
			c.Kick(true)
			break
		}

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
			break
		}

		if len(dataList) == 0 {
			user = pb.NewUserInfo()
			user.Email = req.Email
			user.PhoneNum = req.PhoneNum

			err = basic.AddUserInfo(user, com.Mysql)
			if err != nil {
				log.Error(err)
				break
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

		err = c.Response(in.Seq, pb.MessageID_MID_UserLoginRsp, pb.ToBytes(rsp))
	}

	pb.DeleteUserLoginInfo(rsp.UserLoginInfo)
	pb.DeleteUserLoginRsp(rsp)
	pb.DeleteUserLoginReq(req)
	pb.DeleteUserInfo(user)

	return err
}
