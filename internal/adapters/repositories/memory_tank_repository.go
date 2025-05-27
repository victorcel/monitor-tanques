package repositories

import (
	"context"
	"errors"
	"sync"
	"time"

	"monitor-tanques/internal/core/domain"
)

// Errores comunes para el repositorio
var (
	ErrTankNotFound = errors.New("tank not found")
)

// MemoryTankRepository implementa un repositorio de tanques en memoria
// útil para desarrollo y pruebas
type MemoryTankRepository struct {
	tanks map[string]*domain.Tank
	mutex sync.RWMutex
}

// NewMemoryTankRepository crea una nueva instancia del repositorio en memoria
func NewMemoryTankRepository() *MemoryTankRepository {
	return &MemoryTankRepository{
		tanks: make(map[string]*domain.Tank),
	}
}

// GetTank obtiene un tanque por su ID
func (r *MemoryTankRepository) GetTank(ctx context.Context, id string) (*domain.Tank, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tank, exists := r.tanks[id]
	if !exists {
		return nil, ErrTankNotFound
	}

	// Devolvemos una copia para evitar problemas de concurrencia
	tankCopy := *tank
	return &tankCopy, nil
}

// GetAllTanks obtiene todos los tanques
func (r *MemoryTankRepository) GetAllTanks(ctx context.Context) ([]*domain.Tank, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tanks := make([]*domain.Tank, 0, len(r.tanks))
	for _, tank := range r.tanks {
		tankCopy := *tank
		tanks = append(tanks, &tankCopy)
	}

	return tanks, nil
}

// SaveTank guarda un nuevo tanque
func (r *MemoryTankRepository) SaveTank(ctx context.Context, tank *domain.Tank) error {
	if tank == nil {
		return errors.New("tank cannot be nil")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Si el tanque no tiene fecha de actualización, la establecemos
	if tank.LastUpdated.IsZero() {
		tank.LastUpdated = time.Now()
	}

	// Guardamos una copia para evitar problemas de concurrencia
	tankCopy := *tank
	r.tanks[tank.ID] = &tankCopy

	return nil
}

// UpdateTank actualiza un tanque existente
func (r *MemoryTankRepository) UpdateTank(ctx context.Context, tank *domain.Tank) error {
	if tank == nil {
		return errors.New("tank cannot be nil")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tanks[tank.ID]; !exists {
		return ErrTankNotFound
	}

	// Si el tanque no tiene fecha de actualización, la establecemos
	if tank.LastUpdated.IsZero() {
		tank.LastUpdated = time.Now()
	}

	// Guardamos una copia para evitar problemas de concurrencia
	tankCopy := *tank
	r.tanks[tank.ID] = &tankCopy

	return nil
}

// DeleteTank elimina un tanque por su ID
func (r *MemoryTankRepository) DeleteTank(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tanks[id]; !exists {
		return ErrTankNotFound
	}

	delete(r.tanks, id)
	return nil
}
