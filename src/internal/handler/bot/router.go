package bot

import (
	"fmt"
	"log"
	"strings"

	"silent/src/internal/service/abstract"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	bot           *tgbotapi.BotAPI
	updateService abstract.IUpdateService
	memberService abstract.IMemberService
}

func NewRouter(bot *tgbotapi.BotAPI, updateService abstract.IUpdateService, memberService abstract.IMemberService) *Router {
	return &Router{
		bot:           bot,
		updateService: updateService,
		memberService: memberService,
	}
}

func (r *Router) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := r.bot.GetUpdatesChan(u)

	log.Println("[bot] listening for updates...")
	for update := range updates {
		if update.Message == nil {
			continue
		}
		r.handleMessage(update.Message)
	}
}

func (r *Router) handleMessage(msg *tgbotapi.Message) {
	text := strings.TrimSpace(msg.Text)
	chatID := msg.Chat.ID
	userID := int(msg.From.ID)

	switch {
	case text == "/start":
		r.sendText(chatID, "Привет! Я SilentDaily — твой асинхронный стендап-бот.\n\nВ конце рабочего дня напиши мне что сделал, что планируешь и есть ли блокеры. Утром команда получит дайджест.\n\nКоманды:\n/start — это сообщение\n/help — помощь")

	case text == "/help":
		r.sendText(chatID, "Как пользоваться:\n\n1. Напиши мне в конце дня свой апдейт в свободной форме\n2. Например: «Сегодня закончил авторизацию, завтра займусь тестами, блокеров нет»\n3. Утром в 08:00 тимлид получит дайджест, а ты — свой план на день")

	case strings.HasPrefix(text, "/"):
		r.sendText(chatID, "Неизвестная команда. Напиши /help для помощи.")

	default:
		r.handleDailyUpdate(chatID, userID, text)
	}
}

func (r *Router) handleDailyUpdate(chatID int64, telegramUserID int, text string) {
	if len(strings.TrimSpace(text)) < 10 {
		r.sendText(chatID, "Сообщение слишком короткое. Опиши подробнее что сделал сегодня.")
		return
	}

	err := r.updateService.Submit(telegramUserID, text)
	if err != nil {
		if err.Error() == "member not registered" {
			r.sendText(chatID, fmt.Sprintf(
				"Ты не зарегистрирован в системе.\n\nПопроси тимлида добавить тебя через API:\nPOST /api/v1/members\n{\"telegram_user_id\": %d, \"name\": \"Твоё имя\", \"team_id\": 1, \"is_lead\": false}",
				telegramUserID,
			))
			return
		}
		log.Printf("[bot] submit error for user %d: %v", telegramUserID, err)
		r.sendText(chatID, "Произошла ошибка при сохранении. Попробуй ещё раз.")
		return
	}

	r.sendText(chatID, "Апдейт принят! Обрабатываю его с помощью AI...\n\nУтром в 08:00 тимлид получит дайджест команды, а ты — свой план на день.")
}

func (r *Router) sendText(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := r.bot.Send(msg); err != nil {
		log.Printf("[bot] send error to %d: %v", chatID, err)
	}
}
