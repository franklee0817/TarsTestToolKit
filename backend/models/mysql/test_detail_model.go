package mysql

import (
	"github.com/go-sql-driver/mysql"
)

type PerfTestDetail struct {
	ID         int     `xorm:"id not null pk autoincr INT(10)"`
	TestID     int     `xorm:"test_id not null comment('测试ID') unique(uniq_test) INT(10)"`
	QPS        int     `xorm:"qps not null default 0 comment('QPS') INT(10)"`
	Total      int     `xorm:"total not null default 0 comment('总请求量') INT(10)"`
	Succ       int     `xorm:"succ not null default 0 comment('成功请求量') INT(10)"`
	Failed     int     `xorm:"failed not null default 0 comment('失败请求量') INT(10)"`
	SuccRate   int     `xorm:"succ_rate not null default 0 comment('成功率（放大一万倍），10000代表100%') INT(10)"`
	CostMax    float32 `xorm:"cost_max not null default 0 comment('最大耗时(ms)') FLOAT"`
	CostMin    float32 `xorm:"cost_min not null default 0 comment('最小耗时(ms)') FLOAT"`
	CostAvg    float32 `xorm:"cost_avg not null default 0 comment('平均耗时(ms)') FLOAT"`
	P90        float32 `xorm:"p90 not null default 0 comment('P90（百分之90的请求在多长时间内完成）') FLOAT"`
	P99        float32 `xorm:"p99 not null default 0 comment('P99（百分之99的请求在多长时间内完成）') FLOAT"`
	P999       float32 `xorm:"p999 not null default 0 comment('P999（百分之999的请求在多长时间内完成）') FLOAT"`
	Send       int64   `xorm:"send not null default 0 comment('发送字节数') BIGINT(20)"`
	Recv       int64   `xorm:"recv not null default 0 comment('接收字节数') BIGINT(20)"`
	CostMap    string  `xorm:"cost_map comment('响应耗时和耗时计数的map关系，0~10ms:10, 10~30ms:20') TEXT"`
	RetMap     string  `xorm:"ret_map comment('响应code和code次数的map关系, 0:1000, 400:50') TEXT"`
	CreateTime int     `xorm:"create_time not null default 0 comment('采集时间') unique(uniq_test)  INT(10)"`
}

// TableName 获取数据库对应表名
func (m PerfTestDetail) TableName() string {
	return "tbl_perf_test_detail"
}

// Insert 将当前cpu统计信息落库
func (m *PerfTestDetail) Insert() (int64, error) {
	sess, err := newTestDBSess()
	if err != nil {
		return 0, err
	}

	return sess.Insert(m)
}

// InsertTestDetail 写入测试详细信息，当finished为true时更新测试记录表
func InsertTestDetail(m *PerfTestDetail, finished bool) error {
	sess, err := newTestDBSess()
	if err != nil {
		return err
	}

	err = sess.Begin()
	if err != nil {
		return err
	}
	_, err = sess.Insert(m)
	if err != nil {
		e, ok := err.(*mysql.MySQLError)
		// duplicate entry test, ignore duplicate
		if ok && e.Number == 1062 {
			err = nil
		} else {
			_ = sess.Rollback()
			return err
		}
	}
	if finished {
		perf := new(PerfTests)
		perf.ID = m.TestID
		perf.Finished = 1
		_, err = sess.ID(perf.ID).Update(perf)
		if err != nil {
			_ = sess.Rollback()
			return err
		}
	}

	return sess.Commit()
}

// QueryTestDetail 查询测试详情
func QueryTestDetail(tid, timestamp uint32) ([]PerfTestDetail, error) {
	sess, err := newTestDBSess()
	if err != nil {
		return nil, err
	}

	m := PerfTestDetail{}
	rows := make([]PerfTestDetail, 0)
	id := 0
	readCnt := 100
	// 暂时简单处理，读完为止
	for readCnt == 100 {
		tmp := make([]PerfTestDetail, 0)
		err := sess.Table(m).
			Where("test_id=? AND id>? AND create_time>?", tid, id, timestamp).
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

	return rows, nil
}
