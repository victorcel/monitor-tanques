package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"monitor-tanques/internal/adapters/repositories"
	"monitor-tanques/internal/core/domain"
	"monitor-tanques/internal/core/services"
)

// MockAlertNotifier es un mock para el notificador de alertas para pruebas
type MockAlertNotifier struct {
	AlertsSent  int
	LastTankID  string
	LastMessage string
}

func (m *MockAlertNotifier) SendAlert(ctx context.Context, tankID string, message string) error {
	m.AlertsSent++
	m.LastTankID = tankID
	m.LastMessage = message
	return nil
}

func createTestTank() *domain.Tank {
	return &domain.Tank{
		ID:             uuid.New().String(),
		Name:           "Tanque de Prueba",
		Capacity:       1000.0,
		CurrentLevel:   500.0,
		LiquidType:     "Agua",
		Temperature:    25.0,
		LastUpdated:    time.Now(),
		Status:         "normal",
		AlertThreshold: 10.0,
	}
}

func createTestMeasurement(tankID string, level float64) *domain.Measurement {
	return &domain.Measurement{
		ID:          uuid.New().String(),
		TankID:      tankID,
		Level:       level,
		Temperature: 25.0,
		Timestamp:   time.Now(),
	}
}

func TestTankService_CreateTank(t *testing.T) {
	// Arrange
	tankRepo := repositories.NewMemoryTankRepository()
	measurementRepo := repositories.NewMemoryMeasurementRepository()
	alertNotifier := &MockAlertNotifier{}

	service := services.NewTankService(tankRepo, measurementRepo, alertNotifier)
	ctx := context.Background()

	tank := createTestTank()

	// Act
	err := service.CreateTank(ctx, tank)

	// Assert
	if err != nil {
		t.Fatalf("Error al crear el tanque: %v", err)
	}

	// Verificamos que el tanque se haya guardado correctamente
	savedTank, err := service.GetTank(ctx, tank.ID)
	if err != nil {
		t.Fatalf("Error al obtener el tanque: %v", err)
	}

	if savedTank == nil {
		t.Fatal("El tanque guardado no se encontró")
	}

	if savedTank.ID != tank.ID {
		t.Errorf("ID del tanque incorrecto. Esperado: %s, Obtenido: %s", tank.ID, savedTank.ID)
	}

	if savedTank.Name != tank.Name {
		t.Errorf("Nombre del tanque incorrecto. Esperado: %s, Obtenido: %s", tank.Name, savedTank.Name)
	}
}

func TestTankService_UpdateTank(t *testing.T) {
	// Arrange
	tankRepo := repositories.NewMemoryTankRepository()
	measurementRepo := repositories.NewMemoryMeasurementRepository()
	alertNotifier := &MockAlertNotifier{}

	service := services.NewTankService(tankRepo, measurementRepo, alertNotifier)
	ctx := context.Background()

	// Creamos un tanque para actualizar
	tank := createTestTank()
	err := service.CreateTank(ctx, tank)
	if err != nil {
		t.Fatalf("Error al crear el tanque para la prueba: %v", err)
	}

	// Actualizamos algunos valores
	tank.Name = "Tanque Actualizado"
	tank.LiquidType = "Aceite"

	// Act
	err = service.UpdateTank(ctx, tank)

	// Assert
	if err != nil {
		t.Fatalf("Error al actualizar el tanque: %v", err)
	}

	// Verificamos que los cambios se hayan guardado
	updatedTank, err := service.GetTank(ctx, tank.ID)
	if err != nil {
		t.Fatalf("Error al obtener el tanque actualizado: %v", err)
	}

	if updatedTank == nil {
		t.Fatal("El tanque actualizado no se encontró")
	}

	if updatedTank.Name != "Tanque Actualizado" {
		t.Errorf("El nombre del tanque no se actualizó correctamente. Esperado: %s, Obtenido: %s",
			"Tanque Actualizado", updatedTank.Name)
	}

	if updatedTank.LiquidType != "Aceite" {
		t.Errorf("El tipo de líquido no se actualizó correctamente. Esperado: %s, Obtenido: %s",
			"Aceite", updatedTank.LiquidType)
	}
}

func TestTankService_DeleteTank(t *testing.T) {
	// Arrange
	tankRepo := repositories.NewMemoryTankRepository()
	measurementRepo := repositories.NewMemoryMeasurementRepository()
	alertNotifier := &MockAlertNotifier{}

	service := services.NewTankService(tankRepo, measurementRepo, alertNotifier)
	ctx := context.Background()

	// Creamos un tanque para eliminar
	tank := createTestTank()
	err := service.CreateTank(ctx, tank)
	if err != nil {
		t.Fatalf("Error al crear el tanque para la prueba: %v", err)
	}

	// Act
	err = service.DeleteTank(ctx, tank.ID)

	// Assert
	if err != nil {
		t.Fatalf("Error al eliminar el tanque: %v", err)
	}

	// Verificamos que el tanque ya no exista
	deletedTank, err := service.GetTank(ctx, tank.ID)
	if err == nil || deletedTank != nil {
		t.Errorf("El tanque no se eliminó correctamente")
	}
}

func TestTankService_AddMeasurement(t *testing.T) {
	// Arrange
	tankRepo := repositories.NewMemoryTankRepository()
	measurementRepo := repositories.NewMemoryMeasurementRepository()
	alertNotifier := &MockAlertNotifier{}

	service := services.NewTankService(tankRepo, measurementRepo, alertNotifier)
	ctx := context.Background()

	// Creamos un tanque para añadir mediciones
	tank := createTestTank()
	err := service.CreateTank(ctx, tank)
	if err != nil {
		t.Fatalf("Error al crear el tanque para la prueba: %v", err)
	}

	// Creamos una medición
	measurement := createTestMeasurement(tank.ID, 300.0)

	// Act
	err = service.AddMeasurement(ctx, measurement)

	// Assert
	if err != nil {
		t.Fatalf("Error al añadir la medición: %v", err)
	}

	// Verificamos que el tanque se haya actualizado con el nuevo nivel
	updatedTank, err := service.GetTank(ctx, tank.ID)
	if err != nil {
		t.Fatalf("Error al obtener el tanque actualizado: %v", err)
	}

	if updatedTank.CurrentLevel != 300.0 {
		t.Errorf("El nivel del tanque no se actualizó correctamente. Esperado: %.2f, Obtenido: %.2f",
			300.0, updatedTank.CurrentLevel)
	}
}

func TestTankService_MonitorTank_Critical(t *testing.T) {
	// Arrange
	tankRepo := repositories.NewMemoryTankRepository()
	measurementRepo := repositories.NewMemoryMeasurementRepository()
	alertNotifier := &MockAlertNotifier{}

	service := services.NewTankService(tankRepo, measurementRepo, alertNotifier)
	ctx := context.Background()

	// Creamos un tanque con nivel crítico (por debajo del umbral de alerta)
	tank := createTestTank()
	tank.Capacity = 1000.0
	tank.CurrentLevel = 50.0 // 5%, por debajo del umbral de 10%

	err := service.CreateTank(ctx, tank)
	if err != nil {
		t.Fatalf("Error al crear el tanque para la prueba: %v", err)
	}

	// Act
	err = service.MonitorTank(ctx, tank.ID)

	// Assert
	if err != nil {
		t.Fatalf("Error al monitorear el tanque: %v", err)
	}

	// Verificamos que se haya enviado una alerta
	if alertNotifier.AlertsSent != 1 {
		t.Errorf("Se esperaba 1 alerta enviada, pero se enviaron %d", alertNotifier.AlertsSent)
	}

	if alertNotifier.LastTankID != tank.ID {
		t.Errorf("ID incorrecto en la alerta. Esperado: %s, Obtenido: %s", tank.ID, alertNotifier.LastTankID)
	}
}
