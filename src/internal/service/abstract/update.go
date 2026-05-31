package abstract

type IUpdateService interface {
	Submit(telegramUserID int, rawText string) error
}
