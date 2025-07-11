# Estructura del Proyecto Trackin

## Estructura de Directorios

```
Trackin/
├── cmd/                    # Entry points de la aplicación
│   ├── server/            # HTTP server entry point
│   └── cli/               # CLI commands (generate, etc.)
├── internal/              # Código interno de la aplicación
│   ├── core/             # Lógica de negocio central
│   │   ├── domain/       # Modelos de dominio y interfaces
│   │   ├── managers/     # Business logic managers
│   │   └── services/     # Servicios de aplicación
│   ├── infrastructure/   # Implementaciones de infraestructura
│   │   ├── database/     # Database access layer
│   │   │   ├── model/    # Modelos de base de datos generados
│   │   │   └── query/    # Queries generadas
│   │   ├── messaging/    # RabbitMQ, sender, etc.
│   │   └── geolocation/  # Operaciones geoespaciales
│   ├── interfaces/       # HTTP handlers y APIs
│   │   ├── http/         # REST API handlers
│   │   └── messaging/    # Message handlers y entry points
│   └── container/        # Dependency injection container
├── pkg/                  # Código público reutilizable
├── tests/                # Tests e2e y de integración
│   ├── suites/          # Test suites
│   └── mocks/           # Mocks para testing
├── docs/                 # Documentación
├── configs/              # Archivos de configuración
└── scripts/              # Scripts de deployment, etc.
```

## Descripción de Componentes

### Core (internal/core/)
- **domain/**: Contiene los modelos de dominio y las interfaces de negocio
- **managers/**: Lógica de negocio, incluyendo managers para zonas, data entry, etc.
- **services/**: Servicios de aplicación como monitoring y configuración global

### Infrastructure (internal/infrastructure/)
- **database/**: Capa de acceso a datos con modelos y queries generadas
- **messaging/**: Implementaciones de RabbitMQ y servicios de mensajería
- **geolocation/**: Operaciones geoespaciales y cálculos de ubicación

### Interfaces (internal/interfaces/)
- **http/**: Handlers HTTP y APIs REST
- **messaging/**: Handlers de mensajes y entry points para procesamiento de datos

### Tests (tests/)
- Incluye mocks, suites de pruebas e integración
- Separado del código de producción para mejor organización

## Principios de Arquitectura

1. **Separación de Responsabilidades**: Cada capa tiene una responsabilidad específica
2. **Inversión de Dependencias**: Las capas internas no dependen de las externas
3. **Testabilidad**: Estructura que facilita las pruebas unitarias e integración
4. **Mantenibilidad**: Código organizado y fácil de mantener