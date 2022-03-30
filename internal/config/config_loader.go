package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

func configLoader(configFile string) (cfg *Configuration, err error) {
	// Read Yaml Config
	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("%s", configFile))
	if err != nil {
		return nil, fmt.Errorf("Unable to read config file: %s", err)
	}

	var fileCfg Configuration
	err = yaml.Unmarshal(yamlFile, &fileCfg)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse config file: %s", err)
	}

	// merge config from file to defaults to build final config
	cfg = configDefaults
	err = mergo.Merge(cfg, fileCfg, mergo.WithOverride)
	if err != nil {
		return
	}
	return
}

func initRuntime(cfg *Configuration) {
	if os.Getenv("GOGC") != "" {
		goGc, err := strconv.Atoi(os.Getenv("GOGC"))
		if err != nil {
			log.Printf("[ERROR   ] Unable to parse GOGC variable to int: %s", err.Error())
		} else {
			debug.SetGCPercent(goGc)
			//log.Printf("[INFO    ] Set Runtime GOGC=%s from Environmental variable.", os.Getenv("GOGC"))
		}
	} else {
		debug.SetGCPercent(cfg.Runtime.GOGC)
	}

	if os.Getenv("GOMAXPROCS") != "" {
		goMaxProcs, err := strconv.Atoi(os.Getenv("GOMAXPROCS"))
		if err != nil {
			log.Printf("[ERROR   ] Unable to parse GOMAXPROCS variable to int: %s", err.Error())
		} else {
			runtime.GOMAXPROCS(goMaxProcs)
			//log.Printf("[INFO    ] Set Runtime GOMAXPROCS=%s from Environmental variable.", os.Getenv("GOMAXPROCS"))
		}
	} else {
		numCPU := cfg.Runtime.GOMAXPROCS
		if numCPU == 0 {
			runtime.GOMAXPROCS(numCPU)
		} else {
			runtime.GOMAXPROCS(runtime.NumCPU())
		}
		//log.Printf("[INFO    ] Set Runtime GOMAXPROCS=%d from configuration.", runtime.NumCPU())
	}
}
