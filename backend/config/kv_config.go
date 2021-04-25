package config

import (
	"sync"
	"time"

	"github.com/TarsCloud/TarsGo/tars"
	"gopkg.in/yaml.v3"

	"github.com/franklee0817/t3k/backend/constants"
)

type kvConf struct {
	kv  map[string]string
	mux *sync.RWMutex
}

var kvCfg *kvConf

func init() {
	r := new(sync.Once)
	r.Do(func() {
		kvCfg = &kvConf{
			kv:  make(map[string]string),
			mux: new(sync.RWMutex),
		}
		watchKVCfg()
	})
}

func watchKVCfg() {
	loadKVCfg()
	time.AfterFunc(constants.CfgReloadDur, watchKVCfg)
}

func loadKVCfg() {
	servCfg := tars.GetServerConfig()
	remoteConf := tars.NewRConf(servCfg.App, servCfg.Server, servCfg.BasePath)
	config, err := remoteConf.GetConfig("kv.yaml")
	if err == nil {
		var cfg = make(map[string]string)
		err := yaml.Unmarshal([]byte(config), &cfg)
		if err != nil {
			return
		}
		kvCfg.mux.Lock()
		kvCfg.kv = cfg
		kvCfg.mux.Unlock()
	}
}

// GetKV 读取kv配置
func GetKV(key string) string {
	kvCfg.mux.RLock()
	defer kvCfg.mux.RUnlock()

	return kvCfg.kv[key]
}
