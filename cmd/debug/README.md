# Debug Logs - Herramienta de Depuración de GPS

Esta herramienta permite analizar y limpiar logs de GPS que contienen datos erróneos utilizando la función `IsValidPoint` del paquete de geolocalización.

## 🗄️ Configuración de Base de Datos

### Variables de Entorno

El sistema se conecta a MySQL usando estas variables de entorno:

| Variable | Valor por Defecto | Descripción |
|----------|-------------------|-------------|
| `trackinDbHost` | `vps-2367826-x.dattaweb.com:3306` | Host y puerto de MySQL |
| `trackinDbUser` | `root` | Usuario de la base de datos |
| `trackinDbPass` | `carmind-db` | Contraseña del usuario |

### Opción 1: Archivo .env (Recomendado)

```bash
# 1. Copiar el archivo de ejemplo
cp .env.example .env

# 2. Editar la configuración
nano .env  # o tu editor preferido
```

Contenido del archivo `.env`:
```bash
trackinDbHost=localhost:3306
trackinDbUser=tu_usuario
trackinDbPass=tu_contraseña
```

```bash
# 3. Ejecutar con configuración personalizada
./run_debug.sh -start='2025-09-01 00:00:00' -end='2025-09-15 23:59:59' -imei='123456789' -verbose
```

### Opción 2: Variables de Entorno Directas

**Windows (PowerShell):**
```powershell
$env:trackinDbHost = "localhost:3306"
$env:trackinDbUser = "tu_usuario"  
$env:trackinDbPass = "tu_contraseña"
./debug_logs.exe -start='2025-09-01 00:00:00' -end='2025-09-15 23:59:59' -imei='123456789'
```

**Windows (CMD):**
```cmd
set trackinDbHost=localhost:3306
set trackinDbUser=tu_usuario
set trackinDbPass=tu_contraseña
debug_logs.exe -start="2025-09-01 00:00:00" -end="2025-09-15 23:59:59" -imei="123456789"
```

**Linux/Mac:**
```bash
export trackinDbHost="localhost:3306"
export trackinDbUser="tu_usuario"
export trackinDbPass="tu_contraseña"
./debug_logs.exe -start='2025-09-01 00:00:00' -end='2025-09-15 23:59:59' -imei='123456789'
```

## Compilación

```bash
cd cmd/debug
go build -o debug_logs debug_logs.go
```

## Uso

### Análisis sin eliminar datos (Dry-run)
```bash
./debug_logs -start='2025-12-15 08:00:00' -end='2025-12-31 18:30:00' -imei='356307042441013' -verbose
```

### Eliminar datos inválidos
```bash
./debug_logs -start='2025-12-15 08:00:00' -end='2025-12-31 18:30:00' -imei='356307042441013' -dry-run=false
```

## Parámetros

| Parámetro | Requerido | Descripción | Ejemplo |
|-----------|-----------|-------------|---------|
| `-start` | ✅ | Fecha y hora de inicio | `2025-12-15 08:00:00` |
| `-end` | ✅ | Fecha y hora de fin | `2025-12-31 18:30:00` |
| `-imei` | ✅ | IMEI del dispositivo GPS | `356307042441013` |
| `-dry-run` | ❌ | Solo análisis (por defecto: true) | `false` |
| `-verbose` | ❌ | Información detallada | Incluir para activar |

## Criterios de Validación

La herramienta utiliza la función `IsValidPoint` que valida:

1. **Orden temporal**: La fecha del log actual debe ser posterior al anterior
2. **Velocidad máxima**: No debe exceder los 350 km/h
3. **Distancia mínima**: Movimientos de al menos 5 metros
4. **Tiempo mínimo**: Al menos 1 minuto entre puntos cercanos
5. **División por cero**: Evita errores cuando la diferencia de tiempo es cero

## Salida del Programa

### Análisis Detallado (-verbose)
```
📊 Analizando 1250 logs para IMEI: 356307042441013
📅 Rango: 2025-12-15 08:00:00 a 2025-12-31 18:30:00
🔍 Modo: DRY-RUN (solo análisis)
================================================================================

Log #1 [2025-12-15 08:00:15]: ✅ Válido - Primer log
Log #2 [2025-12-15 08:01:45]: ✅ Válido (0.125 km, 5.0 km/h)
Log #3 [2025-12-15 08:02:00]: ❌ Inválido (velocidad excesiva: 450.2 km/h)
    📍 Lat: -34.123456, Lon: -58.987654
    📏 Distancia: 1.875 km, Velocidad: 450.20 km/h, Tiempo: 0.25 min
```

### Resumen
```
📋 RESUMEN DE VALIDACIÓN
==================================================
📈 Logs totales:     1250
✅ Logs válidos:     1180 (94.4%)
❌ Logs inválidos:   70 (5.6%)
📏 Distancia total:  245.67 km
🚀 Velocidad máxima: 450.20 km/h

🗑️  LOGS INVÁLIDOS ENCONTRADOS:
  • Log #3 [2025-12-15 08:02:00]: ❌ Inválido (velocidad excesiva: 450.2 km/h)
  • Log #15 [2025-12-15 08:15:30]: ❌ Inválido (fecha anterior es posterior)
  ...
```

### Eliminación (dry-run=false)
```
🗑️  Eliminando 70 logs inválidos...
✅ Se eliminaron 70 registros exitosamente
```

## Casos de Uso

### Problema: GPS con datos erróneos
Un dispositivo GPS está enviando coordenadas incorrectas debido a fallas de hardware o interferencias, generando puntos que implican velocidades imposibles o saltos temporales.

### Solución: Limpieza automática
1. **Análisis inicial**: Ejecutar con `-verbose` para revisar qué logs son problemáticos
2. **Backup**: Hacer respaldo de la base de datos antes de eliminar datos
3. **Limpieza**: Ejecutar con `-dry-run=false` para eliminar logs inválidos
4. **Verificación**: Volver a ejecutar el análisis para confirmar la limpieza

## Ejemplos Prácticos

### Analizar una semana específica
```bash
./debug_logs -start='2025-09-01 00:00:00' -end='2025-09-07 23:59:59' -imei='123456789012345' -verbose
```

### Limpiar datos del último mes
```bash
./debug_logs -start='2025-08-15 00:00:00' -end='2025-09-15 23:59:59' -imei='123456789012345' -dry-run=false
```

### Análisis rápido sin detalles
```bash
./debug_logs -start='2025-09-14 00:00:00' -end='2025-09-15 23:59:59' -imei='123456789012345'
```

## Consideraciones de Seguridad

- ⚠️ **Siempre hacer backup** antes de ejecutar con `-dry-run=false`
- 🔍 **Revisar primero** con `-verbose` para entender qué se va a eliminar
- 📊 **Validar resultados** después de la limpieza
- 🚨 **No ejecutar en producción** sin pruebas previas

## Resolución de Problemas

### Error: "No se encontraron logs"
- Verificar que el IMEI existe en la base de datos
- Confirmar que las fechas están en el formato correcto
- Revisar que hay datos en el rango especificado

### Error de conexión a base de datos
- Verificar configuración en `configs/config.json`
- Confirmar que la base de datos está accesible
- Revisar credenciales de conexión
