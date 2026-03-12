package internalhttp

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/app"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/models"
)

type EventsHandler struct {
	host        string
	port        int
	logger      *logger.Logger
	application *app.App
}

func NewEventsHandler(logger *logger.Logger, app *app.App, host string, port int) *EventsHandler {
	return &EventsHandler{
		host:        host,
		port:        port,
		logger:      logger,
		application: app,
	}
}

// GetEvents обрабатывает GET /events?date=...&period=...
func (h *EventsHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из query
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		http.Error(w, "date parameter is required", http.StatusBadRequest)
		return
	}
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		http.Error(w, "invalid date format, use RFC3339", http.StatusBadRequest)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "day" // значение по умолчанию
	}

	var events []*models.Event
	switch period {
	case "day":
		events, err = h.application.FindEventsByDay(r.Context(), date)
	case "week":
		events, err = h.application.FindEventsByWeek(r.Context(), date)
	case "month":
		events, err = h.application.FindEventsByMonth(r.Context(), date)
	default:
		http.Error(w, "unsupported period, use day/week/month", http.StatusBadRequest)
		return
	}

	if err != nil {
		h.logger.Error("failed to get events: " + err.Error())
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		return
	}

	// Формируем ответ. Можно использовать свою структуру, если api.Event не подходит.
	type EventResponse struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
		UserID    string    `json:"user_id"`
	}
	resp := make([]EventResponse, 0, len(events))
	for _, event := range events {
		resp = append(resp, EventResponse{
			ID:        event.ID,
			Title:     event.Title,
			StartTime: event.StartTime,
			EndTime:   event.EndTime,
			UserID:    event.UserID,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("failed to write response: " + err.Error())
	}
}

// PostEvents создаёт новое событие
func (h *EventsHandler) PostEvents(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
		UserID      string    `json:"user_id"`
		Reminder    time.Time `json:"reminder"`
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	event := &models.Event{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		UserID:      req.UserID,
		Reminder:    req.Reminder,
		// ID генерируется в storage
	}

	if err := h.application.CreateEvent(r.Context(), event); err != nil {
		h.logger.Error("failed to create event: " + err.Error())
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"id": event.ID}); err != nil {
		h.logger.Error("failed to encode response: " + err.Error())
	}
}

// PutEventsID обновляет событие по ID
func (h *EventsHandler) PutEventsID(w http.ResponseWriter, r *http.Request, id int64) {
	var req struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
		UserID      string    `json:"user_id"`
		Reminder    time.Time `json:"reminder"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	event := &models.Event{
		ID:          strconv.FormatInt(id, 10),
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		UserID:      req.UserID,
		Reminder:    req.Reminder,
	}

	if err := h.application.UpdateEvent(r.Context(), event); err != nil {
		h.logger.Error("failed to update event: " + err.Error())
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteEventsID удаляет событие по ID
func (h *EventsHandler) DeleteEventsID(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.application.DeleteEvent(r.Context(), strconv.FormatInt(id, 10)); err != nil {
		h.logger.Error("failed to delete event: " + err.Error())
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
