package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/TarsCloud/TarsGo/tars"
	"gopkg.in/yaml.v3"

	"github.com/franklee0817/t3k/backend/constants"
)

const (
	MySQLTestDB = "test_db"
	MySQLTarsDB = "tars_db"
)

// DBConfig 数据库连接配置
type DBConfig struct {
	DBName   string `yaml:"db_name"`
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
	Args     string `yaml:"args"`
}

// DBConfigs 数据库连接配置集合
type DBConfigs struct {
	Configs map[string]DBConfig `yaml:"configs"`
}

var dbCfg *DBConfigs

func init() {
	r := new(sync.Once)
	r.Do(func() {
		dbCfg = &DBConfigs{}
		watchDBCfg()
	})
}

func watchDBCfg() {
	loadDBCfg()
	time.AfterFunc(constants.CfgReloadDur, watchDBCfg)
}

func loadDBCfg() {
	servCfg := tars.GetServerConfig()
	remoteConf := tars.NewRConf(servCfg.App, servCfg.Server, servCfg.BasePath)
	config, err := remoteConf.GetConfig("db.yaml")
	if err == nil {
		var cfg = &DBConfigs{}
		err := yaml.Unmarshal([]byte(config), cfg)
		if err != nil {
			return
		}
		dbCfg = cfg
	}
}

// GetDBCfg 按配置名读取数据库配置
func GetDBCfg(instName string) string {
	cfg, ok := dbCfg.Configs[instName]
	if !ok {
		loadDBCfg()
		cfg = dbCfg.Configs[instName]
	}

	connStr := fmt.Sprintf("%s:%s@(%s)/%s", cfg.UserName, cfg.Password, cfg.Address, cfg.DBName)
	if len(cfg.Args) > 0 {
		connStr += "?" + cfg.Args
	}
	return connStr
}
