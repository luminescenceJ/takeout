package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var envPtr = pflag.String("env", "dev", "Environment: dev or prod")

// AllConfig 整合Config
type AllConfig struct {
	Server     Server
	DataSource DataSource
	//Redis      Redis
	Log    Log
	Jwt    Jwt
	AliOss AliOss
	Path   string ``
	//Wechat     Wechat
}

type DataSource struct {
	Host     string
	Port     string
	UserName string
	Password string
	DBName   string `mapstructure:"db_name"`
	Config   string
}

func (d *DataSource) Dsn() string {
	return d.UserName + ":" + d.Password + "@tcp(" + d.Host + ":" + d.Port + ")/" + d.DBName + "?" + d.Config
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

type AliOss struct {
	EndPoint        string
	AccessKeyId     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	BucketName      string `mapstructure:"bucket_name"`
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
