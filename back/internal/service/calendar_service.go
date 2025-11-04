package service

import (
	"context"

	"github.com/ekosachev/go-backend-template/internal/models"
	"github.com/ekosachev/go-backend-template/internal/repository"
)

type CalendarService struct {
	repo repository.Repository[models.Calendar]
}

func NewCalendarService(repo repository.Repository[models.Calendar]) *CalendarService {
	return &CalendarService{repo: repo}
}

func (s *CalendarService) CreateCalendar(ctx context.Context, calendar *models.Calendar) error {
	return s.repo.Create(ctx, calendar)
}

func (s *CalendarService) GetCalendarByID(ctx context.Context, id any) (*models.Calendar, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *CalendarService) UpdateCalendar(ctx context.Context, calendar *models.Calendar) error {
	return s.repo.Update(ctx, calendar)
}

func (s *CalendarService) DeleteCalendar(ctx context.Context, id any) error {
	where := map[string]any{"id": id}
	return s.repo.Delete(ctx, where)
}
