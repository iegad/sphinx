package test

import (
	"errors"
	"testing"
	"time"

	"github.com/iegad/hydra/pb"
	"github.com/iegad/kraken/log"
	"github.com/iegad/kraken/nw/client"
	"github.com/iegad/kraken/nw/client/tcp"
	"google.golang.org/protobuf/proto"
)

const (
	NODE_ADDR = "127.0.0.1:65401"
)

func TestUserLogin(t *testing.T) {
	cli, err := getClient()
	if err != nil {
		t.Error(err)
		return
	}

	defer cli.Close()

	err = userLogin(cli)
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < 50; i++ {
		err = ping(cli)
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(time.Second)
	}
}

func getClient() (client.IClient, error) {
	return tcp.NewClient(&client.Option{
		Host:    "127.0.0.1:1751",
		Timeout: 30,
	})
}

func userLogin(cli client.IClient) error {
	req := pb.NewUserLoginReq()
	req.Email = "11111@1111.com"
	req.VCode = "12355"

	out := pb.NewPackage()
	out.PID = pb.PackageID_PID_UserDelivery
	out.MID = pb.MessageID_MID_UserLoginReq
	out.Idempotent = time.Now().UnixNano()
	out.Data = pb.ToBytes(req)
	out.ToNodeAddr = NODE_ADDR

	reqbf := pb.ToBytes(out)
	pb.DeleteUserLoginReq(req)
	pb.DeletePackage(out)

	err := cli.Write(reqbf)
	if err != nil {
		return err
	}

	rspbf, err := cli.Read()
	if err != nil {
		return err
	}

	rsp := pb.NewUserLoginRsp()
	defer pb.DeleteUserLoginRsp(rsp)

	in := pb.NewPackage()
	defer pb.DeletePackage(in)

	err = proto.Unmarshal(rspbf, in)
	if err != nil {
		return err
	}

	if in.PID != pb.PackageID_PID_NodeDelivery {
		return errors.New("---1")
	}

	if in.MID != pb.MessageID_MID_UserLoginRsp {
		return errors.New("---2")
	}

	err = proto.Unmarshal(in.Data, rsp)
	if err != nil {
		return err
	}

	if rsp.Code != 0 {
		return errors.New(rsp.Error)
	}

	log.Debug(rsp.UserLoginInfo)
	return nil
}

func ping(cli client.IClient) error {
	out := pb.NewPackage()
	out.PID = pb.PackageID_PID_Ping
	out.Idempotent = time.Now().UnixNano()

	data := pb.ToBytes(out)
	pb.DeletePackage(out)

	err := cli.Write(data)
	if err != nil {
		return err
	}

	in := pb.NewPackage()
	defer pb.DeletePackage(in)
	rspbf, err := cli.Read()
	if err != nil {
		return err
	}

	err = proto.Unmarshal(rspbf, in)
	if err != nil {
		return err
	}

	log.Debug(pb.ToJSON(in))
	return nil
}
