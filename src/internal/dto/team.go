package dto

type CreateTeamRequest struct {
	Name string `json:"name" validate:"required"`
}

type CreateTeamResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
