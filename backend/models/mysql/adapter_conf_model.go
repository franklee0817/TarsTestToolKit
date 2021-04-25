package mysql

import (
	"time"

	"github.com/TarsCloud/TarsGo/tars"

	"github.com/franklee0817/t3k/backend/constants/errors"
)

// AdapterConf tars数据库的微服务运行时配置表
type AdapterConf struct {
	ID           int       `xorm:"id not null pk autoincr INT(11)"`
	Application  string    `xorm:"application default '' unique(application) VARCHAR(50)"`
	ServerName   string    `xorm:"server_name default '' unique(application) VARCHAR(128)"`
	NodeName     string    `xorm:"node_name default '' unique(application) VARCHAR(50)"`
	AdapterName  string    `xorm:"adapter_name default '' unique(application) VARCHAR(100)"`
	RegTs        time.Time `xorm:"registry_timestamp not null default 'CURRENT_TIMESTAMP(3)' index DATETIME(3)"`
	ThreadNum    int       `xorm:"thread_num default 1 INT(11)"`
	Endpoint     string    `xorm:"endpoint default '' index VARCHAR(128)"`
	MaxConn      int       `xorm:"max_connections default 1000 INT(11)"`
	AllowIP      string    `xorm:"allow_ip not null default '' VARCHAR(255)"`
	Servant      string    `xorm:"servant default '' VARCHAR(128)"`
	QueueCap     int       `xorm:"queuecap INT(11)"`
	QueueTimeout int       `xorm:"queuetimeout INT(11)"`
	PostTime     time.Time `xorm:"posttime DATETIME"`
	LastUser     string    `xorm:"lastuser VARCHAR(30)"`
	Proto        string    `xorm:"protocol default 'tars' VARCHAR(64)"`
	HandleGroup  string    `xorm:"handlegroup default '' VARCHAR(64)"`
}

// TableName 获取表名
func (m AdapterConf) TableName() string {
	return "t_adapter_conf"
}

// SetAdapterThread 更新目标服务的运行线程数
func SetAdapterThread(adapterName string, thread int) error {
	sess, err := newTarsDBSess()
	if err != nil {
		return err
	}

	m := &AdapterConf{}
	_, err = sess.Table(m).
		Where("adapter_name=?", adapterName).
		Cols("id").
		Get(m)
	if err != nil {
		return tars.Errorf(errors.ErrMysqlQueryFailed, "no adapter named %q found", adapterName)
	}
	if m.ThreadNum == thread {
		return nil
	}

	m.ThreadNum = thread
	m.PostTime = time.Now()
	_, err = sess.ID(m.ID).Update(m)
	if err != nil {
		return tars.Errorf(errors.ErrMysqlUpdateFailed, "failed to update adapter %q", adapterName)
	}

	return nil
}
