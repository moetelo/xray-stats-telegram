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

	allUserEmails := userState.GetAllUsers()

	userStatsSorted := make([]stats.Stats, 0, len(allUserEmails))
	emptyStatsUsers := make([]string, 0)
	for _, xrayUser := range allUserEmails {
		stats := statsParser.QueryUser(xrayUser, time.Now())

		if stats.Down == 0 && stats.Up == 0 {
			emptyStatsUsers = append(emptyStatsUsers, stats.User)
			continue
		}

		userStatsSorted = append(userStatsSorted, *stats)
	}

	var builder strings.Builder
	for _, stats := range userStatsSorted {
		fmt.Fprintf(&builder, "%s\n%s\n\n", stats.User, stats.ToOneLineString())
	}

	fmt.Fprintln(&builder, "Empty stats users:", strings.Join(emptyStatsUsers, ", "))

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

	stats := statsParser.QueryUser(xrayUser, time.Now())

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   stats.ToString(),
	})
}
