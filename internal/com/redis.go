package com

import (
	"context"
	"runtime"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/iegad/sphinx/internal/cfg"
)

var Redis *redis.Client

func initRedis() error {
	Redis = redis.NewClient(&redis.Options{
		Addr:        cfg.Instance.Redis.Hosts[0],
		Username:    cfg.Instance.Redis.User,
		Password:    cfg.Instance.Redis.Pass,
		DB:          cfg.Instance.Redis.DB,
		DialTimeout: time.Second * 3,
		PoolSize:    runtime.NumCPU(),
	})

	return Redis.Ping(context.TODO()).Err()
}
