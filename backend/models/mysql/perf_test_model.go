package mysql

import "errors"

type PerfTests struct {
	ID        int    `xorm:"id not null pk autoincr INT(10)"`
	ServType  string `xorm:"serv_type comment('服务器类型') VARCHAR(255)"`
	Lang      string `xorm:"lang comment('被测语言') VARCHAR(16)"`
	ServName  string `xorm:"serv_name not null comment('被测服务名称：TestUnits.cpp.TestObj') VARCHAR(64)"`
	FnName    string `xorm:"fn_name comment('被测接口') VARCHAR(64)"`
	Cores     int    `xorm:"cores comment('服务器核数') INT(10)"`
	Threads   int    `xorm:"threads comment('线程数') INT(10)"`
	ConnCnt   int    `xorm:"conn_cnt comment('每个节点的连接数') INT(10)"`
	Frequency int    `xorm:"frequency comment('每个节点的请求速率') INT(10)"`
	KeepAlive int    `xorm:"keep_alive comment('测试时长') INT(10)"`
	PkgLen    int    `xorm:"pkg_len comment('压测包大小：0K,1K...') INT(10)"`
	StartTime int    `xorm:"start_time not null comment('开始时间') index INT(10)"`
	EndTime   int    `xorm:"end_time not null comment('结束时间') INT(10)"`
	Finished  int    `xorm:"finished not null default 0 comment('') TINYINT(1)"`
}

// TableName 获取数据库对应表名
func (m PerfTests) TableName() string {
	return "tbl_perf_tests"
}

// Insert 将当前cpu统计信息落库
func (m *PerfTests) Insert() (int64, error) {
	sess, err := newTestDBSess()
	if err != nil {
		return 0, err
	}

	return sess.Insert(m)
}

// PaginatePerfTests 分页查看PerfTests
func PaginatePerfTests(page, pageSize uint32) (int64, []PerfTests, error) {
	sess, err := newTestDBSess()
	if err != nil {
		return 0, nil, err
	}

	model := PerfTests{}
	rows := make([]PerfTests, 0)
	total, err := sess.Table(model).
		Limit(int(pageSize), int((page-1)*pageSize)).
		FindAndCount(&rows)

	return total, rows, err
}

// GetPerfTest 查询单条历史记录
func GetPerfTest(tid uint32) (*PerfTests, error) {
	sess, err := newTestDBSess()
	if err != nil {
		return nil, err
	}

	row := PerfTests{ID: int(tid)}
	found, err := sess.Table(row).Get(&row)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("no matching row found")
	}

	return &row, nil
}
