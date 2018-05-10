package app

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

//ApplicationArgs -- 应用程序基本信息结构体
type ApplicationArgs struct {
	Name    string          //应用程序名称
	Usage   string          //应用程序描
	Version string          //应用程序版本
	Flags   []cli.Flag      //应用程序参数列表
	Context context.Context //应用程上下文
}

//InitFlagsFunc -- 参数初始化
type InitFlagsFunc func() *ApplicationArgs

//Bootstrap -- 应用程序启动入口
type Bootstrap func(*cli.Context) error

//InitBootParameters -- 初始化应用程序启动参数
func InitBootParameters(initFlagsFunc InitFlagsFunc, bootstrap Bootstrap) (*cli.App, error) {
	app := cli.NewApp()
	applicationArgs := initFlagsFunc()
	applicationArgs.Context = context.Background() //创建一个空的根应用程序上下文
	app.Name = applicationArgs.Name
	app.Usage = applicationArgs.Usage
	app.Version = applicationArgs.Version
	app.Flags = applicationArgs.Flags
	app.Action = func(ctx *cli.Context) error {
		log.SetOutput(os.Stdout)
		formatter := new(prefixed.TextFormatter)
		log.SetFormatter(formatter)
		return bootstrap(ctx)
	}
	return app, nil
}

//RunApplication --
func RunApplication(initFlagsFunc InitFlagsFunc, bootstrap Bootstrap) {
	app, err := InitBootParameters(initFlagsFunc, bootstrap)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("application boot parameters parse error")
	}
	app.Run(os.Args)
}
