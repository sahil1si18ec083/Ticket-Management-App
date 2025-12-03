package services

import (
	"errors"
	"fmt"
	"strconv"
	"ticket-app-gin-golang/models"
	"ticket-app-gin-golang/repositories"
)

type TicketService struct {
	repo     *repositories.TicketRepository
	userRepo *repositories.UserRepository
}

func NewTicketService(repo *repositories.TicketRepository, userRepo *repositories.UserRepository) *TicketService {
	return &TicketService{
		repo:     repo,
		userRepo: userRepo,
	}
}

// --------------- Create Ticket ---------------
func (s *TicketService) CreateTicket(
	userID uint,
	req models.TicketRequest,
) (*models.Ticket, error) {

	if req.Title == "" {
		return nil, errors.New("title is required")
	}

	ticket := models.Ticket{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := s.repo.Create(&ticket); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &ticket, nil
}

// --------------- Get Tickets ---------------
func (s *TicketService) GetUserTickets(
	userID uint,
	role string,
) ([]models.Ticket, error) {

	if role == "USER" {
		return s.repo.GetUserTickets(userID)
	}

	return s.repo.GetAllTickets()
}

// --------------- Get By ID AND USERID ---------------
func (s *TicketService) GetTicketByID(
	userID uint,
	idStr string,
	role string,
) (*models.Ticket, error) {

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, errors.New("invalid ticket id")
	}

	if role == "USER" {
		return s.repo.GetByID(userID, uint(id))
	}

	return s.repo.GetOnlyByID(uint(id))
}
func applyWorkflowRules(from_status models.Status, to_status models.Status, ticket *models.Ticket, userID uint) {
	if from_status == models.StatusNew && to_status == models.StatusInProgress {
		if ticket.AssignedAgentID == nil {
			ticket.AssignedAgentID = &userID
		}
	}
}

// --------------- Update Ticket ---------------
func (s *TicketService) UpdateTicketByID(
	userID uint,
	idStr string,
	req models.TicketUpdateRequest,
	role string,
) (*models.Ticket, error) {
	if role == "USER" {
		return nil, errors.New("not authorized to update  ticket")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	ticket, err := s.repo.GetOnlyByID(uint(id))

	if err != nil {
		return nil, errors.New("ticket not found")
	}
	if ticket.Status == models.StatusClosed {
		return nil, errors.New("cannot update closed ticket")
	}
	if req.Status != nil {
		value, exists := models.StatusMap[string(*req.Status)]
		if !exists {
			return nil, errors.New("invalid status,allowed values are NEW,IN_PROGRESS,WAITING,RESOLVED,CLOSED")
		}
		from_status := ticket.Status
		to_status := value
		if !models.IsValidTransition(from_status, to_status) {
			return nil, errors.New("invalid transition")

		}

		if ticket.Status != value {
			// add applyWorkflowRules here

			applyWorkflowRules(from_status, to_status, ticket, userID)

			ticket.Status = value
		}

	}

	if req.Title != nil {
		ticket.Title = *req.Title
	}

	if req.Content != nil {
		ticket.Content = *req.Content
	}
	fmt.Println(ticket)
	if err := s.repo.Update(ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

// --------------- Delete Ticket ---------------
func (s *TicketService) DeleteTicketByID(
	userID uint,
	idStr string,
	role string,
) error {

	if role == "USER" {
		return errors.New("not authorized to delete ticket")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.New("invalid id")
	}

	ticket, err := s.repo.GetOnlyByID(uint(id))
	if err != nil {
		return errors.New("ticket not found")
	}

	return s.repo.Delete(ticket)
}

func (s *TicketService) UnAssignTicket(userID uint,
	idStr string,
	role string) error {
	if role == "USER" {
		return errors.New("not authorized to unassign ticket")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.New("invalid id")
	}

	ticket, err := s.repo.GetOnlyByID(uint(id))
	if err != nil {
		return errors.New("ticket not found")
	}
	if ticket.Status == models.StatusClosed {
		return errors.New("cannot unassign closed ticket")
	}
	if ticket.AssignedAgentID == nil {
		return errors.New("ticket is already unassigned")
	}
	ticket.AssignedAgentID = nil

	return s.repo.Update(ticket)

}

func (s *TicketService) AssignTicket(userID uint,
	idStr string,
	role string, assigned_agent_id uint) error {
	if role == "USER" {
		return errors.New("not authorized to assign ticket")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.New("invalid id")
	}

	ticket, err := s.repo.GetOnlyByID(uint(id))
	if err != nil {
		return errors.New("ticket not found")
	}
	user, err := s.userRepo.FindById(assigned_agent_id)
	if err != nil {
		// do something
		fmt.Println(user)
		return errors.New("assigned id not found")

	}

	if string(user.Role) == string(models.RoleUser) {

		return errors.New("assigned user must be AGENT or ADMIN")

	}
	if ticket.Status == models.StatusClosed {
		return errors.New("cannot assign closed ticket")
	}
	if string(user.Role) == string(models.RoleAgent) && user.ID != userID {
		return errors.New("agent → can assign only to self")
	}
	if ticket.AssignedAgentID != nil && *ticket.AssignedAgentID == assigned_agent_id {
		return errors.New("ticket already assigned to this agent")
	}

	ticket.AssignedAgentID = &assigned_agent_id

	return s.repo.Update(ticket)

}
