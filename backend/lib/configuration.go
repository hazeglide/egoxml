package egoxml

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

type Configuration struct {
	Savegame  string
	Trigger   []string
	Exclude   []string
	Combine   map[string][]string
	ShipIndex int
}

var config Configuration
var mutex sync.RWMutex

func init() {
	loadConfiguration()
	watchConfiguration()
}
func loadConfiguration() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}
}

func watchConfiguration() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		mutex.Lock()
		loadConfiguration()
		mutex.Unlock()
	})
}

func GetConfig() Configuration {
	mutex.RLock()
	defer mutex.RUnlock()
	return config
}
