package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"xray-stats-telegram/models"
	"xray-stats-telegram/queryDate"
	"xray-stats-telegram/stats"

	"github.com/joho/godotenv"

	"github.com/go-telegram/bot"
	tgModels "github.com/go-telegram/bot/models"
)

var userState *models.UserState

var statsParser *stats.StatsParser

const (
	CommandAdminAll = "/all"
	CommandQuery    = "/query"
)

const (
	HelpCommandAll   = "Get stats for every user for a specific date. Usage: /all [YYYY-MM-DD/empty for today]"
	HelpCommandQuery = "Get your stats for a specific date. Usage: /query [YYYY-MM-DD/empty for today]"
)

func main() {
	godotenv.Load()

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatalln("env var BOT_TOKEN is not set. Use .env file or env var.")
	}

	statsQueryBin := os.Getenv("STATS_QUERY_BIN")
	statsParser = stats.New(statsQueryBin)

	userState = models.NewState(
		"/usr/local/etc/xray-stats-telegram/admins",
		"/usr/local/etc/xray-stats-telegram/users",
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(helpHandler),
		bot.WithMessageTextHandler(CommandAdminAll, bot.MatchTypePrefix, allHandler),
		bot.WithMessageTextHandler(CommandQuery, bot.MatchTypePrefix, queryHandler),
		bot.WithAllowedUpdates([]string{"message", "callback_query"}),
		bot.WithCallbackQueryDataHandler("", bot.MatchType(bot.HandlerTypeCallbackQueryData), allKeyboardHandler),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	fmt.Println("[xray-stats-telegram] Polling for messages...")
	b.Start(ctx)
}

func helpHandler(ctx context.Context, b *bot.Bot, update *tgModels.Update) {
	if update.Message == nil {
		return
	}

	userId := update.Message.From.ID
	if userState.IsAdmin(userId) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Commands:\n\n" + CommandAdminAll + "\n" + HelpCommandAll + "\n\n" + CommandQuery + "\n" + HelpCommandQuery,
		})
		return
	}

	_, isXrayUser := userState.GetXrayEmail(userId)
	if !isXrayUser {
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Commands:\n" + CommandQuery + "\n" + HelpCommandQuery,
	})
}

func allHandler(ctx context.Context, b *bot.Bot, update *tgModels.Update) {
	userId := update.Message.From.ID
	if !userState.IsAdmin(userId) {
		return
	}

	date, hasError := parseQueryDateFromMessage(update.Message.Text)
	if hasError {
		handleBadDateMessage(ctx, b, update)
		return
	}

	allStats := statsParser.Query(date)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        stats.StatsArrayToMessageText(date, allStats),
		ReplyMarkup: dateKeyboard(date),
	})
}

func parseQueryDateFromMessage(messageText string) (queryDate.QueryDate, bool) {
	date := queryDate.Now()
	args := strings.Fields(messageText)
	if len(args) >= 2 {
		possiblyDate := args[1]
		parsedQdate, err := queryDate.Parse(possiblyDate)
		if err != nil {
			return queryDate.QueryDate{}, true
		}

		date = parsedQdate
	}
	return date, false
}

func dateKeyboard(date queryDate.QueryDate) *tgModels.InlineKeyboardMarkup {
	return &tgModels.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgModels.InlineKeyboardButton{
			{
				{Text: "⬅️", CallbackData: date.Prev().String()},
				{Text: "🔄", CallbackData: date.String()},
				{Text: "➡️", CallbackData: date.Next().String()},
			},
		},
	}
}

func allKeyboardHandler(ctx context.Context, b *bot.Bot, update *tgModels.Update) {
	cq := update.CallbackQuery
	if !userState.IsAdmin(cq.From.ID) {
		return
	}

	date, err := queryDate.Parse(cq.Data)
	if err != nil {
		return
	}

	allStats := statsParser.Query(date)

	botMessage := cq.Message.Message
	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      botMessage.Chat.ID,
		MessageID:   botMessage.ID,
		Text:        stats.StatsArrayToMessageText(date, allStats),
		ReplyMarkup: dateKeyboard(date),
	})
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: cq.ID,
		ShowAlert:       false,
	})

	fmt.Println(err)
}

func queryHandler(ctx context.Context, b *bot.Bot, update *tgModels.Update) {
	xrayUser, isXrayUser := userState.GetXrayEmail(update.Message.From.ID)
	if !isXrayUser {
		return
	}

	date, hasError := parseQueryDateFromMessage(update.Message.Text)
	if hasError {
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
		Text:   "Please provide a date in the format YYYY-MM-DD or YYYY-MM.",
	})
}
