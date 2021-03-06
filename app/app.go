package app

import (
	"crypto-colly/api"
	"crypto-colly/common/db"
	"crypto-colly/common/redis"
	"crypto-colly/record"
	"crypto-colly/setting"
	"github.com/panjf2000/ants/v2"
	"github.com/robfig/cron/v3"
	"sync"

	"crypto-colly/const"
	//"crypto-colly/crawler"
	"crypto-colly/models"
	"fmt"
)

type App struct {
	conf *setting.Config
	db *db.Db
	redis *redis.Redis
}


func NewApp(conf *setting.Config,db *db.Db,redis *redis.Redis ) *App {
	return &App{conf: conf,db: db,redis: redis}
}

func (a *App) Do() {
	fmt.Println("===================================Lance App=========================================")
	//设置想要查询的区块链
	blockchain := models.NewBlockChain("BSC",1,constants.Moralis_speedy_node)
	//采用线程池，协程组抓取BSC结点区块数据以获取nft信息
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		num := i.(int)
		record.NewRecordBlock(blockchain, a.db, a.redis,num).Do()
		wg.Done()
	})
	defer p.Release()

	for i := 8;i < 14; i++ {
		wg.Add(1)
		 _ = p.Invoke(i)
	}
	wg.Wait()
	fmt.Println("======================================" +
		"Record block finished=======================" +
		"===============================")

	c := cron.New()
	c.AddFunc("@daily",func(){
		fmt.Println("=====Start querying Bsc market top 100 collections")
		go api.NewCollection(blockchain,constants.BscCollectionListApi,constants.BscCollectionDetaiApi,constants.BscItemDetailApi,a.db).Run()
	})
	c.Start()
	//检测最新生成的区块
	//go NewRecordBlock(&blockchain, a.db, a.redis,14).Do()

	//查询BSC Market 所有上架过的商品
	//go crawler.NewApi(&blockchain,bscItemDetailApi,bscItemListApi,a.db,a.redis).Run()
	done := make(chan bool, 1)
	for {
		select {
		case <-done:
			print("退出程序")
		}
	}
}
