package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/streadway/amqp"
)

// Estructura de notificación para el sistema de biblioteca
type NotificacionBiblioteca struct {
	Tipo          string `json:"tipo"`
	Detalle       string `json:"detalle"`
	UsuarioID     string `json:"usuario_id"`
	LibroID       string `json:"libro_id,omitempty"`
	FechaEvento   string `json:"fecha_evento"`
}

func main() {
	// Configuración inicial
	queueName := "notificaciones_usuario_" + os.Getenv("USUARIO_ID") // Cola personalizada
	if queueName == "notificaciones_usuario_" {
		queueName = "notificaciones_generales" // Cola por defecto
	}

	// Manejo de señales para cierre elegante
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Conexión con reconexión automática
	var conn *amqp.Connection
	var channel *amqp.Channel
	var err error

connect:
	for {
		conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			log.Printf("Error al conectar: %v. Reintentando en 5 segundos...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		channel, err = conn.Channel()
		if err != nil {
			log.Printf("Error al abrir canal: %v. Reintentando...", err)
			conn.Close()
			time.Sleep(5 * time.Second)
			continue connect
		}

		// Declarar cola durable
		_, err = channel.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // args
		)
		if err != nil {
			log.Fatalf("Error al hacer binding: %v", err)
		}

		err = channel.QueueBind(
		    queueName,                   // nombre de la cola
		    queueName,                   // routing key (usamos el mismo nombre)
		    "notificaciones_biblioteca", // exchange
		    false,
		    nil,
		)

		if err != nil {
			log.Printf("Error al declarar cola: %v. Reintentando...", err)
			channel.Close()
			conn.Close()
			time.Sleep(5 * time.Second)
			continue connect
		}

		break
	}

	defer conn.Close()
	defer channel.Close()

	// Consumir mensajes con QoS
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatalf("Error al configurar QoS: %v", err)
	}

	msgs, err := channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack (manualmente)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Error al consumir cola: %v", err)
	}

	log.Printf("📚 Esperando notificaciones en cola '%s'...", queueName)

	// Procesar mensajes en goroutine
	go func() {
		for msg := range msgs {
			var notificacion NotificacionBiblioteca
			if err := json.Unmarshal(msg.Body, &notificacion); err != nil {
				log.Printf("⚠️ Error al decodificar mensaje: %v", err)
				msg.Nack(false, true) // Rechazar pero reintentar
				continue
			}

			log.Printf("📩 Nueva notificación [%s] para usuario %s: %s", 
				notificacion.Tipo, notificacion.UsuarioID, notificacion.Detalle)
			msg.Ack(false) // Confirmar procesamiento
		}
	}()

	// Esperar señal de terminación
	<-signals
	log.Println("👋 Cerrando consumidor de notificaciones...")
}