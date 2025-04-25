# Sistema asincronico medio utilizando GO y RabbitMQ

Este proyecto simula un sistema de notificaciones para una biblioteca universitaria, donde se envían mensajes sobre préstamos, devoluciones y reservas de libros a través de colas de mensajes, los cuales son consumidos por clientes específicos.

Utiliza Go y RabbitMQ para la comunicación asincrónica entre producer y consumer, **sin gRPC ni Docker**.

### Características:

- Producer: Genera notificaciones de biblioteca y las envía a colas específicas en RabbitMQ

- Consumer: Recibe notificaciones personalizadas desde RabbitMQ según el usuario

- Tipos de notificaciones:
  - Préstamos de libros
  - Devoluciones
  - Reservas
  - Recordatorios
  - Multas por retraso

### Estructura del proyecto:

```
SD20251-medium-homework-system.git/
├── producer/
│   └── main.go
└── consumer/
    └── main.go
```

# Instrucciones para ejecutar el proyecto

### Requisitos para todos los sistemas

1. instalar Go: https://golang.org/dl/

2. Instalar RabbitMQ:

**Para Linux:**

```
sudo apt update
sudo apt install rabbitmq-server
sudo systemctl start rabbitmq-server
```

**Para Windows:**

Descargar desde https://rabbitmq.com/install-windows.html

### Configuración y ejecución

1. Clonar el repositorio:

```
git clone https://github.com/FelipeMarchantVargas/SD20251-medium-homework-system.git
cd SD20251-medium-homework-system
```

2. Inicializar modulo Go

```
go mod init medium-homework-system
go mod tidy
```

3. Instalacion de dependencia para utilizar RabbitMQ

```
go get github.com/streadway/amqp
```

### Ejecución en Linux/mac OS

1. Ejecutar el consumer (en una terminal):

```
cd consumer
export USUARIO_ID=user_001  # Cambiar por el ID deseado
go run main.go
```

2. Ejecutar el producer (en otra terminal):

```
cd producer
go run main.go
```

### Ejecución en Windows (CMD)

1. Ejecutar el consumer (en una terminal):

```
cd consumer
set USUARIO_ID=user_001  # Cambiar por el ID deseado
go run main.go
```

2. Ejecutar el producer (en otra terminal):

```
cd producer
go run main.go
```

### Ejecución en Windows (PowerShell)

1. Ejecutar el consumer (en una terminal):

```
cd consumer
$env:USUARIO_ID="user_001"  # Cambiar por el ID deseado
go run main.go
```

2. Ejecutar el producer (en otra terminal):

```
cd producer
go run main.go
```

Esto enviará una tarea a RabbitMQ, que será recibida y procesada por el consumidor.

## Funcionamiento esperado

1. El productor enviará 20 notificaciones aleatorias a diferentes usuarios

2. Cada consumidor recibirá solo las notificaciones dirigidas a su USUARIO_ID

3. Las notificaciones se mostrarán en la consola del consumidor con este formato:

```
📩 Nueva notificación [tipo] para usuario [id]: [mensaje detallado]
```

### Personalización

Puedes modificar los usuarios y libros disponibles editando los slices en producer/main.go:

```
libros := []string{
    "El Quijote",
    "Cien años de soledad",
    // ... agregar más títulos
}

usuarios := []string{
    "user_001",
    "user_002",
    // ... agregar más usuarios
}
```
