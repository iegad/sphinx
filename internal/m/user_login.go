package m

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iegad/hydra/micro"
	"github.com/iegad/hydra/mod/home"
	"github.com/iegad/hydra/pb"
	"github.com/iegad/kraken/log"
	"github.com/iegad/kraken/utils"
	"github.com/iegad/sphinx/internal/com"
	"google.golang.org/protobuf/proto"
)

type UserLogin struct {
}

func (this_ *UserLogin) MID() int32 {
	return pb.MID_UserLoginReq
}

func (this_ *UserLogin) Do(c *micro.User, in *pb.Package) error {
	utils.Assert(c != nil && in != nil, "userLogin.do params are invalid")

	log.Info("-------------------------------------")

	req := pb.NewUserLoginReq()
	err := proto.Unmarshal(in.Data, req)
	if err != nil {
		log.Error(err)
		return err
	}

	// 入参检查
	if len(req.Email) == 0 && len(req.PhoneNum) == 0 {
		return errors.New("account is invalid")
	}

	if len(req.VCode) == 0 {
		return errors.New("vcode is invalid")
	}

	// 查询数据库
	var (
		where = ""
	)
	if len(req.Email) > 0 {
		where = fmt.Sprintf("F_EMAIL='%s'", req.Email)
	} else if len(req.PhoneNum) > 0 {
		where = fmt.Sprintf("F_PHONE_NUM='%s'", req.PhoneNum)
	}

	dataList, err := home.QueryUserInfo(where, 0, 1, "", true, com.Mysql)
	if err != nil {
		return err
	}

	if len(dataList) == 0 {
		user := &home.UserInfo{
			Email:    req.Email,
			PhoneNum: req.PhoneNum,
		}
		err = home.InsertUserInfo(user, com.Mysql)
		if err != nil {
			log.Error(err)
			return err
		}

		dataList = append(dataList, user)
	}

	// 在REDIS中记录会话信息
	// TODO

	// 返回信息
	rsp := &pb.UserLoginRsp{
		UserLoginInfo: &pb.UserLoginInfo{
			UserID: dataList[0].UserID,
			Token:  uuid.New().String(),
		},
	}

	c.UserID = rsp.UserLoginInfo.UserID
	return this_.response(c, rsp, time.Now().UnixNano())
}

func (this_ *UserLogin) response(c *micro.User, rsp *pb.UserLoginRsp, idempotent int64) error {
	pack := pb.NewPackage()
	pack.PID = pb.PID_NodeDelivery
	pack.MID = pb.MID_UserLoginRsp
	pack.Idempotent = idempotent
	pack.ToUserAddrs = []string{c.RemoteAddr()}
	pack.Data = pb.ToBytes(rsp)

	data := pb.ToBytes(pack)
	pb.DeletePackage(pack)
	return c.Write(data)
}
