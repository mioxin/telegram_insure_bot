package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/commands"
	"github.com/mrmioxin/gak_telegram_bot/internal/config"
	"github.com/mrmioxin/gak_telegram_bot/internal/httpclient"
)

func main() {

	conf, options, err := config.Load(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	if !options.Verbose {
		outputLog, err := os.OpenFile(conf.LogFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatal("Cant open output file for logging ", conf.LogFile, ".\n", err)
		}
		defer outputLog.Close()
		log.SetOutput(outputLog)
	}

	bot, err := tgapi.NewBotAPIWithClient(conf.Token, httpclient.NewHTTPClient())
	if err != nil {
		log.Fatal("Telegram API connection failed: ", err)
	}

	bot.Debug = options.Debug
	log.Printf("Authorized on account %s\n", bot.Self.UserName)
	chSignal := make(chan os.Signal, 1)

	// isModifyConfig := make(chan any)
	// go conf.Watch(resources.CONFIG_FILE_NAME, resources.DURATION_WATCH_CONFIG, isModifyConfig)

	c := commands.NewCommander(bot, conf) //, isModifyConfig)

	signal.Notify(chSignal, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(chSignal)
	go func() {
		<-chSignal
		c.Stop()
		log.Println("Cancel bot by OS interruption.")
	}()

	c.Start()
}
