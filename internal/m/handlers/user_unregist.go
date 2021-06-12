package handlers

import (
	"errors"
	"fmt"

	"github.com/iegad/hydra/micro"
	"github.com/iegad/hydra/mod/basic"
	"github.com/iegad/hydra/pb"
	"github.com/iegad/kraken/utils"
	"github.com/iegad/sphinx/internal/com"
	"google.golang.org/protobuf/proto"
)

type UserUnregist struct {
}

func (this_ *UserUnregist) MID() pb.MessageID {
	return pb.MessageID_MID_UserUnregistReq
}

func (this_ *UserUnregist) Do(c *micro.User, in *pb.Package) error {
	utils.Assert(c != nil && in != nil, "userLogin.do params are invalid")

	req := pb.NewUserUnregistReq()
	err := proto.Unmarshal(in.Data, req)
	if err != nil {
		return err
	}

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

	dataList, err := basic.QueryUserInfo(where, 0, 1, "", true, com.Mysql)
	if err != nil {
		return err
	}

	if len(dataList) != 1 {
		return errors.New("F_EMAIL")
	}

	err = basic.RemoveUserInfo(dataList[0].UserID, com.Mysql)
	if err != nil {
		return err
	}

	return nil
}
