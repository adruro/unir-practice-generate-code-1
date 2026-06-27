# TaskFlow — Gestor de Tareas Personal

## Descripción del Proyecto

TaskFlow es una aplicación web de gestión de tareas personales desarrollada en Go. Se distribuye como un binario único ejecutable que incluye todos los recursos necesarios (plantillas HTML, estilos CSS y JavaScript). Al ejecutarse, levanta un servidor web local y crea automáticamente una base de datos SQLite para persistencia.

## Objetivo

Permitir al usuario organizar sus tareas diarias mediante una interfaz web sencilla, con capacidad de categorización, priorización y filtrado.

## Requisitos Funcionales

### RF-01: Gestión de tareas (CRUD)
- Crear nuevas tareas con título, descripción, categoría y prioridad.
- Editar tareas existentes.
- Eliminar tareas.
- Listar todas las tareas.

### RF-02: Estados de tarea
- Marcar tareas como completadas.
- Revertir tareas completadas a pendientes.
- Visualizar claramente el estado de cada tarea.

### RF-03: Categorías
- Asignar una categoría a cada tarea: Trabajo, Personal, Estudio.
- Filtrar tareas por categoría.

### RF-04: Prioridades
- Asignar prioridad: Alta, Media, Baja.
- Filtrar tareas por prioridad.
- Indicación visual del nivel de prioridad (colores).

### RF-05: Filtrado y búsqueda
- Filtrar por estado (pendiente/completada).
- Filtrar por categoría.
- Filtrar por prioridad.
- Combinación de filtros simultáneos.

### RF-06: Persistencia
- Almacenamiento en base de datos SQLite local.
- La base de datos se crea automáticamente en el primer inicio.
- Los datos persisten entre reinicios de la aplicación.

## Requisitos No Funcionales

### RNF-01: Distribución
- Binario único ejecutable (sin dependencias externas en runtime).
- Templates HTML, CSS y JS embebidos en el binario.

### RNF-02: Interfaz
- Diseño web responsive (usable en escritorio y móvil).
- Interfaz limpia e intuitiva.

### RNF-03: Rendimiento
- Tiempo de respuesta inferior a 200ms para cualquier operación.
- Arranque del servidor en menos de 1 segundo.

### RNF-04: Portabilidad
- Compatible con Windows, Linux y macOS.
- Compilación cruzada mediante Makefile.

## Stack Tecnológico

| Componente | Tecnología |
|---|---|
| Lenguaje | Go 1.21+ |
| Servidor web | net/http (librería estándar) |
| Router | Chi (ligero, idiomático) |
| Base de datos | SQLite 3 (mattn/go-sqlite3) |
| Templates | html/template + embed |
| Frontend | HTML5, CSS3, JavaScript vanilla |
| Build | Makefile |

## Estructura del Proyecto

```
taskflow/
├── cmd/
│   └── taskflow/
│       └── main.go            # Punto de entrada
├── internal/
│   ├── handler/               # Handlers HTTP
│   ├── model/                 # Modelos de datos
│   ├── repository/            # Capa de acceso a datos
│   └── service/               # Lógica de negocio
├── web/
│   ├── templates/             # Plantillas HTML
│   └── static/                # CSS, JS
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

## Instrucciones de Ejecución

### Opción 1: Ejecutar el binario
```bash
./taskflow        # Linux/macOS
taskflow.exe      # Windows
```
La aplicación estará disponible en `http://localhost:8080`

### Opción 2: Compilar desde código fuente
```bash
make build        # Compila el binario
make run          # Compila y ejecuta
```
