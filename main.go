package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/iegad/hydra/micro"
	"github.com/iegad/kraken/log"
	"github.com/iegad/sphinx/internal/cfg"
	"github.com/iegad/sphinx/internal/com"
	"github.com/iegad/sphinx/internal/m"
)

func main() {
	// Step 1: 获取当前路径
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Step 2: 获取配置
	err = cfg.Init(root + "/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Step 3: 初始化组件
	err = com.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Step 4: 构建服务节点
	server, err := micro.NewHydra(&micro.Option{
		Project:      cfg.Instance.Server.Project,
		Service:      cfg.Instance.Server.Service,
		EtcdHosts:    cfg.Instance.Etcd.Hosts,
		OnUserClosed: m.OnUserClosed,
		OnIdempotent: m.OnIdempotent,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Step 5: 注册服务句柄
	server.Regist(&m.UserLogin{})
	server.Regist(&m.UserUnregist{})

	// Step 6: 注册退出信号
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-done
		server.Stop()
	}()

	// Step 7: 启动节点服务
	err = server.Run("cerberus/node/tcp")
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
