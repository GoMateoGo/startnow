package config

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"second_hand_mall/internal/core/initialize/db"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	HttpServer HttpServer
	RpcServer  RpcServer
	Db         []db.LinkInfo // 数据库连接信息
	Zap        Zap
	JWT        JWT
}

// 加载本地配置文件
func LoadLocalConfig(filePath string) Config {

	// 获取最后一级名称（包含后缀）
	fileName := filepath.Base(filePath)
	// 获取去掉后缀的文件名称
	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	// 获取最后一级前面的路径
	dirPath := filepath.Dir(filePath)

	viper.SetConfigName(fileNameWithoutExt) // 配置文件名（不包含扩展名）
	viper.SetConfigType("yaml")             // 配置文件类型
	viper.AddConfigPath(dirPath)            // 配置文件所在路径
	//viper.WatchConfig()                     // 开始监听配置文件变化
	//viper.OnConfigChange(changeConfig)           // 配置改变Hook
	if err := viper.ReadInConfig(); err != nil { // 读取配置文件
		log.Panicf("加载本地配置文件错误:%v", err)
	}

	var cfg Config

	// 将读取到的数据反序列化为配置变量
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Panicf("反序列化配置文件失败:%v", err)
	}
	jsonCfg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Printf("转换配置为JSON格式失败: %v", err)
	}
	fmt.Println(string(jsonCfg))

	return cfg
}
