package main

import (
	// "fmt"
	"fmt"
	"log"

	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gocolly/colly"
	"github.com/jjjjackson/ps5-bot/config"
)

var TARGET_URL = "https://www.amazon.co.jp/%E3%82%BD%E3%83%8B%E3%83%BC%E3%83%BB%E3%82%A4%E3%83%B3%E3%82%BF%E3%83%A9%E3%82%AF%E3%83%86%E3%82%A3%E3%83%96%E3%82%A8%E3%83%B3%E3%82%BF%E3%83%86%E3%82%A4%E3%83%B3%E3%83%A1%E3%83%B3%E3%83%88-PlayStation-5-CFI-1000A01/dp/B08GGGBKRQ/ref=sr_1_1?dchild=1&keywords=playstation+5&qid=1617209269&sr=8-1"

func isPS5InStock() bool {
	log.Printf("run isPS5InStock")

	c := colly.NewCollector()

	if err := c.Visit(TARGET_URL); err != nil {
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

	status := isPS5InStock()
	txt := fmt.Sprintf("PS5 In stock : \n %s", TARGET_URL)
	if !status {
		txt = fmt.Sprintf("PS5 Not In stock \n")
	}

	mid, err := strconv.ParseInt(config.MessageID, 10, 64)
	if err != nil {
		log.Fatal("Couldn't parse message ID")
	}

	msg := tgbot.NewMessage(mid, txt)
	_, err = bot.Send(msg)
	if err != nil {
		log.Fatal("Couldn't send message")
	}
}

func printConfig(config *config.Config) {
	log.Println(config.MessageID)
}

func handleLambdaStart() {
	config := config.LoadConfig()
	sendTelegramMessage(config)
}

func main() {
	lambda.Start(handleLambdaStart)
}
