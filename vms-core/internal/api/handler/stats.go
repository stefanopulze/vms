package handler

import (
	"net/http"
	"time"
	"vms-core/internal/api/response"
	"vms-core/internal/infrastructure/exporter/influx"
	"vms-core/internal/service"
)

func NewStats(ic *influx.Client, service *service.Stats) *StatsHandler {
	return &StatsHandler{
		influx:  ic,
		service: service,
	}
}

type StatsHandler struct {
	influx  *influx.Client
	service *service.Stats
}

func (h StatsHandler) GetDayStats(w http.ResponseWriter, r *http.Request) {
	day := r.URL.Query().Get("day")
	if day == "" || len(day) < 1 {
		response.Error(w, r, http.StatusBadRequest, "day is required")
		return
	}

	stats, err := h.service.GetDayPowerUsage(r.Context(), day)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	response.Json(w, http.StatusOK, stats)
}

func (h StatsHandler) DownsamplingDayStats(w http.ResponseWriter, r *http.Request) {
	day := r.URL.Query().Get("day")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if day != "" {
		_, err := h.service.DownsamplingDay(r.Context(), day)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}
	} else if from != "" {
		fromTime, err := time.Parse("2006-01-02", from)
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		toTime, err := time.Parse("2006-01-02", to)
		if err != nil {
			toTime = time.Now()
			return
		}

		for t := fromTime; t.Before(toTime); t = t.AddDate(0, 0, 1) {
			_, err = h.service.DownsamplingDay(r.Context(), t.Format("2006-01-02"))
			if err != nil {
				response.Error(w, r, http.StatusInternalServerError, err)
				return
			}
		}

	} else {
		response.Error(w, r, http.StatusBadRequest, "day or from and to are required")
	}

	response.Header(w, http.StatusAccepted)
}
