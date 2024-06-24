package main

import (
	"context"
	"fmt"
	"time"

	flexpool "github.com/pablolagosm/go-flexpool" // Asegúrese de reemplazar esto con el nombre real de su paquete
)

func main() {
	// Crear un nuevo pool con 5 trabajadores y capacidad para 10 tareas
	pool := flexpool.New(5, 10)

	// Crear un contexto con timeout
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Función para crear tareas de ejemplo
	createTask := func(id int, priority flexpool.Priority) func() error {
		return func() error {
			fmt.Printf("Ejecutando tarea %d con prioridad %v\n", id, priority)
			time.Sleep(time.Second) // Simular trabajo
			fmt.Printf("Tarea %d completada\n", id)
			return nil
		}
	}

	// Enviar algunas tareas al pool
	for i := 0; i < 10; i++ {
		priority := flexpool.LowPriority
		if i%3 == 0 {
			priority = flexpool.HighPriority
		} else if i%2 == 0 {
			priority = flexpool.MediumPriority
		}

		err := pool.Submit(ctx, createTask(i, priority), priority)
		if err != nil {
			fmt.Printf("Error al enviar tarea %d: %v\n", i, err)
		}
	}

	// Esperar un poco para que se procesen algunas tareas
	time.Sleep(1 * time.Second)

	// Cambiar el tamaño del pool
	err := pool.Resize(ctx, 3)
	if err != nil {
		fmt.Printf("Error al redimensionar el pool: %v\n", err)
	} else {
		fmt.Println("Pool redimensionado a 3 trabajadores")
	}

	// Esperar otro poco
	time.Sleep(3 * time.Second)

	// Cerrar el pool
	err = pool.Shutdown(ctx)
	if err != nil {
		fmt.Printf("Error al cerrar el pool: %v\n", err)
	} else {
		fmt.Println("Pool cerrado exitosamente")
	}
}
