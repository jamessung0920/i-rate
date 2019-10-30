package crawler

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
	// "github.com/PuerkitoBio/goquery"
	"app/database"
	"app/database/models"
)

type Rate struct {
	buy string
	sell string
}

func Crawler(index uint, currency string) {
	c := colly.NewCollector()

	var rate []string
	// Find and visit all links
	c.OnHTML("tr td:nth-of-type(1), tr td:nth-of-type(3), tr td:nth-of-type(4)", func(e *colly.HTMLElement) {
		// fmt.Println("本行買入:", e.Text)
		// e.Request.Visit(e.Attr("href"))
		// e.DOM.Children().First().Find("tr>td").Parent().Each(func(_ int, s *goquery.Selection) {
		// 	fmt.Println(s)
		// 	symbol := s.Find("td").Text()
		// 	fmt.Println(symbol)
		// })
		rate = append(rate, e.Text)
	})

	// c.OnHTML("table.wikitable", func(e *colly.HTMLElement) {
	// 	e.DOM.Children().First().Find("tr>td").Parent().Each(func(_ int, s *goquery.Selection) {
	// 		symbol := s.Find("td a").First().Text()
	// 		fmt.Println(symbol)
	// 	})
	// })

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
		fmt.Println("rate slice", rate)
		fmt.Println(len(rate), rate[len(rate) - 3], rate[len(rate) - 2], rate[len(rate) - 1])

		layout := "2006/01/02 15:04:05"
		formatLayout := "2006-01-02 15:04:05"
		rateTimeString := rate[len(rate) - 3]
		rateTime, _ := time.Parse(layout, rateTimeString)
		rateTimeFormat := rateTime.Format(formatLayout)
		rateTimeFormatTime, _ := time.Parse(formatLayout, rateTimeFormat)
		fmt.Println(rateTimeFormat)
		database.GetDB().Create(&models.Rate{
			RateTime: rateTimeFormatTime,
			BuyRate: rate[len(rate) - 2],
			SellRate: rate[len(rate) - 1],
			CrawlFrom: "taiwan-bank",
			CurrencyID: index,
		})
	})

	c.Visit("https://rate.bot.com.tw/xrt/quote/day/" + currency)
	// c.Visit("https://en.wikipedia.org/wiki/List_of_S&P_500_companies")
}