package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"xray-stats-telegram/models"
	"xray-stats-telegram/stats"

	"github.com/joho/godotenv"

	"github.com/go-telegram/bot"
	tgModels "github.com/go-telegram/bot/models"
)

var userState *models.UserState

var statsParser *stats.StatsParser

func main() {
	godotenv.Load()

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatalln("env var BOT_TOKEN is not set. Use .env file or env var.")
	}

	userState = models.NewStateFromConfigs(
		"/usr/local/etc/xray-stats-telegram/admins",
		"/usr/local/etc/xray-stats-telegram/users",
	)

	trafficDataDirectory, err := os.ReadFile("/usr/local/etc/xray-stats/directory")
	if err != nil {
		panic(err)
	}

	statsParser = stats.New(strings.TrimSpace(string(trafficDataDirectory)))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		bot.WithMessageTextHandler("/all", bot.MatchTypePrefix, allHandler),
		bot.WithMessageTextHandler("/stats", bot.MatchTypePrefix, statsHandler),
	}

	b, err := bot.New(botToken, opts...)
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

	var builder strings.Builder
	builder.WriteString("Today:\n")

	for _, xrayUser := range *allUserEmails {
		stats := statsParser.GetToday(xrayUser)
		fmt.Fprintf(&builder, "%s\n%s\n\n", xrayUser, stats.ToOneLineString())
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   builder.String(),
	})
}

func statsHandler(ctx context.Context, b *bot.Bot, update *tgModels.Update) {
	userId := update.Message.From.ID
	xrayUser, isXrayUser := userState.GetXrayEmail(userId)

	if !isXrayUser {
		return
	}

	stats := statsParser.GetToday(xrayUser)
	text := "Today:\n" + stats.ToString()

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
}
