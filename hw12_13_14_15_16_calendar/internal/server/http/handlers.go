package internalhttp

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/app"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/models"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/server/http/api"
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

func (h *EventsHandler) GetEvents(w http.ResponseWriter, r *http.Request, params api.GetEventsParams) {
	var (
		events []*models.Event
		err    error
	)
	date := params.Date.Time

	period := api.Day
	if params.Period != nil {
		period = *params.Period
	}

	switch period {
	case api.Day:
		events, err = h.application.FindEventsByDay(r.Context(), date)
	case api.Week:
		events, err = h.application.FindEventsByWeek(r.Context(), date)
	case api.Month:
		events, err = h.application.FindEventsByMonth(r.Context(), date)
	default:
		h.logger.Error("unsupported period value")
		http.Error(w, "Unsupported period", http.StatusInternalServerError)
		return
	}

	if err != nil {
		h.logger.Error("failed to get events: " + err.Error())
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		return
	}

	resp := make([]api.Event, 0, len(events))
	for _, event := range events {
		resp = append(resp, api.Event{
			Id:       event.ID,
			Title:    event.Title,
			DateTime: event.DateTime,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("failed to write response: " + err.Error())
	}
}

func (h *EventsHandler) PostEvents(w http.ResponseWriter, r *http.Request) {
	var req api.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	event := &models.Event{
		ID:       req.Id,
		Title:    req.Title,
		DateTime: time.Now(), // временно, пока не разберемся с полем
	}

	if err := h.application.CreateEvent(r.Context(), event); err != nil {
		h.logger.Error("failed to create event: " + err.Error())
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *EventsHandler) PutEventsID(w http.ResponseWriter, r *http.Request, id int64) {
	var req api.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	event := &models.Event{
		ID:       id,
		Title:    req.Title,
		DateTime: req.DateTime,
	}

	if err := h.application.Update(r.Context(), strconv.FormatInt(id, 10), event); err != nil {
		h.logger.Error("failed to update event: " + err.Error())
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *EventsHandler) DeleteEventsID(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.application.DeleteByID(r.Context(), strconv.FormatInt(id, 10)); err != nil {
		h.logger.Error("failed to delete event: " + err.Error())
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
