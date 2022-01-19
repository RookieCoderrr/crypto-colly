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

func NewNftModel(conn *db.Db) *NftModel{
	return &NftModel{conn: conn}
}

func (n *NftModel) CreateNft(blockchainId int64, address, ercType, name, symbol string, blockHeight uint64, tx string) (int64, error) {
	data := NFTCollect{
		BlockchainId:  blockchainId,
		Contract:      address,
		ErcType:       ercType,
		Title:         name,
		NFTName:       name,
		NFTSymbol:     symbol,
		BlockHeight:   blockHeight,
		Tx:            tx,
		IsAutoInclude: true,
		Status:        1,
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
