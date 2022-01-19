package main

import (
	"context"
	cfg "crypto-colly/common/config"
	"crypto-colly/common/db"
	"crypto-colly/common/redis"
	"crypto-colly/config"
	"flag"
)

var confFile = flag.String("c","config.yml","config file")

func main() {
	var (
		conf   *config.Config
		redisConn   *redis.Redis
		dbConn *db.Db
		err    error
	)
	flag.Parse()
	conf  = new(config.Config)
	cfg.NewConfig(conf).Read(*confFile)
	redisConn = redis.InitializeRedisLocalClient(&conf.Redis)
	dbConn, err = db.InitializeMongoLocalClient(context.TODO(),&conf.Db)
	redisConn.Test()
	if err != nil {
		panic(err)
	}
	NewApp(conf,dbConn,redisConn).Do()
}