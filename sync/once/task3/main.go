package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Интерфейс для всех плагинов
type Plugin interface {
	Execute() string
}

// Управляет инициализацией и доступом к плагинам
type PluginManager struct {
	plugins map[string]*pluginEntry
	mu      sync.RWMutex
}

type pluginEntry struct {
	//Добавить необходимые поля для однократной инициализации
	sOnce      sync.Once
	initFn     func() (Plugin, error)
	initPlugin Plugin
	err        error
}

// NewPluginManager создает новый менеджер плагинов
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]*pluginEntry),
	}
}

// RegisterPlugin регистрирует новый плагин
func (pm *PluginManager) RegisterPlugin(name string, initFn func() (Plugin, error)) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.plugins[name] = &pluginEntry{
		initFn: initFn,
	}
}

// GetPlugin возвращает инициализированный плагин
func (pm *PluginManager) GetPlugin(name string) (Plugin, error) {
	pm.mu.RLock()
	entry, ok := pm.plugins[name]
	pm.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("плагин %s не существует", name)
	}

	entry.sOnce.Do(func() {
		entry.initPlugin, entry.err = entry.initFn()
	})

	return entry.initPlugin, entry.err
	// Реализовать:
	// 1. Проверку существования плагина
	// 2. Потокобезопасную однократную инициализацию
	// 3. Обработку и кэширование ошибок
	// 4. Возврат кэшированного результата
}

// DemoPlugin реализация плагина
type DemoPlugin struct{}

func (p *DemoPlugin) Execute() string {
	return "DemoPlugin executed successfully!"
}

func initDemo() (Plugin, error) {
	// Имитация длительной инициализации
	//t := time.Now()
	time.Sleep(500 * time.Millisecond)
	//log.Println(time.Since(t))
	return &DemoPlugin{}, nil
}

func main() {
	pm := NewPluginManager()

	pm.RegisterPlugin("demo", initDemo)
	pm.RegisterPlugin("broken", func() (Plugin, error) {
		return nil, fmt.Errorf("simulated error")
	})

	var wg sync.WaitGroup

	// Тестирование рабочего плагина
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {

			defer wg.Done()
			p, err := pm.GetPlugin("demo")
			if err != nil {
				log.Printf("Goroutine %d error: %v", id, err)
				return
			}
			log.Printf("Goroutine %d: %s", id, p.Execute())

		}(i)
	}

	// Тестирование плагина с ошибкой
	for i := 5; i < 7; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, err := pm.GetPlugin("broken")
			if err != nil {
				log.Printf("Goroutine %d error: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

}
