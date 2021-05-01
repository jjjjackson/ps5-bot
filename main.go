package main

import (
	"fmt"
	"log"

	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gocolly/colly"
	"github.com/jjjjackson/ps5-bot/config"
)

var TARGET_URL = []string{
	"https://www.amazon.co.jp/dp/B08GGGBKRQ/ref=cm_sw_r_oth_api_glt_i_WYS1X7GWPSY386FSVN65",
	"https://www.amazon.co.jp/dp/B08GGGCH3Y/ref=cm_sw_r_oth_api_glt_i_SXG07XY44Y5BX9PZ8Q2M",
	"https://www.amazon.co.jp/dp/B091D2HGKP/ref=cm_sw_r_oth_api_glt_i_JNPERGZM9BE65XCY7J0X",
	"https://www.amazon.co.jp/dp/B091D2959B/ref=cm_sw_r_oth_api_glt_i_K18C0CRF8PYNFQZZP5V5",
}

func isPS5InStock(url string) bool {
	log.Printf("run isPS5InStock")

	c := colly.NewCollector()

	if err := c.Visit(url); err != nil {
		log.Fatal("Couldn't visit amazon")
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36")
	})

	status := false
	c.OnHTML("#add-to-cart-button", func(e *colly.HTMLElement) {
		status = true
	})

	return status
}

func sendTelegramMessage(config *config.Config) {

	bot, err := tgbot.NewBotAPI(config.TelegramToken)
	if err != nil {
		log.Fatal("Couldn't new bot")
	}

	for _, url := range TARGET_URL {
		status := isPS5InStock(url)
		txt := fmt.Sprintf("PS5 In stock : \n %s", TARGET_URL)

		mid, err := strconv.ParseInt(config.MessageID, 10, 64)
		if err != nil {
			log.Fatal("Couldn't parse message ID")
		}

		msg := tgbot.NewMessage(mid, txt)

		if status {
			if _, err := bot.Send(msg); err != nil {
				log.Fatal("Couldn't send message")
			}
		}

	}
}

func handleLambdaStart() {
	config := config.LoadConfig()
	sendTelegramMessage(config)
}

func main() {
	lambda.Start(handleLambdaStart)
}
