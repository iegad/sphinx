package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/iegad/hydra/micro"
	"github.com/iegad/kraken/log"
	"github.com/iegad/sphinx/internal/cfg"
	"github.com/iegad/sphinx/internal/com"
	"github.com/iegad/sphinx/internal/m"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = cfg.Init(root + "/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = com.Init()
	if err != nil {
		log.Fatal(err)
	}

	server, err := micro.NewHydra(&micro.Option{
		Project:   cfg.Instance.Server.Project,
		Service:   cfg.Instance.Server.Service,
		EtcdHosts: cfg.Instance.Etcd.Hosts,
	})
	if err != nil {
		log.Fatal(err)
	}

	server.Regist(&m.UserLogin{})

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT)
	go func() {
		<-done
		server.Stop()
	}()

	err = server.Run("cerberus/node/tcp")
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
