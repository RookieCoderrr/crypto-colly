package main

import (
	"crypto-colly/api"
	"crypto-colly/common/db"
	"crypto-colly/common/redis"
	"crypto-colly/config"
	"github.com/robfig/cron/v3"

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
	bscItemDetailApi = "https://www.binance.com/bapi/nft/v1/friendly/nft/nft-trade/product-detail"
	bscItemListApi = "https://www.binance.com/bapi/nft/v1/friendly/nft/product-list"
	bscCollectionListApi = "https://www.binance.com/bapi/nft/v1/public/nft/ranking/top-collections/1/100"
	bscCollectionDetaiApi = "https://www.binance.com/bapi/nft/v1/friendly/nft/layer-product-list"
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

	//多线程抓取BSC结点区块数据以获取nft信息
	for i := 8; i < 14; i++ {
		go NewRecordBlock(&blockchain, a.db, a.redis,i).Do()
	}
	c := cron.New()
	c.AddFunc("@daily",func(){
		fmt.Println("=====Start querying Bsc market top 100 collections")
		go api.NewCollection(&blockchain,bscCollectionListApi,bscCollectionDetaiApi,bscItemDetailApi,a.db).Run()
	})
	c.Start()
	//检测最新生成的区块
	//go NewRecordBlock(&blockchain, a.db, a.redis,14).Do()

	//查询BSC Market 所有上架过的商品
	//go crawler.NewApi(&blockchain,bscItemDetailApi,bscItemListApi,a.db,a.redis).Run()


	// 查询BSC Market top collection nft
	//go api.NewCollection(&blockchain,bscCollectionListApi,bscCollectionDetaiApi,bscItemDetailApi,a.db).Run()

	done := make(chan bool, 1)
	for {
		select {
		case <-done:
			print("退出程序")
		}
	}
}
