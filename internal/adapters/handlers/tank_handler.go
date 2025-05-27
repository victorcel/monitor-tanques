package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"monitor-tanques/internal/core/domain"
	"monitor-tanques/internal/core/ports"
	"monitor-tanques/pkg/logger"
)

// TankHandler maneja las peticiones HTTP relacionadas con los tanques
type TankHandler struct {
	tankService ports.TankService
	logger      logger.Logger
}

// NewTankHandler crea una nueva instancia del manejador de tanques
func NewTankHandler(tankService ports.TankService, logger logger.Logger) *TankHandler {
	return &TankHandler{
		tankService: tankService,
		logger:      logger,
	}
}

// RegisterRoutes registra las rutas del manejador en el router
func (h *TankHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/tanks", h.GetAllTanks).Methods(http.MethodGet)
	router.HandleFunc("/api/tanks/{id}", h.GetTank).Methods(http.MethodGet)
	router.HandleFunc("/api/tanks", h.CreateTank).Methods(http.MethodPost)
	router.HandleFunc("/api/tanks/{id}", h.UpdateTank).Methods(http.MethodPut)
	router.HandleFunc("/api/tanks/{id}", h.DeleteTank).Methods(http.MethodDelete)
	router.HandleFunc("/api/tanks/{id}/measurements", h.AddMeasurement).Methods(http.MethodPost)
}

// GetAllTanks devuelve todos los tanques
func (h *TankHandler) GetAllTanks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tanks, err := h.tankService.GetAllTanks(ctx)
	if err != nil {
		h.logger.Error("Failed to get tanks", "error", err)
		http.Error(w, "Error al obtener los tanques", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tanks); err != nil {
		h.logger.Error("Failed to encode tanks", "error", err)
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
		return
	}
}

// GetTank devuelve un tanque específico
func (h *TankHandler) GetTank(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	tank, err := h.tankService.GetTank(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get tank", "error", err, "id", id)
		http.Error(w, "Error al obtener el tanque", http.StatusInternalServerError)
		return
	}

	if tank == nil {
		http.Error(w, "Tanque no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tank); err != nil {
		h.logger.Error("Failed to encode tank", "error", err)
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
		return
	}
}

// CreateTank crea un nuevo tanque
func (h *TankHandler) CreateTank(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var tank domain.Tank

	if err := json.NewDecoder(r.Body).Decode(&tank); err != nil {
		h.logger.Error("Failed to decode request body", "error", err)
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	// Generamos un ID único si no se proporcionó
	if tank.ID == "" {
		tank.ID = uuid.New().String()
	}

	if err := h.tankService.CreateTank(ctx, &tank); err != nil {
		h.logger.Error("Failed to create tank", "error", err)
		http.Error(w, "Error al crear el tanque", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(tank); err != nil {
		h.logger.Error("Failed to encode response", "error", err)
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
		return
	}
}

// UpdateTank actualiza un tanque existente
func (h *TankHandler) UpdateTank(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	var tank domain.Tank
	if err := json.NewDecoder(r.Body).Decode(&tank); err != nil {
		h.logger.Error("Failed to decode request body", "error", err)
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	// Aseguramos que el ID en el cuerpo coincida con el de la URL
	tank.ID = id

	if err := h.tankService.UpdateTank(ctx, &tank); err != nil {
		h.logger.Error("Failed to update tank", "error", err, "id", id)
		http.Error(w, "Error al actualizar el tanque", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tank); err != nil {
		h.logger.Error("Failed to encode response", "error", err)
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
		return
	}
}

// DeleteTank elimina un tanque
func (h *TankHandler) DeleteTank(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.tankService.DeleteTank(ctx, id); err != nil {
		h.logger.Error("Failed to delete tank", "error", err, "id", id)
		http.Error(w, "Error al eliminar el tanque", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddMeasurement añade una medición a un tanque
func (h *TankHandler) AddMeasurement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	tankID := vars["id"]

	var measurement domain.Measurement
	if err := json.NewDecoder(r.Body).Decode(&measurement); err != nil {
		h.logger.Error("Failed to decode request body", "error", err)
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	// Asignamos el ID del tanque de la URL
	measurement.TankID = tankID

	// Generamos un ID único si no se proporcionó
	if measurement.ID == "" {
		measurement.ID = uuid.New().String()
	}

	// Establecemos la marca de tiempo si no se proporcionó
	if measurement.Timestamp.IsZero() {
		measurement.Timestamp = time.Now()
	}

	if err := h.tankService.AddMeasurement(ctx, &measurement); err != nil {
		h.logger.Error("Failed to add measurement", "error", err, "tankID", tankID)
		http.Error(w, "Error al añadir la medición", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(measurement); err != nil {
		h.logger.Error("Failed to encode response", "error", err)
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
		return
	}
}
