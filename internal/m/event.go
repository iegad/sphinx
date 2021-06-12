package m

import (
	"github.com/iegad/hydra/micro"
	"github.com/iegad/hydra/pb"
	"github.com/iegad/kraken/log"
)

func OnUserClosed(user *micro.User) {
	log.Debug("%d[%s] has disconnected", user.UserID, user.RemoteAddr())
}

func OnIdempotent(user *micro.User, pack *pb.Package) {
	log.Debug("%v => %v", user.UserID, pb.ToJSON(pack))
}
