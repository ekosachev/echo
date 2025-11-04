package service

import (
	"context"

	"github.com/ekosachev/go-backend-template/internal/models"
	"github.com/ekosachev/go-backend-template/internal/repository"
)

type EventService struct {
	repo repository.Repository[models.Event]
}

func NewEventService(repo repository.Repository[models.Event]) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) CreateEvent(ctx context.Context, event *models.Event) error {
	return s.repo.Create(ctx, event)
}

func (s *EventService) GetEventByID(ctx context.Context, id any) (*models.Event, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *EventService) UpdateEvent(ctx context.Context, event *models.Event) error {
	return s.repo.Update(ctx, event)
}

func (s *EventService) DeleteEvent(ctx context.Context, id any) error {
	where := map[string]any{"id": id}
	return s.repo.Delete(ctx, where)
}
