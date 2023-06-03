package app

var (
	config *Config
)

type initializer struct {
	o   any
	err error
}

func (ini initializer) initialize(fn func() (any, error)) any {
	if ini.err == nil {
		ini.o, ini.err = fn()
	}
	return ini.o
}

// func prepareEnv(ctx *cli.Context) error {
// 	// 全局参数指定的配置文件
// 	configFile := ctx.String("c")
//
// 	ini := initializer{}
// 	// 解析配置文件
// 	config = ini.initialize(func() (any, error) {
// 		return LoadConf(configFile)
// 	}).(*Config)
//
// 	if ini.err != nil {
// 		return fmt.Errorf("环境初始化失败，错误: %v", ini.err)
// 	}
// 	return nil
// }
