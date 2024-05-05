package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

func init() {
	// 初始化 viper

	configPath := os.Getenv("CONFIG_PATH")
	fmt.Println("configPath:", configPath)
	if configPath != "" {
		filename := fmt.Sprintf("%s/dev.yml", configPath)
		viper.SetConfigFile(filename)
	} else {
		viper.SetConfigName("dev")      // 配置文件名称(无扩展名)
		viper.SetConfigType("yml")      // 如果配置文件的名称中没有扩展名，则需要配置此项
		viper.AddConfigPath("./config") // 当前目录中查找配置文件
	}
	// 读取配置数据
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&Configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}
