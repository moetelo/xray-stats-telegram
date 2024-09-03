package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"xray-stats-telegram/internal"
	"xray-stats-telegram/models"
	"xray-stats-telegram/stats"

	"github.com/alexflint/go-arg"
	"github.com/joho/godotenv"

	"github.com/go-telegram/bot"
	tgModels "github.com/go-telegram/bot/models"
)

type Args struct {
	UsersJsonPath        string `arg:"-u,--users-json,required" help:"Path to the users.json file"`
	TrafficDataDirectory string `arg:"-t,--traffic-data,required" help:"Path to the traffic data directory"`
}

var userState *models.UserState

var statsParser *stats.StatsParser

func main() {
	godotenv.Load()

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatalln("env var BOT_TOKEN is not set. Use .env file or env var.")
	}

	var args Args
	arg.MustParse(&args)

	var userMap models.UsersJson
	err := internal.ReadJson(args.UsersJsonPath, &userMap)
	if err != nil {
		panic(err)
	}

	userState = models.NewState(&userMap)
	statsParser = stats.New(args.TrafficDataDirectory)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		bot.WithMessageTextHandler("/all", bot.MatchTypePrefix, allHandler),
		bot.WithMessageTextHandler("/stats", bot.MatchTypePrefix, statsHandler),
	}

	b, err := bot.New(os.Getenv("BOT_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}

	fmt.Println("[xray-stats-telegram] Polling for messages...")
	b.Start(ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *tgModels.Update) {
	userId := update.Message.From.ID
	isAdmin := userState.IsAdmin(userId)
	_, isXrayUser := userState.GetXrayEmail(userId)

	if isAdmin {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Commands:\n/stats\n/all",
		})
		return
	}

	if !isXrayUser {
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Command:\n/stats",
	})
}

func allHandler(ctx context.Context, b *bot.Bot, update *tgModels.Update) {
	userId := update.Message.From.ID
	isAdmin := userState.IsAdmin(userId)

	if !isAdmin {
		return
	}

	allUserEmails := userState.GetAllUsers()

	builder := strings.Builder{}
	builder.WriteString("Today:\n")

	for _, xrayUser := range *allUserEmails {
		builder.WriteString(xrayUser)
		builder.WriteRune('\n')

		stats := statsParser.GetToday(xrayUser)
		builder.WriteString(stats.ToOneLineString())
		builder.WriteString("\n\n")
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   builder.String(),
	})
}

func statsHandler(ctx context.Context, b *bot.Bot, update *tgModels.Update) {
	userId := update.Message.From.ID
	isAdmin := userState.IsAdmin(userId)
	xrayUser, isXrayUser := userState.GetXrayEmail(userId)

	if !isAdmin && !isXrayUser {
		return
	}

	stats := statsParser.GetToday(xrayUser)
	text := "Today:\n" + stats.ToString()

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
}
