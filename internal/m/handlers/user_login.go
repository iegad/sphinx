package handlers

import (
	"fmt"

	"github.com/iegad/hydra/micro"
	"github.com/iegad/hydra/mod/basic"
	"github.com/iegad/hydra/mod/cache"
	"github.com/iegad/hydra/pb"
	"github.com/iegad/kraken/log"
	"github.com/iegad/kraken/utils"
	"github.com/iegad/sphinx/internal/com"
	"google.golang.org/protobuf/proto"
)

/**
 * 用户登录:
 *   用户登录有两种方式
 *   1: 冷登录: 需要带入参数 phoneNum/email, vcode
 *   2: 热登录: 需要带入参数 userID, account, token, termInfo
 */

type UserLogin struct {
}

func (this_ *UserLogin) MID() pb.MessageID {
	return pb.MessageID_MID_UserLoginReq
}

func (this_ *UserLogin) Do(user *micro.User, in *pb.Package) error {
	var (
		rsp *pb.UserLoginRsp
		req = pb.NewUserLoginReq()
		err = proto.Unmarshal(in.Data, req)
	)

	for dwf := true; dwf; dwf = false {
		if err != nil {
			log.Error("返序列化错误: %v", err)
			break
		}

		if len(req.VCode) == 5 {
			rsp, err = this_.login1(user, req)
		} else {
			rsp, err = this_.login2(user, req)
		}

		if err != nil {
			log.Error("登录失败: %v", err)
			break
		}

		err = user.Response(in.Seq, pb.MessageID_MID_UserLoginRsp, pb.ToBytes(rsp))
		if err != nil {
			log.Error("响应消息失败: %v", err)
		}
	}

	if rsp != nil {
		if rsp.UserSession != nil {
			pb.DeleteUserInfo(rsp.UserSession.UserInfo)
		}
		pb.DeleteUserSession(rsp.UserSession)
	}

	pb.DeleteUserLoginRsp(rsp)
	pb.DeleteUserLoginReq(req)

	return nil
}

// login1 冷登录
func (this_ *UserLogin) login1(user *micro.User, req *pb.UserLoginReq) (*pb.UserLoginRsp, error) {
	var (
		err      error
		rsp      = pb.NewUserLoginRsp()
		dataList []*pb.UserInfo
	)

	for dwf := true; dwf; dwf = false {
		// Step 1: 入参检查
		if len(req.DeviceCode) != 16 {
			rsp.Code = -1
			break
		}

		if len(req.Email) == 0 && len(req.PhoneNum) == 0 {
			rsp.Code = -2
			break
		}

		// Step 2: DB查询用户是否存在
		where := ""
		if len(req.Email) > 0 {
			where = fmt.Sprintf("F_EMAIL='%s'", req.Email)
		} else {
			where = fmt.Sprintf("F_PHONE_NUM='%s'", req.PhoneNum)
		}

		dataList, err = basic.QueryUserInfo(where, 0, 1, "", false, com.Mysql)
		if err != nil {
			log.Error("QueryUserInfo 失败: %v", err)
			break
		}

		// Step 3: 不存在, 则新建用户
		if len(dataList) == 0 {
			userInfo := pb.NewUserInfo()
			userInfo.Email = req.Email
			userInfo.PhoneNum = req.PhoneNum
			err = basic.AddUserInfo(userInfo, com.Mysql)
			if err != nil {
				rsp.Code = -100
				pb.DeleteUserInfo(userInfo)
				log.Error("AddUserInfo 失败: %v", err)
				break
			}

			dataList = append(dataList, userInfo)
		}

		// Step 4: 构建Session
		ss := pb.NewUserSession()
		ss.OSType = req.OSType
		ss.Token = utils.UUID_Bytes()
		ss.DeviceCode = req.DeviceCode
		ss.MountAddr = user.MountAddr()
		ss.UserInfo = dataList[0]

		err = cache.SetUserSess(ss, com.Redis)
		if err != nil {
			rsp.Code = -101
			pb.DeleteUserSession(ss)
			log.Error("SetUserSess 失败: %v", err)
			break
		}

		// Step 5: 构建response
		rsp.UserSession = ss
	}

	return rsp, nil
}

// login2 热登录
func (this_ *UserLogin) login2(user *micro.User, req *pb.UserLoginReq) (*pb.UserLoginRsp, error) {
	var (
		err error
		ss  *pb.UserSession
		rsp = pb.NewUserLoginRsp()
	)

	for dwf := true; dwf; dwf = false {
		// Step 1: 入参检查
		if req.UserID <= 0 {
			rsp.Code = -1
			break
		}

		if len(req.Token) != 16 {
			rsp.Code = -2
			break
		}

		if len(req.DeviceCode) != 16 {
			rsp.Code = -3
			break
		}

		if len(req.Email) == 0 && len(req.PhoneNum) == 0 {
			rsp.Code = -4
			break
		}

		// Step 2: 获取会话
		ss, err = cache.GetUserSess(req.UserID, com.Redis)
		if err != nil {
			rsp.Code = -100
			pb.DeleteUserSession(ss)
			log.Error("GetUserSess 失败: %v", err)
			break
		}

		if utils.BytesToString(ss.Token) != utils.BytesToString(req.Token) {
			rsp.Code = -2
			pb.DeleteUserSession(ss)
			break
		}

		if utils.BytesToString(ss.DeviceCode) != utils.BytesToString(req.DeviceCode) {
			rsp.Code = -3
			pb.DeleteUserSession(ss)
			break
		}

		if ss.UserInfo.PhoneNum != req.PhoneNum || ss.UserInfo.Email != req.Email {
			rsp.Code = -4
			pb.DeleteUserSession(ss)
			break
		}

		// Step 3: 重置会话
		ss.Token = utils.UUID_Bytes()
		ss.MountAddr = user.MountAddr()

		err = cache.SetUserSess(ss, com.Redis)
		if err != nil {
			rsp.Code = -100
			log.Error("SetUserSess 失败: %v", err)
			pb.DeleteUserSession(ss)
			break
		}

		// Step 4: 构建response
		rsp.UserSession = ss
	}

	return rsp, nil
}
