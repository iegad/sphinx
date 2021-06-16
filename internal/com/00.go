package com

import (
	"errors"

	"github.com/iegad/sphinx/internal/cfg"
)

func Init() error {
	if cfg.Instance == nil {
		return errors.New("未正确加载配置文件")
	}

	err := initMysql()
	if err != nil {
		return err
	}

	return initRedis()
}
