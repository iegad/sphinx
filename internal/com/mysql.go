package com

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/iegad/sphinx/internal/cfg"
)

var Mysql *sql.DB

func initMysql() error {
	config := mysql.NewConfig()
	config.Net = "tcp"
	config.Addr = fmt.Sprintf("%s:%d", cfg.Instance.MySql.Addr, cfg.Instance.MySql.Port)
	config.User = cfg.Instance.MySql.User
	config.Passwd = cfg.Instance.MySql.Pass
	config.Params = map[string]string{"charset": "utf8mb4"}

	if Mysql != nil {
		Mysql.Close()
	}

	var err error

	Mysql, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return err
	}

	return nil
}
