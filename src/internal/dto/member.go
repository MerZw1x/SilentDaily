package dto

type RegisterMemberRequest struct {
	TelegramUserID int    `json:"telegram_user_id" validate:"required"`
	Name           string `json:"name" validate:"required"`
	TeamID         int    `json:"team_id" validate:"required"`
	IsLead         bool   `json:"is_lead"`
}

type RegisterMemberResponse struct {
	Status string `json:"status"`
}
