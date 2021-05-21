package m

import (
	"github.com/iegad/kraken/log"
	"github.com/iegad/kraken/nw/server"
	"github.com/iegad/kraken/piper"
)

type Sphinx struct {
}

func (this_ *Sphinx) OnConnected(conn server.IConn) {
	log.Info("%s has connected", conn.RouteKey())
}

func (this_ *Sphinx) OnDisconnect(conn server.IConn) {
	log.Info("%s has disconnected", conn.RouteKey())
}

func (this_ *Sphinx) OnStop(svr *piper.Server) {
	log.Info("server is stopped")
}

func (this_ *Sphinx) OnRun(svr *piper.Server) {
	log.Info("server is running")
}

func (this_ *Sphinx) Decode(c server.IConn, data []byte) ([]byte, error) {
	return data, nil
}

func (this_ *Sphinx) Encode(c server.IConn, data []byte) ([]byte, error) {
	return data, nil
}
