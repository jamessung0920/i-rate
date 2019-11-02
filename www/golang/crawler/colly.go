package crawler

import (
	"fmt"
	"time"
	"strings"

	"github.com/gocolly/colly"
	"github.com/PuerkitoBio/goquery"

	"app/common"
	"app/database"
	"app/database/models"
)

type Rate struct {
	buy string
	sell string
}

func Crawler(site string) {
	c := colly.NewCollector()

	// var rateData []models.Rate
	// var currencyRate models.Rate
	var rateTimeFormatTime time.Time
	_, currencyList := common.GetCurrencyList()
	var visitSite, tagTarget string = "", ""

	switch site {
	case "taiwan-bank":
		visitSite = "https://rate.bot.com.tw/xrt?Lang=zh-TW"
		tagTarget = "main"
	case "ctbc-bank":
		visitSite = "https://www.bestxrate.com/bankrate/twctbc.html"
		tagTarget = ".container"
	case "esun-bank":
		visitSite = "https://www.bestxrate.com/bankrate/twesun.html"
		tagTarget = ".container"
	default:
		visitSite = "https://rate.bot.com.tw/xrt?Lang=zh-TW"
		tagTarget = "main"
	}
	
	c.OnHTML(tagTarget, func(e *colly.HTMLElement) {
		fmt.Println("hshhhshshshshshshshshshh")
		if site == "taiwan-bank" {
			e.DOM.Find("span.time").Each(func(_ int, s *goquery.Selection) {
				rateTimeFormatTime = common.StringToTime(strings.TrimSpace(s.Text()))
			})
		}
		e.DOM.Find("table tbody tr").Each(func(_ int, s *goquery.Selection) {
			rowCurrencyData := strings.Fields(s.Text())
			if site == "ctbc-bank" || site == "esun-bank" {
				rateTimeFormatTime = common.StringToTime(rowCurrencyData[6] + " " + rowCurrencyData[7])
			}
			for key, val := range currencyList {
				// fmt.Println(key, rowCurrencyData[0], rowCurrencyData[1], rowCurrencyData[2], val)
				if strings.Contains(rowCurrencyData[1], val) {
					// rateData = append(rateData, currencyRate)
					database.GetDB().Create(&models.Rate{
						RateTime: rateTimeFormatTime,
						BuyRate: rowCurrencyData[4],
						SellRate: rowCurrencyData[5],
						CrawlFrom: site,
						CurrencyID: uint(key+1),
					})
				}
			}
			fmt.Println(rowCurrencyData)
			fmt.Println("===============")
		})
		// rate = append(rate, e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		// fmt.Println("rate slice", rateData)
	})

	c.Visit(visitSite)
	// c.Visit("https://en.wikipedia.org/wiki/List_of_S&P_500_companies")
}