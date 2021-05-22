package m

import (
	"github.com/iegad/hydra/micro"
	"github.com/iegad/kraken/log"
)

func UserClosed(user *micro.User) {
	log.Debug("%d[%s] has disconnected", user.UserID, user.RemoteAddr())
}
