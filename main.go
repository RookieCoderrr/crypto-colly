package main

import (
	"context"
	"crypto-colly/app"
	cfg "crypto-colly/common/config"
	"crypto-colly/common/db"
	"crypto-colly/common/redis"
	"crypto-colly/setting"
	"flag"
)

var confFile = flag.String("c","setting.yml","setting file")

func main() {
	var (
		conf   *setting.Config
		redisConn   *redis.Redis
		dbConn *db.Db
		err    error
	)
	flag.Parse()
	conf  = new(setting.Config)
	cfg.NewConfig(conf).Read(*confFile)
	redisConn = redis.InitializeRedisLocalClient(&conf.Redis)
	dbConn, err = db.InitializeMongoLocalClient(context.TODO(),&conf.Db)
	redisConn.Test()
	if err != nil {
		panic(err)
	}
	app.NewApp(conf,dbConn,redisConn).Do()
}