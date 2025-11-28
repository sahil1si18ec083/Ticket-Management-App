package repositories

import (
	"ticket-app-gin-golang/models"

	"gorm.io/gorm"
)

type TicketRepository struct {
	DB *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{
		DB: db,
	}
}

func (r *TicketRepository) Create(ticket *models.Ticket) error {
	return r.DB.Create(ticket).Error
}

func (r *TicketRepository) GetUserTickets(userID uint) ([]models.Ticket, error) {
	var tickets []models.Ticket

	err := r.DB.
		Where("user_id = ?", userID).
		Find(&tickets).Error

	return tickets, err
}

func (r *TicketRepository) GetAllTickets() ([]models.Ticket, error) {
	var tickets []models.Ticket

	err := r.DB.Find(&tickets).Error
	return tickets, err
}

func (r *TicketRepository) GetByID(userID uint, ticketID uint) (*models.Ticket, error) {
	var ticket models.Ticket

	err := r.DB.
		Where("user_id = ? AND id = ?", userID, ticketID).
		First(&ticket).Error

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *TicketRepository) GetOnlyByID(ticketID uint) (*models.Ticket, error) {
	var ticket models.Ticket

	err := r.DB.
		Where("id = ?", ticketID).
		First(&ticket).Error

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *TicketRepository) Update(ticket *models.Ticket) error {
	return r.DB.Save(ticket).Error
}

func (r *TicketRepository) Delete(ticket *models.Ticket) error {
	return r.DB.Delete(ticket).Error
}
