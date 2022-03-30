package config

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func init() {
	config.Store(Config)
}

var cfgmu sync.Mutex

// Atomic Value Storage
var config atomic.Value

var Config Configuration

// Set (Upload) configuration data to the Atomic in-memory storage
func setConfig(data *Configuration) {
	// Lock for update operations
	cfgmu.Lock()
	if data == nil {
		//log.Println("[WARNING ] Channels Config Data is nil....ignoring....")
		cfgmu.Unlock()
		return
	}
	config.Store(*data)
	cfgmu.Unlock()
	//log.Println("[INFO    ] Main Configuration Storage updated....")
}

func GetConfig() Configuration {
	return config.Load().(Configuration)
}

func Loader(configFile string) (err error) {
	// Load Main Configuration
	cfgData, err := configLoader(configFile)
	if err != nil {
		return
	}
	// Init Runtime
	initRuntime(cfgData)
	// Merge Defaults into parsed configuration
	//data, err := setDefaults(defaults, cfgData)
	//if err != nil {
	//	return
	//}
	// Set config Data to Atomic Storage
	setConfig(cfgData)
	return
}

func (c Configuration) GetListenConfig() string {
	return fmt.Sprintf("%s:%d", c.Listen.Host, c.Listen.Port)
}

func (c Configuration) GetCors() Cors {
	return c.Cors
}
