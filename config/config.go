package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var envPtr = pflag.String("env", "dev", "Environment: dev or prod")

// AllConfig 整合Config
type AllConfig struct {
	Server Server
	//DataSource DataSource
	//Redis      Redis
	Log Log
	Jwt Jwt
	//AliOss     AliOss
	//Wechat     Wechat
}

type Server struct {
	Port  string
	Level string
}

type Log struct {
	Level    string
	FilePath string
}

type Jwt struct {
	Admin JwtOption
	User  JwtOption
}

type JwtOption struct {
	Secret string
	TTL    string
	Name   string
}

func InitLoadConfig() *AllConfig {
	pflag.Parse()
	config := viper.New()
	config.AddConfigPath("./config")
	config.SetConfigName(fmt.Sprintf("application-%s", *envPtr))
	config.SetConfigType("yaml")

	var configData *AllConfig
	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Use Viper ReadInConfig Fatal error config err:%s \n", err))
	}
	err = config.Unmarshal(&configData)
	if err != nil {
		panic(fmt.Errorf("Use Viper Unmarshal Fatal error config err:%s \n", err))
	}

	fmt.Printf("配置文件信息：%+v", configData) // 打印配置文件信息
	return configData
}
