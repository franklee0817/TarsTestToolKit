package mysql

// CpuStats tbl_cpu_stats表
type CpuStats struct {
	ID         int   `xorm:"id not null pk autoincr comment('id') INT(10)"`
	TestID     int   `xorm:"test_id not null comment('test id') index(idx_test_id_ts) INT(10)"`
	Total      int64 `xorm:"total not null comment('total time of cpu stats') BIGINT(20)"`
	Idle       int64 `xorm:"idle not null comment('cpu idle time') BIGINT(20)"`
	Used       int64 `xorm:"used not null comment('cpu used time') BIGINT(20)"`
	CreateTime int64 `xorm:"create_time not null index(idx_test_id_ts) BIGINT(20)"`
}

// TableName 获取数据库对应表名
func (m CpuStats) TableName() string {
	return "tbl_cpu_stats"
}

// Insert 将当前cpu统计信息落库
func (m *CpuStats) Insert() (int64, error) {
	sess, err := newTestDBSess()
	if err != nil {
		return 0, err
	}

	return sess.Insert(m)
}

// StoreStats 存储cpu和内存的统计数据
func StoreStats(cpu []CpuStats, mem MemStats) error {
	sess, err := newTestDBSess()
	if err != nil {
		return err
	}

	if err := sess.Begin(); err != nil {
		return err
	}

	_, err = sess.InsertMulti(&cpu)
	if err != nil {
		_ = sess.Rollback()
		return err
	}

	_, err = sess.Insert(mem)
	if err != nil {
		_ = sess.Rollback()
	}

	return sess.Commit()
}

// QueryCpuStats 查询测试期间全部是CPU数据采集
func QueryCpuStats(tid uint32) ([]CpuStats, error) {
	sess, err := newTestDBSess()
	if err != nil {
		return nil, err
	}

	m := CpuStats{}
	rows := make([]CpuStats, 0)
	id := 0
	readCnt := 100
	// 暂时简单处理，读完为止
	for readCnt == 100 {
		tmp := make([]CpuStats, 0)
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
