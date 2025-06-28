# Trackin

## Generador de código

Para correr el generador de codigo:
```
go run ./cmd/generate/
```

## Docker

### Buildear imagen Docker
```bash
docker build -t fonzeca/trackin:legacy .
```

### Pushear imagen al registro Docker
```bash
docker push fonzeca/trackin:legacy
```

### Comandos completos
```bash
# Buildear y pushear en una sola secuencia
docker build -t fonzeca/trackin:legacy .
docker push fonzeca/trackin:legacy
```

