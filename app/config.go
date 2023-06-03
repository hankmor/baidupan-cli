package app

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	BaiduPan BaiduPanConfig `mapstructure:"baidu-pan" yaml:"baidu-pan"`
}

type BaiduPanConfig struct {
	AppId     string `mapstructure:"app-id" yaml:"app-id"`
	AppKey    string `mapstructure:"app-key" yaml:"app-key"`
	SecretKey string `mapstructure:"secret-key" yaml:"secret-key"`
	SignKey   string `mapstructure:"sign-key" yaml:"sign-key"`
}

func LoadConf(f string) (*Config, error) {
	fmt.Printf("正在加载配置文件：%s...\n", f)

	var cfg = &Config{}
	var v = viper.New()
	v.SetConfigFile(f)
	v.SetConfigType("yaml")
	var err = v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("读取配置文件错误: %v", err)
	}
	// 监听配置文件是否变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("配置文件已更改，正在更新...\n")
		// 重新解析配置文件到结构体，解析错误不 panic
		if err = v.Unmarshal(cfg); err != nil {
			fmt.Printf("更新配置文件错误: %v\n", err)
		}
	})
	// 解析配置文件到结构体
	if err = v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件错误: %v\n", err)
	}
	return cfg, nil
}
