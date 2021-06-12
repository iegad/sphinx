package handlers

import (
	"github.com/iegad/hydra/micro"
	"github.com/iegad/hydra/mod/basic"
	"github.com/iegad/hydra/pb"
	"github.com/iegad/kraken/log"
	"github.com/iegad/sphinx/internal/com"
	"google.golang.org/protobuf/proto"
)

type ModifyUserInfo struct {
}

func (this_ *ModifyUserInfo) MID() pb.MessageID {
	return pb.MessageID_MID_ModifyUserInfoReq
}

func (this_ *ModifyUserInfo) Do(user *micro.User, in *pb.Package) error {
	var (
		rsp *pb.ModifyUserInfoRsp
		req = pb.NewModifyUserInfoReq()
		err = proto.Unmarshal(in.Data, req)
	)

	for dwf := true; dwf; dwf = false {
		if err != nil {
			user.Kick(true)
			break
		}

		rsp = pb.NewModifyUserInfoRsp()

		if req.UserInfo == nil {
			rsp.Code = -1
			break
		}

		if req.UserInfo.UserID <= 0 {
			rsp.Code = -2
			break
		}

		if len(req.Token) != 36 {
			rsp.Code = -3
			break
		}

		err = basic.ModifyUserInfo(req.UserInfo, com.Mysql)
		if err != nil {
			log.Error(err)
			rsp.Code = -102
		}
	}

	if user.Valid() {
		err = user.Response(in.Seq, pb.MessageID_MID_ModifyUserInfoRsp, pb.ToBytes(rsp))
	}

	if err != nil {
		log.Error(err)
	}

	pb.DeleteModifyUserInfoReq(req)
	pb.DeleteModifyUserInfoRsp(rsp)

	return nil
}
