package models

type Blockchain struct {
	Name    string `json:"name"`
	ChainId int64    `json:"chain_id"`
	RPC     string `json:"rpc"`
	Graph   string `json:"graph"`
	Sort    int    `json:"sort"`
	Status  int    `json:"status" default:"0"`
}
type NFTCollect struct {
	BlockchainId      int64  `json:"blockchain_id"`
	Contract          string `json:"contract"`
	Title             string `json:"title"`
	Status            int    `json:"status"`
	CreateTime        int64  `json:"create_time"`
	UpdateTime        int64  `json:"update_time"`
	NFTName           string `json:"nft_name"`
	NFTSymbol         string `json:"nft_symbol"`
	ErcType           string `json:"erc_type"`
	BlockHeight       uint64 `json:"block_height"`    // 发现块高
	Tx                string `json:"tx"`              // 交易hash
	IsAutoInclude     bool   `json:"is_auto_include"` // 是否为自动收录
	MarketPlace       int    `json:"market_place"`
	isPopular         bool   `json:"is_popular"`
}