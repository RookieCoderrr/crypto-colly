package main

import (
	"context"
	"crypto-colly/common/chainutils"
	"crypto-colly/common/db"
	"crypto-colly/common/redis"
	"crypto-colly/contract/erc721"
	"crypto-colly/models"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
	"time"
)

const (
	ProcessBlockHeightPrefix = "process_block_height_"
	NotifyWatchKey           = "PUBSUB_WATCH"
)

type RecordBlock struct {
	chain              *models.Blockchain
	db                 *db.Db
	redis              *redis.Redis
	client             *ethclient.Client
	model              *models.NftModel
	processBlockHeight *big.Int
	currentBlockHeight *big.Int
	crawling           bool
	startTime          time.Time
}
func NewRecordBlock(chain *models.Blockchain, db *db.Db, redis *redis.Redis) *RecordBlock{
	fmt.Println("newRecordBlock")
	return &RecordBlock{
		chain: chain,
		db: db,
		redis: redis,
		model: models.NewNftModel(db),
		processBlockHeight: new(big.Int),
		currentBlockHeight: new(big.Int),
		startTime: time.Now(),
	}

}
func (r *RecordBlock) Do(){
	fmt.Println("do")
	client, err := ethclient.Dial(r.chain.RPC)
	if err != nil {
		fmt.Sprintf("连接(%s)[%s]失败: %s\n", r.chain.Name, r.chain.RPC, err)
		return
	}
	r.client = client

	fmt.Sprintf("连接(%s)[%s]成功\n", r.chain.Name, r.chain.RPC)
	lastProcessBlockHeight, err := r.getProcessedBlockHeight()
	if err != nil {
		fmt.Sprintf("(%s)获取上次处理块高失败: %s\n", r.chain.Name, err)
		return
	}
	r.processBlockHeight = lastProcessBlockHeight
	fmt.Sprintf("(%s)开始爬取合约，上次处理块高: %s\n", r.chain.Name, lastProcessBlockHeight.String())

	go r.autoGetCurrentBlockHeight()
	r.autoCrawl()

}
func (r *RecordBlock) getProcessedBlockHeight() (*big.Int, error) {
	var (
		blockHeight = new(big.Int)
		err         error
	)

	result, err := r.redis.Do("GET", ProcessBlockHeightPrefix+strings.ToLower(r.chain.Name))
	if err != nil {
		return blockHeight, err
	}

	if result == nil {
		return blockHeight, nil
	}
	//fmt.Println(123123213123123)
	blockHeight.SetString(string(result.([]byte)), 10)
	return blockHeight, nil
}

func (r *RecordBlock) autoGetCurrentBlockHeight() {
	fmt.Println("autoGetCurrentBlockHeight")
	tick := time.Tick(3 * time.Second)
	for {
		select {
		case <-tick:
			r.getCurrentBlockHeight()
		}
	}
}

func (r *RecordBlock) getCurrentBlockHeight() {
	header, err := r.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		fmt.Sprintf("(%s)获取当前块高失败: %s\n", r.chain.Name, err)
		return
	}

	r.currentBlockHeight = header.Number
	//fmt.Println(r.currentBlockHeight)

	var diff = new(big.Int).Sub(r.currentBlockHeight, r.processBlockHeight)
	if diff.Cmp(big.NewInt(10)) > 0 {
		fmt.Sprintf("(%s)待处理: %s 块\n", r.chain.Name, diff.String())
	}
}

func (r *RecordBlock) autoCrawl() {
	fmt.Println("autoCrawl")
	tick := time.Tick(3 * time.Second)
	for {
		select {
		case <-tick:
			if !r.crawling && r.processBlockHeight.Cmp(r.currentBlockHeight) <= 0 {
				go r.crawl()
			}
		}
	}
}

func (r *RecordBlock) crawl() {
	r.crawling = true
	for {
		block, err := r.client.BlockByNumber(context.Background(), r.processBlockHeight)
		if err != nil {
			fmt.Sprintf("(%s)[%d]获取块数据失败: %s\n", r.chain.Name, r.processBlockHeight, err)
			break
		}
		fmt.Println(block.Transactions())
		fmt.Println(block.Hash())

		for _, tx := range block.Transactions() {
			fmt.Println("transaction")
			// 试试这个， to为空就当是合约判断
			if tx.To() == nil {
				fmt.Println("==============================================================================================")
				err := r.analyzeTx(tx)
				if err != nil {
					continue
				}
			}

			// 只有当tx data 足够大的时候，才被解析，否则跳过
			//if len(tx.Data()) > 7000 {
			//
			//	err := a.analyzeTx(tx)
			//	if err != nil {
			//		continue
			//	}
			//}
		}

		err = r.saveProcessedBlockHeight(r.processBlockHeight)
		if err != nil {
			fmt.Sprintf("(%s)[%d]保存处理进度失败: %s\n", r.chain.Name, r.processBlockHeight, err)
			break
		}

		r.processBlockHeight = new(big.Int).Add(r.processBlockHeight, big.NewInt(1))
		fmt.Println(r.processBlockHeight)
		if r.processBlockHeight.Cmp(r.currentBlockHeight) > 0 {
			break
		}
	}
	r.crawling = false
}
func (r *RecordBlock) saveProcessedBlockHeight(blockHeight *big.Int) error {
	_, err := r.redis.Do("SET", ProcessBlockHeightPrefix+strings.ToLower(r.chain.Name), blockHeight.String())
	fmt.Sprintf("Save block height: %d",blockHeight)
	return err
}

func (r *RecordBlock) analyzeTx(tx *types.Transaction) error {
	fmt.Println("analyzeTx")
	receipt, err := r.client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		fmt.Sprintf("(%s)[%d]解析交易失败(%s): %s\n", r.chain.Name, r.processBlockHeight, tx.Hash().Hex(), err)
		return err
	}

	//tx, isPedding, err := a.client.TransactionByHash(context.Background(), tx.Hash())
	//receipt, err := bind.WaitMined(context.Background(), a.client, tx)

	if receipt.ContractAddress.Hex() != "0x0000000000000000000000000000000000000000" {
		ercType, err := chainutils.JudgmentErcType(receipt.ContractAddress, r.client)
		if err != nil {
			fmt.Sprintf("(%s)[%d]判断合约类型失败(tx: %s, contract: %s): %s\n", r.chain.Name,
				r.processBlockHeight, tx.Hash().Hex(), receipt.ContractAddress.Hex(), err)
			return err
		}

		switch ercType {
		case chainutils.Erc721:
			err := r.recordErc721(receipt.ContractAddress, tx.Hash().Hex())
			if err != nil {
				fmt.Sprintf("(%s)[%d]保存合约失败(tx: %s, contract: %s, erc_type: %s): %s\n", r.chain.Name,
					r.processBlockHeight, tx.Hash().Hex(), receipt.ContractAddress.Hex(), "erc721", err)
				return err
			}
			break
		case chainutils.Unknown:
			break
		}
	}

	return nil
}

func (r *RecordBlock) recordErc721(address common.Address, tx string) error {
	addr := strings.ToLower(address.Hex())
	instance, _ := erc721.InitInstance(r.chain.RPC, address.Hex())
	name, _ := instance.Name(&bind.CallOpts{})
	symbol, _ := instance.Symbol(&bind.CallOpts{})
	_, err := r.model.CreateNft(r.chain.ChainId, addr, "erc721", name, symbol, r.processBlockHeight.Uint64(), tx)
	fmt.Sprintf("(%s)[%d]新收录(contract: %s, erc_type: %s, name: %s, symbol: %s)\n", r.chain.Name,
		r.processBlockHeight, addr, "erc721", name, symbol)
	return err
}

