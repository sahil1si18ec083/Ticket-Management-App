package models

import "time"

type TicketRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
type TicketUpdateRequest struct {
	Title           *string `json:"title"`
	Content         *string `json:"content"`
	Status          *Status `json:"status"`
	AssignedAgentID *uint   `json:"assigned_agent_id"`
}
type TicketResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type TicketAssignRequest struct {
	AssignedAgentID *uint `json:"assigned_agent_id"`
}
