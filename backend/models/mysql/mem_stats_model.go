package mysql

// MemStats tbl_mem_stats表
type MemStats struct {
	ID         int   `xorm:"id not null pk autoincr comment('id') INT(10)"`
	TestID     int   `xorm:"test_id not null comment('test id') index(idx_test_id_ts) INT(10)"`
	Total      int64 `xorm:"total not null BIGINT(20)"`
	Used       int64 `xorm:"used not null BIGINT(20)"`
	Cached     int64 `xorm:"cached not null BIGINT(20)"`
	Free       int64 `xorm:"free not null BIGINT(20)"`
	Active     int64 `xorm:"active not null BIGINT(20)"`
	Inactive   int64 `xorm:"inactive not null BIGINT(20)"`
	SwapTotal  int64 `xorm:"swap_total not null BIGINT(20)"`
	SwapUsed   int64 `xorm:"swap_used not null BIGINT(20)"`
	SwapFree   int64 `xorm:"swap_free not null BIGINT(20)"`
	CreateTime int64 `xorm:"create_time not null index(idx_test_id_ts) BIGINT(20)"`
}

// TableName 获取数据库对应表名
func (m MemStats) TableName() string {
	return "tbl_mem_stats"
}

// Insert 将当前cpu统计信息落库
func (m *MemStats) Insert() (int64, error) {
	sess, err := newTestDBSess()
	if err != nil {
		return 0, err
	}

	return sess.Insert(m)
}

// QueryMemStats 查询测试期间全部的内存数据采集
func QueryMemStats(tid uint32) ([]MemStats, error) {
	sess, err := newTestDBSess()
	if err != nil {
		return nil, err
	}

	m := MemStats{}
	rows := make([]MemStats, 0)
	id := 0
	readCnt := 100
	// 暂时简单处理，读完为止
	for readCnt == 100 {
		tmp := make([]MemStats, 0)
		err := sess.Table(m).
			Where("test_id=? AND id>?", tid, id).
			Limit(100).
			OrderBy("id ASC").
			Find(&tmp)
		if err != nil {
			return nil, err
		}
		readCnt = len(tmp)
		rows = append(rows, tmp...)
		if len(rows) > 0 {
			id = rows[len(rows)-1].ID
		}
	}

	return rows, err
}
