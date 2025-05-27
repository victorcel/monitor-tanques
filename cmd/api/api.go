package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"monitor-tanques/internal/adapters/handlers"
	"monitor-tanques/internal/adapters/repositories"
	"monitor-tanques/internal/core/services"
	"monitor-tanques/pkg/logger"
)

// Config contiene la configuración de la API
type Config struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// DefaultConfig retorna una configuración predeterminada para la API
func DefaultConfig() Config {
	return Config{
		Port:            "8080",
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		ShutdownTimeout: 5 * time.Second,
	}
}

// API es el componente principal de la aplicación que maneja el servidor HTTP
type API struct {
	server *http.Server
	router *mux.Router
	logger logger.Logger
	config Config
}

// NewAPI crea una nueva instancia de la API
func NewAPI(config Config, logger logger.Logger) *API {
	router := mux.NewRouter()

	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}

	return &API{
		server: server,
		router: router,
		logger: logger,
		config: config,
	}
}

// SetupRoutes configura todas las rutas de la API
func (a *API) SetupRoutes() {
	// Creamos los repositorios (adaptadores de salida)
	tankRepo := repositories.NewMemoryTankRepository()
	measurementRepo := repositories.NewMemoryMeasurementRepository()

	// Creamos un notificador de alertas mock (podría ser reemplazado por uno real)
	alertNotifier := &mockAlertNotifier{logger: a.logger}

	// Creamos el servicio principal (puerto)
	tankService := services.NewTankService(tankRepo, measurementRepo, alertNotifier)

	// Creamos los handlers (adaptadores de entrada)
	tankHandler := handlers.NewTankHandler(tankService, a.logger)

	// Registramos las rutas
	tankHandler.RegisterRoutes(a.router)

	// Añadimos middleware para logging
	a.router.Use(a.loggingMiddleware)

	// Ruta de comprobación de estado
	a.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)
}

// loggingMiddleware registra información sobre cada solicitud HTTP
func (a *API) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		a.logger.Info("Request started",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		next.ServeHTTP(w, r)

		a.logger.Info("Request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}

// Start inicia el servidor HTTP
func (a *API) Start() error {
	// Configurar manejo de interrupciones para cierre graceful
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		a.logger.Info("Recibida señal para apagar el servidor, cerrando conexiones...")

		// Creamos un contexto con timeout para el cierre
		ctx, cancel := context.WithTimeout(context.Background(), a.config.ShutdownTimeout)
		defer cancel()

		if err := a.server.Shutdown(ctx); err != nil {
			a.logger.Error("Error al cerrar el servidor:", "error", err)
		}

		close(idleConnsClosed)
	}()

	a.logger.Info("Servidor iniciado", "port", a.config.Port)

	if err := a.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed
	a.logger.Info("Servidor apagado correctamente")
	return nil
}

// mockAlertNotifier es una implementación simple del puerto AlertNotifier para desarrollo
type mockAlertNotifier struct {
	logger logger.Logger
}

// SendAlert envía una alerta (en este caso, solo la registra)
func (n *mockAlertNotifier) SendAlert(ctx context.Context, tankID string, message string) error {
	n.logger.Warn("ALERTA", "tank_id", tankID, "message", message)
	return nil
}
