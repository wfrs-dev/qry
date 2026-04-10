# qry

Cliente mínimo de línea de comandos para consultar bases de datos relacionales.

## Características

- Guarda conexiones con nombre para reutilizarlas
- Soporta MySQL, PostgreSQL y SQL Server
- Ejecuta consultas desde la línea de comandos o desde la entrada estándar (stdin)
- Muestra los resultados en formato tabla o lista
- Aplica automáticamente un límite de 200 filas en consultas `SELECT` sin límite explícito

## Instalación

### Desde el código fuente

Requiere Go 1.21+ y [Task](https://taskfile.dev).

```bash
git clone https://gitlab.com/wfrsgo/qry.git
cd qry
task build-linux   # Linux
task build-windows # Windows
```

El binario se genera en `./dist/{os}-{arch}/qry`.

## Uso

```
qry <comando> [opciones]
```

### Comandos

#### `add` - Agregar una conexión

Registra una nueva conexión a una base de datos de forma interactiva.

```bash
qry add <nombre>
```

Ejemplo:

```bash
qry add mibd
# Seleccionar driver: mysql / postgres / sqlserver
# Ingresar DSN
```

#### `list` - Listar conexiones guardadas

Muestra todas las conexiones registradas. Las contraseñas del DSN se ocultan.

```bash
qry list
```

#### `run` - Ejecutar una consulta

Ejecuta una consulta SQL en la conexión especificada.

```bash
qry run <nombre> -q "SELECT * FROM usuarios"
```

| Opción | Corto | Descripción |
|--------|-------|-------------|
| `--query` | `-q` | Consulta SQL a ejecutar |
| `--list`  | `-l` | Muestra el resultado en formato lista (por defecto: tabla) |

Si no se especifica `-q`, la consulta se lee desde `stdin`:

```bash
echo "SELECT * FROM usuarios" | qry run mibd
cat consulta.sql | qry run mibd
```

## Drivers soportados

| Driver      | Formato de DSN |
|-------------|----------------|
| `mysql`     | `usuario:contraseña@host:puerto/basedatos` |
| `postgres`  | `postgres://usuario:contraseña@host:puerto/basedatos` |
| `sqlserver` | `sqlserver://usuario:contraseña@host:puerto?database=basedatos` |

## Configuración

Las conexiones se guardan en `~/.config/qry/connections.json`.

```json
{
  "connections": {
    "mibd": {
      "driver": "postgres",
      "dsn": "postgres://usuario:contraseña@localhost:5432/mibasedatos"
    }
  }
}
```

## Desarrollo

```bash
task audit        # Verificaciones de calidad (vet, tests, vulnerabilidades)
task build-linux  # Compilar para Linux
task precommit    # Formatear y analizar archivos staged antes de un commit
task clean        # Eliminar archivos generados
```

## Licencia

Este proyecto está bajo la licencia MIT. Para más información, consulte el archivo [LICENSE](./LICENSE).
