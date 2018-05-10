package server

//Status -- API网关运行状态
//* 内存使用统计信息(使用量、百分比)
//* CPU使用统计信息(使用量、百分比)
//* goroutine协程使用统计信息（会话处理协程）
type Status struct {
	CPU       struct{} //CPU使用统计信息
	Memory    struct{} //内存使用统计信息
	Goroutine struct{} //协程使用统计信息
}

const (
	//OK --
	OK = 200
	//SYSTEMERROR --
	SYSTEMERROR = 500
	//SYSTEMBUSYERROR --
	SYSTEMBUSYERROR = 500.13
	//TIMEOUT --
	TIMEOUT = 504
	//ILLAGERREQUEST --
	ILLAGERREQUEST = 403
	//NOTFIND --
	NOTFIND = 404
)
