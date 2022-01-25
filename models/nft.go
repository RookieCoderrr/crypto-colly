package models

import (
	"context"
	"crypto-colly/common/db"
	"fmt"
	"time"
)

type NftModel struct {
	conn *db.Db
}

type NFTCollect struct {
	BlockchainId      int64  `json:"blockchain_id"`
	Contract          string `json:"contract"`
	CreateTime        int64  `json:"create_time"`
	UpdateTime        int64  `json:"update_time"`
	NFTName           string `json:"nft_name"`
	NFTSymbol         string `json:"nft_symbol"`
	ErcType           string `json:"erc_type"`
	BlockHeight       uint64 `json:"block_height"`    // 发现块高
	Tx                string `json:"tx"`              // 交易hash
	MarketPlace       int    `json:"market_place"`
	IsPopular         bool   `json:"is_popular"`
}

func NewNftModel(conn *db.Db) *NftModel{
	return &NftModel{conn: conn}
}


func (n *NftModel) CreateNft(blockchainId int64, address, ercType, name, symbol string, blockHeight uint64, tx string) (int64, error) {
	data := NFTCollect{
		BlockchainId:  blockchainId,
		Contract:      address,
		ErcType:       ercType,
		NFTName:       name,
		NFTSymbol:     symbol,
		BlockHeight:   blockHeight,
		Tx:            tx,
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
	}

	_,err := n.conn.GetConn().Database("nft").Collection("info").InsertOne(context.TODO(),data)
	if err != nil {
		fmt.Println("Insert nft error")
		return 0, err
	}

	return data.BlockchainId, nil
}
