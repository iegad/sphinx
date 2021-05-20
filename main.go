package main

import (
	"github.com/iegad/hydra/micro"
	"github.com/iegad/kraken/log"
)

func main() {
	service, err := micro.NewHydra(&micro.Option{
		Project:   "sphinx",
		Service:   "sphinx",
		EtcdHosts: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		log.Fatal(err)
	}

	
	service.Run("cerberus")
}
