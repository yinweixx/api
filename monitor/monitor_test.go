package monitor

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMonitorCron(t *testing.T) {
	monitor := NewGatewayMonitor()
	monitor.Crontab = []Crontab{
		Crontab{
			Name:       "test",
			Expression: "*/3 * * * * *",
			Func: func() {
				fmt.Println("Time:", time.Now())
			},
		},
	}
	monitor.StartTasks(false)
	var a sync.WaitGroup
	a.Add(1)
	a.Wait()
}
