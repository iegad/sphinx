package m

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iegad/hydra/mod/home"
	"github.com/iegad/hydra/pb"
	"github.com/iegad/kraken/log"
	"github.com/iegad/kraken/piper"
	"github.com/iegad/sphinx/internal/com"
)

func (this_ *Sphinx) UserLogin(c *piper.Context, req *pb.UserLoginReq, rsp *pb.UserLoginRsp) error {
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
		where = fmt.Sprintf("F_PHONE='%s'", req.PhoneNum)
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

	// 返回信息
	rsp.UserLoginInfo = &pb.UserLoginInfo{
		UserID: dataList[0].UserID,
		Token:  uuid.New().String(),
	}

	return nil
}
