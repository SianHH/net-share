package component

import (
	"gopkg.in/yaml.v3"
	"log"
	"net-share/pkg/utils"
	"net-share/server/config"
	"net-share/server/global"
	"os"
	"path/filepath"
)

func InitConfig() {
	loadConfigFile("data/config.yaml", &global.App)
}

func loadConfigFile(configFilePath string, result any) {
	_ = os.MkdirAll(filepath.Dir(configFilePath), os.ModeDir)
	stat, err := os.Stat(configFilePath)
	if err != nil {
		log.Println("配置文件不存在", err)
		initConfigFile(configFilePath)
		return
	}
	if stat.IsDir() {
		log.Println("配置文件不能为文件夹")
		os.Exit(1)
	}
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Println("读取配置文件错误", err)
		os.Exit(1)
	}
	if err := yaml.Unmarshal(file, result); err != nil {
		log.Println("解析配置文件错误", err)
		os.Exit(1)
	}
}

func initConfigFile(configFilePath string) {
	if err := os.MkdirAll(filepath.Dir(configFilePath), 0666); err != nil {
		log.Println("初始化配置文件失败", err)
		os.Exit(1)
	}
	global.App = config.App{
		Addr:        "0.0.0.0:8080",
		Mode:        "prod",
		Account:     "admin",
		Password:    "123456",
		JwtKey:      utils.RandStr(32, utils.AllDict),
		Ip:          "127.0.0.1",
		Domain:      "example.com",
		HostPort:    "2096",
		Entrypoint:  "18080",
		ForwardPort: "2097",
		Ports: []string{
			"20001-21000",
			"30001",
			"30002",
		},
		Logger: config.Logger{
			File:    "application.log",
			Level:   0,
			Console: true,
		},
	}
	marshal, err := yaml.Marshal(global.App)
	if err != nil {
		log.Println("初始化配置文件失败", err)
		os.Exit(1)
	}
	if err := os.WriteFile(configFilePath, marshal, 0755); err != nil {
		log.Println("初始化配置文件失败", err)
		os.Exit(1)
	}
	log.Println("初始化配置文件完成")
}
