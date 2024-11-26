package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
	"xray-stats-telegram/models"
	"xray-stats-telegram/stats"

	"github.com/joho/godotenv"

	"github.com/go-telegram/bot"
	tgModels "github.com/go-telegram/bot/models"
)

var userState *models.UserState

var statsParser *stats.StatsParser

const (
	CommandAll   = "/all"
	CommandQuery = "/query"
)

func main() {
	godotenv.Load()

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatalln("env var BOT_TOKEN is not set. Use .env file or env var.")
	}

	statsQueryBin := os.Getenv("STATS_QUERY_BIN")
	statsParser = stats.New(statsQueryBin)

	userState = models.NewStateFromConfigs(
		"/usr/local/etc/xray-stats-telegram/admins",
		"/usr/local/etc/xray-stats-telegram/users",
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		bot.WithMessageTextHandler(CommandAll, bot.MatchTypePrefix, allHandler),
		bot.WithMessageTextHandler(CommandQuery, bot.MatchTypePrefix, queryHandler),
		bot.WithAllowedUpdates([]string{"message"}),
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
	if userState.IsAdmin(userId) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Command:\n" + CommandAll + "\n" + CommandQuery,
		})
		return
	}

	_, isXrayUser := userState.GetXrayEmail(userId)
	if !isXrayUser {
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Command:\n" + CommandQuery,
	})
}

func allHandler(ctx context.Context, b *bot.Bot, update *tgModels.Update) {
	userId := update.Message.From.ID
	isAdmin := userState.IsAdmin(userId)

	if !isAdmin {
		return
	}

	date, err := parseDate(update.Message.Text)
	if err != nil {
		handleBadDateMessage(ctx, b, update)
		return
	}

	allStats := statsParser.Query(date)

	var builder strings.Builder
	for _, stats := range allStats {
		fmt.Fprintf(&builder, "%s\n%s\n\n", stats.User, stats.ToOneLineString())
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   builder.String(),
	})
}

func queryHandler(ctx context.Context, b *bot.Bot, update *tgModels.Update) {
	xrayUser, isXrayUser := userState.GetXrayEmail(update.Message.From.ID)
	if !isXrayUser {
		return
	}

	date, err := parseDate(update.Message.Text)
	if err != nil {
		handleBadDateMessage(ctx, b, update)
		return
	}

	stats := statsParser.QueryUser(xrayUser, date)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   stats.ToString(),
	})
}

func handleBadDateMessage(ctx context.Context, b *bot.Bot, update *tgModels.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Please provide a date in the format YYYY-MM-DD.",
	})
}

func parseDate(messageText string) (time.Time, error) {
	args := strings.Fields(messageText)
	if len(args) < 2 {
		return time.Now(), nil
	}

	return time.Parse(time.DateOnly, args[1])
}
