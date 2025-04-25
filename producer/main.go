package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

type NotificacionBiblioteca struct {
	Tipo          string `json:"tipo"`
	Detalle       string `json:"detalle"`
	UsuarioID     string `json:"usuario_id"`
	LibroID       string `json:"libro_id,omitempty"`
	FechaEvento   string `json:"fecha_evento"`
}

func main() {
	// Tipos de notificación para biblioteca
	tiposNotificacion := []string{
		"prestamo",
		"devolucion",
		"reserva",
		"recordatorio",
		"multa",
	}

	libros := []string{
		"El Quijote",
		"Cien años de soledad",
		"1984",
		"Don Juan Tenorio",
		"La sombra del viento",
	}

	usuarios := []string{
		"user_001",
		"user_002",
		"user_003",
		"user_004",
	}

	// Conexión RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Error al conectar: %v", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error al abrir canal: %v", err)
	}
	defer channel.Close()

	// Declarar exchange
	err = channel.ExchangeDeclare(
		"notificaciones_biblioteca", // name
		"direct",                   // type
		true,                       // durable
		false,                      // auto-deleted
		false,                      // internal
		false,                      // no-wait
		nil,                        // args
	)
	if err != nil {
		log.Fatalf("Error al declarar exchange: %v", err)
	}

	// Generar notificaciones aleatorias
	for i := 1; i <= 20; i++ {
		// Seleccionar tipo aleatorio
		tipo := tiposNotificacion[rand.Intn(len(tiposNotificacion))]
		usuario := usuarios[rand.Intn(len(usuarios))]
		libro := libros[rand.Intn(len(libros))]

		notificacion := NotificacionBiblioteca{
			Tipo:        tipo,
			Detalle:     generarDetalle(tipo, libro),
			UsuarioID:   usuario,
			LibroID:     libro,
			FechaEvento: time.Now().Format("2006-01-02 15:04:05"),
		}

		body, err := json.Marshal(notificacion)
		if err != nil {
			log.Printf("Error al codificar notificación: %v", err)
			continue
		}

		// Publicar con routing key específica
		routingKey := "notificaciones_usuario_" + usuario
		err = channel.Publish(
			"notificaciones_biblioteca", // exchange
			routingKey,                 // routing key
			false,                      // mandatory
			false,                      // immediate
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         body,
				DeliveryMode: amqp.Persistent,
				Timestamp:    time.Now(),
			},
		)
		if err != nil {
			log.Printf("Error al publicar mensaje: %v", err)
			continue
		}

		log.Printf("📤 Notificación enviada a '%s': %s", routingKey, notificacion.Detalle)
		time.Sleep(time.Duration(1+rand.Intn(3)) * time.Second)
	}

	log.Println("✅ Todas las notificaciones de biblioteca enviadas")
}

func generarDetalle(tipo, libro string) string {
	switch tipo {
	case "prestamo":
		return "El libro '" + libro + "' ha sido prestado con éxito. Fecha de devolución: " + 
			time.Now().AddDate(0, 0, 15).Format("02/01/2006")
	case "devolucion":
		return "El libro '" + libro + "' ha sido devuelto correctamente"
	case "reserva":
		return "Has reservado el libro '" + libro + "'. Te avisaremos cuando esté disponible"
	case "recordatorio":
		return "Recordatorio: El libro '" + libro + "' debe ser devuelto antes de " + 
			time.Now().AddDate(0, 0, 2).Format("02/01/2006")
	case "multa":
		return "Has recibido una multa por retraso en la devolución del libro '" + libro + "'"
	default:
		return "Mensaje importante sobre el libro '" + libro + "'"
	}
}