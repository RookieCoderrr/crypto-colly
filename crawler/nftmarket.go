package crawler

import (
	"crypto-colly/common/db"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"net/http"
	"strconv"
	"strings"
)

type NftMarket struct {
	url string
	db *db.Db
}

func NewNftMarket (url string, db *db.Db) *NftMarket{
	return &NftMarket{
		url:url,
		db :db,
	}
}
func (n *NftMarket) crawl()  {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.MaxDepth(2),
		colly.Debugger(&debug.LogDebugger{}),
		)
	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})
	detailCollector := c.Clone()
	//rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:1081")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//c.SetProxyFunc(rp)
	c.OnHTML("body", func(e *colly.HTMLElement){
		fmt.Println("body")
		e.ForEach("table tr", func(_ int, el *colly.HTMLElement){
			detailURL := "https://bscscan.com"+el.ChildAttr("td:nth-child(2) h3 div a","href")
			detailCollector.Visit(detailURL)

		})
		nextPage := e.ChildAttr("div[id=ContentPlaceHolder1_divPagination] ul li:nth-child(4) a", "href")
		fmt.Println("nextPage:", nextPage)
		if pageReg.MatchString(nextPage){
			c.Visit(e.Request.AbsoluteURL(nextPage))
		}
	})

	detailCollector.OnHTML("div[id=ContentPlaceHolder1_divSummary]>div:nth-child(1)", func(e *colly.HTMLElement){
		totalSupply, err := strconv.ParseFloat(strings.ReplaceAll(e.ChildText("div:nth-child(1)>div>div:nth-child(2)>div:nth-child(1)>div:nth-child(2)>span:nth-of-type(1) "),",",""),64)
		if err != nil {
			fmt.Println("convert totalSupply error")
		}
		symbol := e.ChildText("div:nth-child(1)>div>div:nth-child(2)>div:nth-child(1)>div:nth-child(2)>b")
		holders,err:= strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(e.ChildText("div[id=ContentPlaceHolder1_tr_tokenHolders] div div div div ")," addresses",""),",",""),64)
		if err != nil {
			fmt.Println("convert holders error")
		}
		contractHash := e.ChildText("div:nth-child(2)>div>div:nth-child(2)>div:nth-child(1)>div:nth-child(2) a")
		fmt.Println(totalSupply)
		fmt.Println(symbol)
		fmt.Println(holders)
		fmt.Println(contractHash)
		raw := NFTInfo{
			ContractHash :contractHash,
			Symbol :symbol,
			TotalSupply :totalSupply,
			TotalHolders :holders,
		}
		fmt.Println("raw:",raw)

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.Visit(fmt.Sprintf("%s?p=%d", n.url, 1))
}