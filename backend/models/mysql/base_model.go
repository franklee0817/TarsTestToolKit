package mysql

import (
	_ "github.com/go-sql-driver/mysql"

	"xorm.io/xorm"

	"github.com/franklee0817/t3k/backend/config"
)

func newSession(cfgName string) (*xorm.Session, error) {
	connStr := config.GetDBCfg(cfgName)
	engine, err := xorm.NewEngine("mysql", connStr)
	if err != nil {
		return nil, err
	}

	return engine.NewSession(), nil
}

func newTestDBSess() (*xorm.Session, error) {
	return newSession(config.MySQLTestDB)
}

func newTarsDBSess() (*xorm.Session, error) {
	return newSession(config.MySQLTarsDB)
}
