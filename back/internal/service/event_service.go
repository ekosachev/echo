package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ekosachev/go-backend-template/internal/models"
	"github.com/ekosachev/go-backend-template/internal/repository"
	"github.com/ekosachev/go-backend-template/pkg/timezone"
	"github.com/google/uuid"
)

type EventService struct {
	eventRepo    repository.Repository[models.Event]
	userRepo     repository.Repository[models.User]
	calendarRepo repository.Repository[models.Calendar]
}

func NewEventService(
	eventRepo repository.Repository[models.Event],
	userRepo repository.Repository[models.User],
	calendarRepo repository.Repository[models.Calendar],
) *EventService {
	return &EventService{
		eventRepo:    eventRepo,
		userRepo:     userRepo,
		calendarRepo: calendarRepo,
	}
}

type CreateEventRequest struct {
	CalendarID  uuid.UUID `json:"calendar_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	Timezone    string    `json:"timezone"`
	AllDay      bool      `json:"all_day"`
	Priority    int       `json:"priority"`
	Color       string    `json:"color"`
}

func (s *EventService) CreateEvent(ctx context.Context, userID uuid.UUID, req *CreateEventRequest) (*models.Event, error) {
	if !timezone.IsValidTimezone(req.Timezone) {
		return nil, fmt.Errorf("invalid timezone: %s", req.Timezone)
	}

	calendar, err := s.calendarRepo.FindByID(ctx, req.CalendarID)
	if err != nil {
		return nil, fmt.Errorf("calendar not found: %w", err)
	}
	if calendar.OwnerID != userID {
		return nil, fmt.Errorf("access denied")
	}

	startUTC, err := timezone.ConvertToUTC(req.StartAt, req.Timezone)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to UTC start time: %w", err)
	}

	endUTC, err := timezone.ConvertFromUTC(req.EndAt, req.Timezone)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to UTC end time: %w", err)
	}

	event := &models.Event{
		CalendarId:  req.CalendarID,
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		StartAt:     startUTC,
		EndAt:       endUTC,
		Timezone:    req.Timezone,
		AllDay:      req.AllDay,
		Priority:    req.Priority,
		Color:       req.Color,
		Meta:        make(map[string]interface{}),
	}

	if err := s.eventRepo.Create(ctx, event); err != nil {
		return nil, err
	}

	return event, nil
}

type EventResponse struct {
	CalendarID  uuid.UUID              `json:"calendar_id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Location    string                 `json:"location"`
	StartAt     time.Time              `json:"start_at"`
	EndAt       time.Time              `json:"end_at"`
	Timezone    string                 `json:"timezone"`
	AllDay      bool                   `json:"all_day"`
	Priority    int                    `json:"priority"`
	Color       string                 `json:"color"`
	Meta        map[string]interface{} `json:"meta"`
}

func (s *EventService) GetEventForUser(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) (*models.Event, error) {
	event, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userTimezone := timezone.GetUserTimezone(user.Timezone)
	localEvent := *event
	localEvent.StartAt, _ = timezone.ConvertFromUTC(event.StartAt, userTimezone)
	localEvent.EndAt, _ = timezone.ConvertFromUTC(event.EndAt, userTimezone)

	return &localEvent, nil
}

func (s *EventService) UpdateEventWithTimezone(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, req *CreateEventRequest,
) (*models.Event, error) {
	event, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	calendar, err := s.calendarRepo.FindByID(ctx, req.CalendarID)
	if err != nil {
		return nil, fmt.Errorf("calendar not found: %w", err)
	}
	if calendar.OwnerID != userID {
		return nil, fmt.Errorf("invalid timezone: %s", req.Timezone)
	}
	if !timezone.IsValidTimezone(req.Timezone) {
		return nil, fmt.Errorf("invalid timezone: %s", req.Timezone)
	}

	startUTC, err := timezone.ConvertToUTC(req.StartAt, req.Timezone)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to UTC start time: %w", err)
	}

	endUTC, err := timezone.ConvertFromUTC(req.EndAt, req.Timezone)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to UTC end time: %w", err)
	}

	event.CalendarId = req.CalendarID
	event.Title = req.Title
	event.Description = req.Description
	event.Location = req.Location
	event.StartAt = startUTC
	event.EndAt = endUTC
	event.Timezone = req.Timezone
	event.AllDay = req.AllDay
	event.Priority = req.Priority
	event.Color = req.Color
	event.Meta = make(map[string]interface{})

	if err := s.eventRepo.Update(ctx, event); err != nil {
		return nil, err
	}

	return event, nil
}

func (s *EventService) GetEventsByCalendar(ctx context.Context, calendarID uuid.UUID, userID uuid.UUID) ([]*models.Event, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userTimezone := timezone.GetUserTimezone(user.Timezone)
	events := []*models.Event{}

	for _, event := range events {
		event.StartAt, _ = timezone.ConvertFromUTC(event.StartAt, userTimezone)
		event.EndAt, _ = timezone.ConvertFromUTC(event.EndAt, userTimezone)
	}

	return events, nil
}

func (s *EventService) GetEventByID(ctx context.Context, id any) (*models.Event, error) {
	return s.eventRepo.FindByID(ctx, id)
}

func (s *EventService) UpdateEvent(ctx context.Context, event *models.Event) error {
	return s.eventRepo.Update(ctx, event)
}

func (s *EventService) DeleteEvent(ctx context.Context, id any) error {
	where := map[string]any{"id": id}
	return s.eventRepo.Delete(ctx, where)
}
