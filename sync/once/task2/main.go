package main

import (
	"fmt"
	"sync"
)

type configManager struct {
	config map[string]string
	cmOnce sync.Once
}

func newConfigManager() *configManager {
	return &configManager{}
}

func (cm *configManager) LoadConfig() {
	cm.cmOnce.Do(func() {
		// Имитация загрузки конфигурации
		cm.config = map[string]string{
			"app_name":  "MyApp",
			"port":      "8080",
			"log_level": "debug",
		}
		//fmt.Println("config manager loaded")
	})
}

func (cm *configManager) Get(key string) string {
	cm.LoadConfig()
	return cm.config[key]
}

func (cm *configManager) PrintConfig() {
	cm.LoadConfig()
	for k, v := range cm.config {
		fmt.Println(k, v)
	}
}

func main() {
	cm := newConfigManager()
	keys := []string{"app_name", "port", "log_level"}
	cm.PrintConfig()
	for _, key := range keys {
		fmt.Println(cm.Get(key))
	}
}
