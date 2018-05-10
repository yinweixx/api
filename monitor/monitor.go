package monitor

import "github.com/robfig/cron"

//GatewayMonitor -- 网关监控服务结构体
type GatewayMonitor struct {
	Cron    *cron.Cron
	Crontab []Crontab
}

//Crontab -- 监控任务定义结构体
type Crontab struct {
	Name       string
	Expression string
	Func       GatewayMonitorFunc
}

//GatewayMonitorFunc -- 监控任务处理器
type GatewayMonitorFunc func()

//NewGatewayMonitor -- 创建API网关监控
func NewGatewayMonitor() *GatewayMonitor {
	return &GatewayMonitor{
		Cron: cron.New(),
	}
}

//StartTasks -- 开始监控任务
func (_this *GatewayMonitor) StartTasks(runTaskAfterAdded bool) {
	_this.Cron.Start()
	for _, crontab := range _this.Crontab {
		_this.Cron.AddFunc(crontab.Expression, crontab.Func)
	}
	if runTaskAfterAdded {
		for _, crontab := range _this.Crontab {
			crontab.Func()
		}
	}
}
