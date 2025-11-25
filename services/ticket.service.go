package services

import (
	"errors"
	"fmt"
	"strconv"
	"ticket-app-gin-golang/models"
	"ticket-app-gin-golang/repositories"
)

type TicketService struct {
	repo *repositories.TicketRepository
}

func NewTicketService(repo *repositories.TicketRepository) *TicketService {
	return &TicketService{repo: repo}
}

// --------------- Create Ticket ---------------
func (s *TicketService) CreateTicket(userID uint, req models.TicketRequest) (*models.Ticket, error) {

	if req.Title == "" {
		return nil, errors.New("title is required")
	}

	ticket := models.Ticket{
		Title:   req.Title,
		Content: req.Content,
		Status:  req.Status,
		UserID:  userID,
	}

	err := s.repo.Create(&ticket)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &ticket, nil
}

// --------------- Get Tickets ---------------
func (s *TicketService) GetUserTickets(userID uint) ([]models.Ticket, error) {
	return s.repo.GetUserTickets(userID)
}

// --------------- Get By ID ---------------
func (s *TicketService) GetTicketByID(userID uint, idStr string) (*models.Ticket, error) {

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, errors.New("invalid ticket id")
	}

	return s.repo.GetByID(userID, uint(id))
}

// --------------- Update Ticket ---------------
func (s *TicketService) UpdateTicketByID(userID uint, idStr string, req models.TicketUpdateRequest) (*models.Ticket, error) {

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	ticket, err := s.repo.GetByID(userID, uint(id))
	if err != nil {
		return nil, errors.New("ticket not found")
	}

	if req.Title != nil {
		ticket.Title = *req.Title
	}
	if req.Content != nil {
		ticket.Content = *req.Content
	}
	if req.Status != nil {
		ticket.Status = *req.Status
	}

	err = s.repo.Update(ticket)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

// --------------- Delete Ticket ---------------
func (s *TicketService) DeleteTicketByID(userID uint, idStr string) error {

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.New("invalid id")
	}

	ticket, err := s.repo.GetByID(userID, uint(id))
	if err != nil {
		return errors.New("ticket not found")
	}

	return s.repo.Delete(ticket)
}
