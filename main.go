package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/iegad/kraken/log"
	"github.com/iegad/kraken/piper"
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

	server, err := piper.NewServer(&m.Sphinx{}, &piper.ServerOption{
		ID:        cfg.Instance.Server.ID,
		Protocol:  cfg.Instance.Server.Protocol,
		Service:   cfg.Instance.Server.Service,
		EtcdHosts: cfg.Instance.Etcd.Hosts,
		Host:      cfg.Instance.Server.Host,
		MaxConn:   cfg.Instance.Server.MaxConn,
		Timeout:   cfg.Instance.Server.Timeout,
	})
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT)
	go func() {
		<-done
		server.Stop()
	}()

	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
