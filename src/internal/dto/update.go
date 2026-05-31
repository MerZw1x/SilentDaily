package dto

type SubmitUpdateRequest struct {
	TelegramUserID int
	RawText        string
}

type SubmitUpdateResponse struct {
	Status string
}
