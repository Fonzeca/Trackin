package manager

import (
	"testing"
)

func TestManagerContainer_Singleton(t *testing.T) {
	// Reset para asegurar estado limpio
	defer ResetSingleton()

	// Primera instancia
	container1 := GetManagerContainer()

	// Segunda instancia (debería ser la misma)
	container2 := GetManagerContainer()

	// Verificar que es el mismo objeto (singleton)
	if container1 != container2 {
		t.Errorf("Container no es singleton - instancias diferentes")
	}
}

func TestManagerContainer_DependencyInjection(t *testing.T) {
	// Reset para asegurar estado limpio
	defer ResetSingleton()

	container := GetManagerContainer()

	// Verificar que todos los managers están inicializados
	routesManager := container.GetRoutesManager()
	if routesManager == nil {
		t.Errorf("RoutesManager no está inicializado")
	}

	zonasManager := container.GetZonasManager()
	if zonasManager == nil {
		t.Errorf("ZonasManager no está inicializado")
	}

	dataEntryManager := container.GetDataEntryManager()
	if dataEntryManager == nil {
		t.Errorf("DataEntryManager no está inicializado")
	}
}

func TestManagerContainer_Reset(t *testing.T) {
	// Crear una instancia
	container1 := GetManagerContainer()

	// Reset del singleton
	ResetSingleton()

	// Crear nueva instancia después del reset
	container2 := GetManagerContainer()

	// Deben ser instancias diferentes después del reset
	if container1 == container2 {
		t.Errorf("Reset no funcionó - misma instancia después de reset")
	}
}

func TestManagerContainer_ThreadSafety(t *testing.T) {
	// Reset para asegurar estado limpio
	defer ResetSingleton()

	// Crear múltiples goroutines para probar thread-safety
	containers := make(chan *ManagerContainer, 10)

	for i := 0; i < 10; i++ {
		go func() {
			container := GetManagerContainer()
			containers <- container
		}()
	}

	// Recoger todas las instancias
	firstContainer := <-containers
	for i := 1; i < 10; i++ {
		container := <-containers
		if container != firstContainer {
			t.Errorf("Thread-safety falla - instancias diferentes en goroutines")
		}
	}
}

// Benchmark para verificar performance del singleton
func BenchmarkManagerContainer_GetInstance(b *testing.B) {
	defer ResetSingleton()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetManagerContainer()
	}
}
