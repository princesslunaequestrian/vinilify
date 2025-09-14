package main

import (
	"context"
	"log"
	"os"

	tg "github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/princesslunaequestrian/vinilify/types"
)

// Global hashmap for storing user requests
var users map[int64]*types.User

func main() {

	ctx := context.Background()

	// Initializing the map of the users
	users = make(map[int64]*types.User)

	// Telegram API token
	botToken, err := getToken()
	if err != nil {
		log.Fatal("Failed to find token file:", err)
	}

	// Initializing the bot
	bot, err := tg.NewBot(botToken, tg.WithDefaultDebugLogger())
	if err != nil {
		log.Fatal(err)
	}

	// Launching updates and handler
	updates, err := bot.UpdatesViaLongPolling(ctx, &tg.GetUpdatesParams{Timeout: 10})
	if err != nil {
		log.Panic(err)
	}

	// Launching bot handler
	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		bh.Stop()
		//ctx, cancel := context.WithCancel(context.Background())
	}()

	bh.Handle(
		handleStart,
		th.CommandEqual("start"),
	)

	bh.Handle(
		handleGenerateVideo,
		th.CommandEqual("generate"),
	)

	bh.Handle(
		handleUpload,
		th.Any(),
	)

	bh.Start()

}

func getToken() (string, error) {
	dat, err := os.ReadFile("token")
	return string(dat), err
}
