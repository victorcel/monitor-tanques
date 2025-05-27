package repositories

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"monitor-tanques/internal/core/domain"
)

// MemoryMeasurementRepository implementa un repositorio de mediciones en memoria
// útil para desarrollo y pruebas
type MemoryMeasurementRepository struct {
	measurements map[string][]*domain.Measurement // clave: tankID, valor: slice de mediciones
	mutex        sync.RWMutex
}

// NewMemoryMeasurementRepository crea una nueva instancia del repositorio en memoria
func NewMemoryMeasurementRepository() *MemoryMeasurementRepository {
	return &MemoryMeasurementRepository{
		measurements: make(map[string][]*domain.Measurement),
	}
}

// SaveMeasurement guarda una nueva medición
func (r *MemoryMeasurementRepository) SaveMeasurement(ctx context.Context, measurement *domain.Measurement) error {
	if measurement == nil {
		return errors.New("measurement cannot be nil")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Si la medición no tiene fecha, la establecemos
	if measurement.Timestamp.IsZero() {
		measurement.Timestamp = time.Now()
	}

	// Guardamos una copia para evitar problemas de concurrencia
	measurementCopy := *measurement

	// Añadimos la medición al slice correspondiente al tankID
	if _, exists := r.measurements[measurement.TankID]; !exists {
		r.measurements[measurement.TankID] = make([]*domain.Measurement, 0)
	}

	r.measurements[measurement.TankID] = append(r.measurements[measurement.TankID], &measurementCopy)

	// Ordenamos las mediciones por timestamp (más recientes primero)
	sort.Slice(r.measurements[measurement.TankID], func(i, j int) bool {
		return r.measurements[measurement.TankID][i].Timestamp.After(
			r.measurements[measurement.TankID][j].Timestamp,
		)
	})

	return nil
}

// GetMeasurementsByTankID obtiene las mediciones para un tanque específico
func (r *MemoryMeasurementRepository) GetMeasurementsByTankID(ctx context.Context, tankID string, limit int) ([]*domain.Measurement, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if tankID == "" {
		return nil, errors.New("tank ID cannot be empty")
	}

	measurements, exists := r.measurements[tankID]
	if !exists {
		return make([]*domain.Measurement, 0), nil
	}

	// Limitamos la cantidad de mediciones a devolver si es necesario
	result := measurements
	if limit > 0 && limit < len(measurements) {
		result = measurements[:limit]
	}

	// Creamos copias de las mediciones para evitar problemas de concurrencia
	copies := make([]*domain.Measurement, len(result))
	for i, m := range result {
		copy := *m
		copies[i] = &copy
	}

	return copies, nil
}

// GetLastMeasurement obtiene la última medición para un tanque específico
func (r *MemoryMeasurementRepository) GetLastMeasurement(ctx context.Context, tankID string) (*domain.Measurement, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if tankID == "" {
		return nil, errors.New("tank ID cannot be empty")
	}

	measurements, exists := r.measurements[tankID]
	if !exists || len(measurements) == 0 {
		return nil, nil
	}

	// Las mediciones están ordenadas con la más reciente primero
	lastMeasurement := *measurements[0]
	return &lastMeasurement, nil
}
