package main

import (
	"monitor-tanques/cmd/api"
	"monitor-tanques/pkg/logger"
)

func main() {
	// Inicializamos el logger
	log := logger.NewSimpleLogger()

	// Configuramos la API
	config := api.DefaultConfig()

	// Creamos la instancia de la API
	app := api.NewAPI(config, log)

	// Configuramos las rutas
	app.SetupRoutes()

	// Iniciamos el servidor
	if err := app.Start(); err != nil {
		log.Fatal("Error al iniciar la aplicación", "error", err)
	}
}
