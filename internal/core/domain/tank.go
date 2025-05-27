package domain

import (
	"time"
)

// Tank representa la entidad principal de nuestro dominio - un tanque que almacena líquidos
type Tank struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Capacity       float64   `json:"capacity"`      // Capacidad total en litros
	CurrentLevel   float64   `json:"current_level"` // Nivel actual en litros
	LiquidType     string    `json:"liquid_type"`   // Tipo de líquido almacenado
	Temperature    float64   `json:"temperature"`   // Temperatura en grados Celsius
	LastUpdated    time.Time `json:"last_updated"`
	Status         string    `json:"status"`          // normal, warning, critical
	AlertThreshold float64   `json:"alert_threshold"` // Umbral para alertas (porcentaje)
}

// GetLevelPercentage calcula el porcentaje de llenado del tanque
func (t *Tank) GetLevelPercentage() float64 {
	if t.Capacity <= 0 {
		return 0
	}
	return (t.CurrentLevel / t.Capacity) * 100
}

// IsLevelCritical determina si el nivel del tanque está en estado crítico
func (t *Tank) IsLevelCritical() bool {
	percentage := t.GetLevelPercentage()
	return percentage <= t.AlertThreshold
}

// UpdateStatus actualiza el estado del tanque basado en las condiciones actuales
func (t *Tank) UpdateStatus() {
	percentage := t.GetLevelPercentage()

	switch {
	case percentage <= t.AlertThreshold:
		t.Status = "critical"
	case percentage <= t.AlertThreshold*2:
		t.Status = "warning"
	default:
		t.Status = "normal"
	}
}

// Measurement representa una medición del nivel del tanque en un momento específico
type Measurement struct {
	ID          string    `json:"id"`
	TankID      string    `json:"tank_id"`
	Level       float64   `json:"level"`
	Timestamp   time.Time `json:"timestamp"`
	Temperature float64   `json:"temperature"`
}
