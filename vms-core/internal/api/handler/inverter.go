package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"vms-core/internal/api/model"
	"vms-core/internal/api/response"
	"vms-core/internal/cache"
	"vms-core/internal/event"
	"vms-core/internal/voltronic"
)

const (
	cacheAggregateInfo = "aggregate_info"
	cacheRatingInfo    = "rating_info"
)

func NewInverter(inverter *voltronic.Client, qs *cache.QuerySnapshot) *InverterHandler {
	return &InverterHandler{
		inverter:      inverter,
		cache:         cache.New(),
		querySnapshot: qs,
	}
}

type InverterHandler struct {
	inverter      *voltronic.Client
	cache         *cache.Cache
	event         event.Publisher
	querySnapshot *cache.QuerySnapshot
}

func (h InverterHandler) QueryTime(w http.ResponseWriter, r *http.Request) {
	t, err := h.inverter.QueryTime()
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	data := map[string]any{
		"status": "ok",
		"time":   t,
	}

	response.Json(w, http.StatusOK, data)
}

func (h InverterHandler) UpdateTime(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	if err := h.inverter.UpdateTime(t); err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	response.Header(w, http.StatusAccepted)
}

func (h InverterHandler) QueryRatingInfo(w http.ResponseWriter, r *http.Request) {
	if cached, found := h.cache.Get(cacheRatingInfo); found {
		response.Json(w, http.StatusOK, cached)
		return
	}

	data, err := h.inverter.QueryPIRI()
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	j := model.ApiRatingInfo{
		DeviceRatingInfo:          data,
		SourcePriorityEnum:        data.ChargerSourcePriorityEnum(),
		ChargerSourcePriorityEnum: data.ChargerSourcePriorityEnum(),
	}

	h.cache.Set(cacheRatingInfo, j, 3*time.Second)

	response.Json(w, http.StatusOK, j)
}

func (h InverterHandler) QueryMode(w http.ResponseWriter, r *http.Request) {
	mode, err := h.inverter.QueryMode()
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	data := map[string]any{
		"status": "ok",
		"mode":   mode.Mode,
	}
	response.Json(w, http.StatusOK, data)
}

func (h InverterHandler) QueryWarnings(w http.ResponseWriter, r *http.Request) {
	data, err := h.inverter.QueryWarning()
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	response.Json(w, http.StatusOK, data)
}

func (h InverterHandler) QueryStatus(w http.ResponseWriter, r *http.Request) {
	gs := h.querySnapshot.GetGeneralStatus()
	if gs == nil || gs.Timestamp.IsZero() {
		data, err := h.inverter.QueryPIGS()
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		gs = data
	}

	response.Json(w, http.StatusOK, gs)
}

func (h InverterHandler) AggregateInfo(w http.ResponseWriter, _ *http.Request) {
	if cached, found := h.cache.Get(cacheAggregateInfo); found {
		response.Json(w, http.StatusOK, cached)
		return
	}

	fw, _ := h.inverter.QueryFirmware()
	cpu2Fw, _ := h.inverter.QuerySecondCpuFirmware()
	remotePanelFw, _ := h.inverter.QueryRemotePanelFirmware()
	serial, _ := h.inverter.QuerySerial()
	serialNumber, _ := h.inverter.QuerySerialNumber()
	modelName, _ := h.inverter.QueryModelName()
	generalModelName, _ := h.inverter.QueryGeneralModelName()

	data := map[string]any{
		"firmware":            fw,
		"secondCpuFirmware":   cpu2Fw,
		"remotePanelFirmware": remotePanelFw,
		"serial":              serial,
		"serialNumber":        serialNumber,
		"modelName":           modelName,
		"generalModelName":    generalModelName,
	}
	h.cache.Set(cacheAggregateInfo, data, 1*time.Minute)

	response.Json(w, http.StatusOK, data)
}

func (h InverterHandler) QueryCommand(w http.ResponseWriter, r *http.Request) {
	cmd := r.URL.Query().Get("cmd")
	data, err := h.inverter.SendCommand(cmd)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	j := map[string]any{
		"data": string(data),
	}

	response.Json(w, http.StatusOK, j)
}

func (h InverterHandler) QueryFlags(w http.ResponseWriter, r *http.Request) {
	flags, err := h.inverter.QueryFlags()
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	response.Json(w, http.StatusOK, flags)
}

func (h InverterHandler) UpdateFlags(w http.ResponseWriter, r *http.Request) {
	var flags voltronic.DeviceFlags
	if err := json.NewDecoder(r.Body).Decode(&flags); err != nil {
		response.Error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.inverter.UpdateFlags(flags); err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	h.event.Publish(event.UpdateFlags, flags)

	response.Header(w, http.StatusAccepted)
}

func (h InverterHandler) UpdateSourcePriority(w http.ResponseWriter, r *http.Request) {
	var mode model.UpdateSourcePriority
	if err := json.NewDecoder(r.Body).Decode(&mode); err != nil {
		response.Error(w, r, http.StatusBadRequest, err)
		return
	}

	// FIXME: validate mode string

	if err := h.inverter.UpdateSourcePriority(mode.Source); err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	response.Header(w, http.StatusAccepted)
}

func (h InverterHandler) UpdateChargerSourcePriority(w http.ResponseWriter, r *http.Request) {
	var mode model.UpdateSourcePriority
	if err := json.NewDecoder(r.Body).Decode(&mode); err != nil {
		response.Error(w, r, http.StatusBadRequest, err)
		return
	}

	// FIXME: validate mode string

	if err := h.inverter.UpdateChargerPriority(mode.Source); err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	response.Header(w, http.StatusAccepted)
}

func (h InverterHandler) QueryMaxAcChargingCurrentValues(w http.ResponseWriter, r *http.Request) {
	current, err := h.inverter.QueryMaxAcChargingCurrentValues()
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	response.Json(w, http.StatusOK, map[string]any{
		"values": current,
	})
}

func (h InverterHandler) UpdateMaxAcChargingCurrent(w http.ResponseWriter, r *http.Request) {
	var mode model.UpdateChargingCurrent
	if err := json.NewDecoder(r.Body).Decode(&mode); err != nil {
		response.Error(w, r, http.StatusBadRequest, err)
		return
	}

	// FIXME: validate mode string

	if err := h.inverter.SetMaxAcChargingCurrent(mode.Current); err != nil {
		slog.Error(fmt.Sprintf("Error setting max charging current: %s", err))
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	response.Header(w, http.StatusAccepted)
}
