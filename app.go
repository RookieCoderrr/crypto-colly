package main

import (
	"crypto-colly/common/db"
	"crypto-colly/common/redis"
	"crypto-colly/config"
	//"crypto-colly/crawler"
	"crypto-colly/models"
	"fmt"
)

type App struct {
	conf *config.Config
	db *db.Db
	redis *redis.Redis
}
const (
	moralis_speedy_node = "https://speedy-nodes-nyc.moralis.io/b9c9a1b11e9b39df1b2c3baf/bsc/mainnet"
	moralis_server = "https://cpzwdcel5tyw.usemoralis.com:2053/server"
	quick_node = "https://quiet-white-tree.bsc.quiknode.pro/ae4802ce03ff19567834f9e82226b3dab9b92f00/"
	bscUrl = "https://bscscan.com/tokens-nft"
	bscDetailApi = "https://www.binance.com/bapi/nft/v1/friendly/nft/nft-trade/product-detail"
	bscListApi = "https://www.binance.com/bapi/nft/v1/friendly/nft/product-list"
)

func NewApp(conf *config.Config,db *db.Db,redis *redis.Redis ) *App{
	return &App{conf: conf,db: db,redis: redis}
}

func (a *App) Do() {
	fmt.Println("======Lance App======")
	blockchain := models.Blockchain{
		Name: "BSC",
		ChainId: 1,
		RPC: moralis_speedy_node,
		}
	for i := 1; i < 6; i++ {
		go NewRecordBlock(&blockchain, a.db, a.redis,i).Do()
	}
	//go crawler.NewApi(&blockchain,bscDetailApi,bscListApi,a.db,a.redis).Run()
	//go crawler.NewNftMarket(bscUrl,a.db)
	done := make(chan bool, 1)
	for {
		select {
		case <-done:
			print("退出程序")
		}
	}

}