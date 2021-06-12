package m

import (
	"github.com/iegad/hydra/micro"
	"github.com/iegad/hydra/pb"
	"github.com/iegad/kraken/log"
)

// OnUserClosed 用户连接断开事件
func OnUserClosed(user *micro.User) {
	log.Debug("%d[%s] has disconnected", user.UserID, user.RemoteAddr())
}

// OnIdempotent 包重复事件
func OnIdempotent(user *micro.User, pack *pb.Package) {
	log.Debug("%v => %v", user.UserID, pb.ToJSON(pack))
}
