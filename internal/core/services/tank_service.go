package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"monitor-tanques/internal/core/domain"
	"monitor-tanques/internal/core/ports"
)

// Errores comunes que puede devolver el servicio
var (
	ErrTankNotFound = errors.New("tank not found")
	ErrInvalidTank  = errors.New("invalid tank data")
)

// TankServiceImpl implementa la interfaz TankService
type TankServiceImpl struct {
	tankRepo        ports.TankRepository
	measurementRepo ports.MeasurementRepository
	alertNotifier   ports.AlertNotifier
}

// NewTankService crea una nueva instancia del servicio de tanques
func NewTankService(
	tankRepo ports.TankRepository,
	measurementRepo ports.MeasurementRepository,
	alertNotifier ports.AlertNotifier,
) ports.TankService {
	return &TankServiceImpl{
		tankRepo:        tankRepo,
		measurementRepo: measurementRepo,
		alertNotifier:   alertNotifier,
	}
}

// GetTank obtiene un tanque por su ID
func (s *TankServiceImpl) GetTank(ctx context.Context, id string) (*domain.Tank, error) {
	if id == "" {
		return nil, ErrInvalidTank
	}

	tank, err := s.tankRepo.GetTank(ctx, id)
	if err != nil {
		return nil, err
	}

	if tank == nil {
		return nil, ErrTankNotFound
	}

	// Obtenemos la última medición para actualizar el estado actual
	lastMeasurement, err := s.measurementRepo.GetLastMeasurement(ctx, id)
	if err == nil && lastMeasurement != nil {
		tank.CurrentLevel = lastMeasurement.Level
		tank.Temperature = lastMeasurement.Temperature
		tank.LastUpdated = lastMeasurement.Timestamp
		tank.UpdateStatus()
	}

	return tank, nil
}

// GetAllTanks obtiene todos los tanques
func (s *TankServiceImpl) GetAllTanks(ctx context.Context) ([]*domain.Tank, error) {
	tanks, err := s.tankRepo.GetAllTanks(ctx)
	if err != nil {
		return nil, err
	}

	// Actualizamos los estados de todos los tanques con las últimas mediciones
	for _, tank := range tanks {
		lastMeasurement, err := s.measurementRepo.GetLastMeasurement(ctx, tank.ID)
		if err == nil && lastMeasurement != nil {
			tank.CurrentLevel = lastMeasurement.Level
			tank.Temperature = lastMeasurement.Temperature
			tank.LastUpdated = lastMeasurement.Timestamp
			tank.UpdateStatus()
		}
	}

	return tanks, nil
}

// CreateTank crea un nuevo tanque
func (s *TankServiceImpl) CreateTank(ctx context.Context, tank *domain.Tank) error {
	if tank == nil || tank.Name == "" || tank.Capacity <= 0 {
		return ErrInvalidTank
	}

	// Aseguramos que tenga los valores predeterminados adecuados
	tank.Status = "normal"
	if tank.AlertThreshold <= 0 {
		tank.AlertThreshold = 10.0 // Valor predeterminado: 10%
	}
	tank.LastUpdated = time.Now()

	return s.tankRepo.SaveTank(ctx, tank)
}

// UpdateTank actualiza un tanque existente
func (s *TankServiceImpl) UpdateTank(ctx context.Context, tank *domain.Tank) error {
	if tank == nil || tank.ID == "" {
		return ErrInvalidTank
	}

	// Verificamos que el tanque exista
	existingTank, err := s.tankRepo.GetTank(ctx, tank.ID)
	if err != nil {
		return err
	}

	if existingTank == nil {
		return ErrTankNotFound
	}

	// Actualizamos el estado basado en los valores actuales
	tank.UpdateStatus()
	tank.LastUpdated = time.Now()

	return s.tankRepo.UpdateTank(ctx, tank)
}

// DeleteTank elimina un tanque por su ID
func (s *TankServiceImpl) DeleteTank(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidTank
	}

	// Verificamos que el tanque exista
	existingTank, err := s.tankRepo.GetTank(ctx, id)
	if err != nil {
		return err
	}

	if existingTank == nil {
		return ErrTankNotFound
	}

	return s.tankRepo.DeleteTank(ctx, id)
}

// MonitorTank monitorea un tanque específico y genera alertas si es necesario
func (s *TankServiceImpl) MonitorTank(ctx context.Context, tankID string) error {
	tank, err := s.GetTank(ctx, tankID)
	if err != nil {
		return err
	}

	// Verificamos si el nivel es crítico y enviamos una alerta
	if tank.IsLevelCritical() {
		message := "¡Alerta! El tanque " + tank.Name + " está en nivel crítico (" +
			"nivel: " + fmt.Sprintf("%.2f%%", tank.GetLevelPercentage()) + "). " +
			"Se requiere atención inmediata."

		return s.alertNotifier.SendAlert(ctx, tankID, message)
	}

	return nil
}

// AddMeasurement añade una nueva medición para un tanque
func (s *TankServiceImpl) AddMeasurement(ctx context.Context, measurement *domain.Measurement) error {
	if measurement == nil || measurement.TankID == "" || measurement.Level < 0 {
		return errors.New("invalid measurement data")
	}

	// Verificamos que el tanque exista
	tank, err := s.tankRepo.GetTank(ctx, measurement.TankID)
	if err != nil {
		return err
	}

	if tank == nil {
		return ErrTankNotFound
	}

	// Asignamos la marca de tiempo si no está establecida
	if measurement.Timestamp.IsZero() {
		measurement.Timestamp = time.Now()
	}

	// Guardamos la medición
	if err := s.measurementRepo.SaveMeasurement(ctx, measurement); err != nil {
		return err
	}

	// Actualizamos el tanque con los nuevos valores
	tank.CurrentLevel = measurement.Level
	tank.Temperature = measurement.Temperature
	tank.LastUpdated = measurement.Timestamp
	tank.UpdateStatus()

	if err := s.tankRepo.UpdateTank(ctx, tank); err != nil {
		return err
	}

	// Verificamos si necesitamos enviar alertas
	return s.MonitorTank(ctx, measurement.TankID)
}

// GetTankStatus obtiene el estado actual de un tanque
func (s *TankServiceImpl) GetTankStatus(ctx context.Context, tankID string) (string, error) {
	tank, err := s.GetTank(ctx, tankID)
	if err != nil {
		return "", err
	}

	return tank.Status, nil
}
