package service

import (
	"context"
	"fmt"

	"github.com/ekosachev/go-backend-template/internal/models"
	"github.com/ekosachev/go-backend-template/internal/repository"
	"github.com/google/uuid"
)

type CalendarService struct {
	repo repository.Repository[models.Calendar]
}

func NewCalendarService(repo repository.Repository[models.Calendar]) *CalendarService {
	return &CalendarService{repo: repo}
}

func (s *CalendarService) CreateCalendar(ctx context.Context, calendar *models.Calendar) error {
	if calendar.Name == "" {
		return fmt.Errorf("calendar name is required")
	}
	if calendar.OwnerID == uuid.Nil {
		return fmt.Errorf("owner ID is required")
	}
	return s.repo.Create(ctx, calendar)
}

func (s *CalendarService) GetUserCalendars(ctx context.Context, userID uuid.UUID) ([]*models.Calendar, error) {
	where := map[string]any{"owner_id": userID}
	calendar, err := s.repo.FindOne(ctx, where)
	if err != nil {
		return []*models.Calendar{}, nil
	}
	return []*models.Calendar{calendar}, nil
}

func (s *CalendarService) GetCalendarForUser(ctx context.Context, calendarID, userID uuid.UUID) (*models.Calendar, error) {
	calendar, err := s.repo.FindByID(ctx, calendarID)
	if err != nil {
		return nil, fmt.Errorf("calendar not found: %w", err)
	}

	if calendar.OwnerID != userID {
		return nil, fmt.Errorf("access denied")
	}

	return calendar, nil
}

func (s *CalendarService) GetCalendarByID(ctx context.Context, id any) (*models.Calendar, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *CalendarService) UpdateCalendar(ctx context.Context, calendar *models.Calendar) error {
	if calendar.Name == "" {
		return fmt.Errorf("calendar name is required")
	}

	return s.repo.Update(ctx, calendar)
}

func (s *CalendarService) DeleteCalendar(ctx context.Context, id any) error {
	where := map[string]any{"id": id}
	return s.repo.Delete(ctx, where)
}

func (s *CalendarService) CalendarExists(ctx context.Context, calendarID, userID uuid.UUID) (bool, error) {
	calendar, err := s.repo.FindByID(ctx, calendarID)
	if err != nil {
		return false, nil
	}

	return calendar.OwnerID == userID, nil
}
