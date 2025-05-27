package ports

import (
	"context"

	"monitor-tanques/internal/core/domain"
)

// TankRepository define el puerto para operaciones de persistencia de tanques
type TankRepository interface {
	GetTank(ctx context.Context, id string) (*domain.Tank, error)
	GetAllTanks(ctx context.Context) ([]*domain.Tank, error)
	SaveTank(ctx context.Context, tank *domain.Tank) error
	UpdateTank(ctx context.Context, tank *domain.Tank) error
	DeleteTank(ctx context.Context, id string) error
}

// MeasurementRepository define el puerto para operaciones de persistencia de mediciones
type MeasurementRepository interface {
	SaveMeasurement(ctx context.Context, measurement *domain.Measurement) error
	GetMeasurementsByTankID(ctx context.Context, tankID string, limit int) ([]*domain.Measurement, error)
	GetLastMeasurement(ctx context.Context, tankID string) (*domain.Measurement, error)
}

// TankService define el puerto para el servicio de tanques
type TankService interface {
	GetTank(ctx context.Context, id string) (*domain.Tank, error)
	GetAllTanks(ctx context.Context) ([]*domain.Tank, error)
	CreateTank(ctx context.Context, tank *domain.Tank) error
	UpdateTank(ctx context.Context, tank *domain.Tank) error
	DeleteTank(ctx context.Context, id string) error
	MonitorTank(ctx context.Context, tankID string) error
	AddMeasurement(ctx context.Context, measurement *domain.Measurement) error
	GetTankStatus(ctx context.Context, tankID string) (string, error)
}

// AlertNotifier define el puerto para enviar notificaciones/alertas
type AlertNotifier interface {
	SendAlert(ctx context.Context, tankID string, message string) error
}
