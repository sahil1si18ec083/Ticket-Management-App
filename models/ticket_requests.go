package models

import "time"

type TicketRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
}
type TicketUpdateRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Status  *string `json:"status"`
}
type TicketResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
