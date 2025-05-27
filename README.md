# Sistema de Monitoreo de Tanques

Este proyecto implementa un sistema backend para el monitoreo de tanques que guardan líquidos, siguiendo las mejores prácticas de desarrollo como Clean Code, Arquitectura Hexagonal, principios SOLID, y con soporte para pruebas automatizadas, Docker y medición de cobertura.

## Características

- **Monitoreo en tiempo real** de niveles y temperatura de tanques de líquidos
- **Sistema de alertas** cuando los niveles son críticos
- **API REST** para integración con otros sistemas
- **Arquitectura Hexagonal** (Ports & Adapters) para facilitar la extensibilidad y mantenimiento
- **Containerización** con Docker para despliegue simplificado
- **Pruebas automatizadas** para garantizar la calidad del código

## Arquitectura

El proyecto sigue la Arquitectura Hexagonal (también conocida como Ports & Adapters), que separa claramente la lógica de negocio de los detalles de implementación externa, facilitando el mantenimiento y las pruebas.

### Estructura del proyecto

```
monitor-tanques/
├── cmd/                    # Punto de entrada de aplicaciones
│   └── api/                # Aplicación API REST
├── configs/                # Archivos de configuración
├── deployments/            # Archivos de despliegue
├── docs/                   # Documentación
├── internal/               # Código interno no exportable
│   ├── adapters/           # Adaptadores (implementaciones de puertos)
│   │   ├── handlers/       # Handlers HTTP
│   │   └── repositories/   # Implementaciones de repositorios
│   └── core/               # Núcleo de la aplicación
│       ├── domain/         # Modelos y entidades de dominio
│       ├── ports/          # Interfaces (puertos)
│       └── services/       # Servicios de dominio (lógica de negocio)
├── pkg/                    # Bibliotecas exportables
│   ├── config/             # Utilidades de configuración
│   └── logger/             # Sistema de logging
├── scripts/                # Scripts útiles
├── test/                   # Tests
│   ├── integration/        # Tests de integración
│   └── unit/               # Tests unitarios
├── Dockerfile              # Definición de imagen Docker
├── go.mod                  # Dependencias de Go
└── main.go                 # Punto de entrada principal
```

### Componentes principales

- **Dominio**: Contiene las entidades principales (`Tank` y `Measurement`) con su lógica inherente.
- **Puertos**: Define las interfaces que permiten la comunicación entre el núcleo y el mundo exterior.
- **Servicios**: Implementa la lógica de negocio principal, utilizando los puertos para las operaciones.
- **Adaptadores**: Conecta el núcleo con tecnologías específicas (bases de datos, APIs, etc.).

## Principios de diseño aplicados

- **Clean Code**: Código legible, nombres descriptivos, funciones pequeñas con responsabilidad única.
- **SOLID**:
  - **S**: Cada clase tiene una única razón para cambiar.
  - **O**: El sistema es extensible sin modificar código existente.
  - **L**: Las implementaciones pueden sustituir a sus interfaces.
  - **I**: Interfaces específicas para cada cliente.
  - **D**: Dependencia de abstracciones, no de implementaciones concretas.
- **Inmutabilidad**: Se favorecen estructuras de datos inmutables para mejorar la concurrencia.
- **Manejo explícito de errores**: Errores claros y descriptivos.

## Requisitos

- Go 1.24 o superior
- Docker (opcional para containerización)

## Configuración e instalación

### Ejecución local

1. Clona el repositorio:
   ```bash
   git clone <url-del-repositorio>
   cd monitor-tanques
   ```

2. Instala las dependencias:
   ```bash
   go mod download
   ```

3. Ejecuta la aplicación:
   ```bash
   go run main.go
   ```

La API estará disponible en http://localhost:8080.

### Ejecución con Docker

1. Construye la imagen:
   ```bash
   docker build -t monitor-tanques .
   ```

2. Ejecuta el contenedor:
   ```bash
   docker run -p 8080:8080 monitor-tanques
   ```

## Pruebas

### Ejecutar pruebas unitarias

```bash
go test -v ./test/unit/...
```

### Análisis de cobertura

```bash
go test -coverprofile=coverage.out ./test/unit/...
go tool cover -html=coverage.out
```

## API REST

La API expone los siguientes endpoints:

### Tanques

- **GET** `/api/tanks`: Obtener todos los tanques.
- **GET** `/api/tanks/{id}`: Obtener un tanque específico.
- **POST** `/api/tanks`: Crear un nuevo tanque.
  ```json
  {
    "name": "Tanque Principal",
    "capacity": 1000.0,
    "current_level": 500.0,
    "liquid_type": "Agua",
    "temperature": 25.0,
    "alert_threshold": 10.0
  }
  ```
- **PUT** `/api/tanks/{id}`: Actualizar un tanque existente.
- **DELETE** `/api/tanks/{id}`: Eliminar un tanque.

### Mediciones

- **POST** `/api/tanks/{id}/measurements`: Añadir una nueva medición a un tanque.
  ```json
  {
    "level": 450.0,
    "temperature": 26.5
  }
  ```

### Estado

- **GET** `/health`: Verificar el estado del servicio.

## Desarrollo

### Extensión del sistema

Para añadir una nueva funcionalidad, siga estos pasos:

1. Si es necesario, añada nuevas entidades o modifique las existentes en `internal/core/domain`.
2. Defina o actualice las interfaces (puertos) en `internal/core/ports`.
3. Implemente la lógica de negocio en los servicios en `internal/core/services`.
4. Añada los adaptadores necesarios en `internal/adapters`.
5. Añada pruebas unitarias para los nuevos componentes.

### Persistencia de datos

Actualmente, el sistema utiliza repositorios en memoria para desarrollo y pruebas. Para implementar persistencia real:

1. Cree nuevos adaptadores en `internal/adapters/repositories` que implementen las interfaces del puerto correspondiente.
2. Utilice su sistema de base de datos preferido (SQL, NoSQL, etc.).
3. Actualice la configuración en `cmd/api/api.go` para utilizar los nuevos repositorios.

## Licencia

Este proyecto está licenciado bajo BSD. Consulte el archivo LICENSE para más detalles.

## Contribución

Las contribuciones son bienvenidas. Por favor, siga estos pasos:

1. Fork el repositorio
2. Cree una rama para su funcionalidad (`git checkout -b feature/amazing-feature`)
3. Haga commit de sus cambios (`git commit -m 'Add some amazing feature'`)
4. Push a la rama (`git push origin feature/amazing-feature`)
5. Abra un Pull Request
