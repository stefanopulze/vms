package handler

import (
	"net/http"
	"vms-core/internal/api/response"
	"vms-core/internal/serial"
)

func NewStatus(port serial.Serial) *StatusHandler {
	return &StatusHandler{
		port: port,
	}
}

type StatusHandler struct {
	port serial.Serial
}

func (sh StatusHandler) Health(w http.ResponseWriter, _ *http.Request) {
	response.Header(w, http.StatusOK)
}

func (sh StatusHandler) Status(w http.ResponseWriter, _ *http.Request) {
	data := map[string]any{
		"status": "ok",
	}

	if qp, ok := sh.port.(*serial.Queue); ok {
		data["queue"] = qp.QueueLength()
	}

	response.Json(w, http.StatusOK, data)
}
