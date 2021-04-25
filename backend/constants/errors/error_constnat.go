package errors

const (
	ErrCodeInvalidReq  = 1000 // 无效请求
	ErrCodeParam       = 1001 // 参数错误
	ErrCodeInternalErr = 2000 // 内部错误

	ErrMysqlQueryFailed  = 3000 // mysql 查询失败
	ErrMysqlUpdateFailed = 3001 // mysql 更新失败

	ErrCodeInvalidResp = 4001 // 无效响应

	ErrCodeCommon = 5000 // 一般错误

	BmSucc              = 0
	BmSequence          = 1
	BmException         = -1
	BmInitParam         = -101
	BmErrorUrl          = -102
	BmPacketError       = -1000
	BmPacketEncode      = -1001
	BmPacketDecode      = -1002
	BmPacketParam       = -1003
	BmPacketOverflow    = -1004
	BmSockError         = -2000
	BmSockInvalid       = -2001
	BmSockNconnected    = -2002
	BmSockConnError     = -2003
	BmSockConnTimeout   = -2004
	BmSockSendError     = -2005
	BmSockRecvError     = -2006
	BmSockRecvTimeout   = -2007
	BmShmErrGet         = -3000
	BmShmErrAttach      = -3001
	BmShmErrInit        = -3002
	BmShmErrClear       = -3003
	BmSocketErrBase     = -4000
	BmServerErrParam    = -10000
	BmAdminErrNotfind   = -10001 // 未找到此链接
	BmAdminErrRunning   = -10002 // 服务接口未运行压测
	BmAdminErrStartup   = -10003 // 启动失败
	BmAdminErrShutdown  = -10004 // 关闭失败
	BmAdminErrEncode    = -10005 // 编码失败
	BmAdminErrDecode    = -10006 // 编码失败
	BmAdminErrSocket    = -10007 // 网络失败
	BmAdminErrTask      = -10008 // 服务接口正在运行中
	BmAdminErrProto     = -10009 // 未找到此协议
	BmNodeErrRunning    = -20001 // 正在运行
	BmNodeErrResource   = -20002 // 资源不足
	BmNodeErrCasematch  = -20003 // 用例参数和内容不匹配
	BmNodeErrConnection = -20004 // 链接数设置不合理(整数倍，且不要超过500倍)
	BmNodeErrEndpoint   = -20005 // 目标服务器配置不正确
	BmNodeErrRpccall    = -20006 // 目标服务器调用失败
)
