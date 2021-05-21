package test

import (
	"testing"

	"github.com/iegad/hydra/pb"
	"github.com/iegad/kraken/nw"
	"github.com/iegad/kraken/piper"
)

func TestUserLogin(t *testing.T) {
	stub, err := piper.NewClient(&piper.ClientOption{
		Protocol: nw.PROTOCOL_KCP,
		Host:     "127.0.0.1:10000",
		Timeout:  5,
	})
	if err != nil {
		t.Error(err)
		return
	}

	defer stub.Close()

	rsp := &pb.UserLoginRsp{}
	err = stub.Call("UserLogin", &pb.UserLoginReq{
		Email: "123@qq.com",
		VCode: "123456",
	}, rsp)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(pb.ToJSON(rsp))
}
