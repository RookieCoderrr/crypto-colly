package main

import (
	"crypto-colly/common/db"
	"crypto-colly/common/redis"
	"crypto-colly/config"
	"crypto-colly/models"
	"fmt"
)

type App struct {
	conf *config.Config
	db *db.Db
	redis *redis.Redis
}

func NewApp(conf *config.Config,db *db.Db,redis *redis.Redis ) *App{
	return &App{conf: conf,db: db,redis: redis}
}

func (a *App) Do() {
	fmt.Println("Success")
	blockchain := models.Blockchain{
		Name: "bsc",
		ChainId: 1,
		RPC: "https://quiet-white-tree.bsc.quiknode.pro/ae4802ce03ff19567834f9e82226b3dab9b92f00/",
		}
	go NewRecordBlock(&blockchain, a.db, a.redis).Do()
	done := make(chan bool, 1)
	for {
		select {
		case <-done:
			print("退出程序")
		}
	}

}