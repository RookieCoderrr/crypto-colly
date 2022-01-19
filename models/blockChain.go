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
	Email             string `json:"email"`
	Logo              string `json:"logo"`
	BlockchainId      int64  `json:"blockchain_id"`
	Contract          string `json:"contract"`
	Note              string `json:"note"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	Website           string `json:"website"`
	Telegram          string `json:"telegram"`
	Twitter           string `json:"twitter"`
	Weibo             string `json:"weibo"`
	Discord           string `json:"discord"`
	CatId             int64  `json:"cat_id"`
	Status            int    `json:"status"`
	IsShow            int    `json:"is_show"`
	Sort              int64  `json:"sort"`
	CreateTime        int64  `json:"create_time"`
	UpdateTime        int64  `json:"update_time"`
	NFTName           string `json:"nft_name"`
	NFTSymbol         string `json:"nft_symbol"`
	NFTCount          int64  `json:"nft_count"`
	NFTBurnCount      int64  `json:"nft_burn_count"`
	BlockchainName    string `json:"blockchain_name" `
	BlockchainChainId int    `json:"blockchain_chain_id" `
	CatName           string `json:"cat_name" `
	CatSlug           string `json:"cat_slug" `
	ErcType           string `json:"erc_type"`
	JsonFormat        string `json:"json_format"`
	Facebook          string `json:"facebook"`
	Github            string `json:"github"`
	OpenseaStatus     int    `json:"opensea_status"`
	BlockHeight       uint64 `json:"block_height"`    // 发现块高
	Tx                string `json:"tx"`              // 交易hash
	IsAutoInclude     bool   `json:"is_auto_include"` // 是否为自动收录
}