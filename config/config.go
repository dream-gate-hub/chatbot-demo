package config

import (
	"github.com/spf13/viper"
	"path"
	"runtime"
)

type Config struct {
	Mysql struct {
		Host     string
		Port     string
		Username string
		Password string
		DBname   string
	}
	Openai struct {
		Apikey string
	}
}

func GetAllConfig() *Config {
	viper.SetConfigName("config")                   // 配置文件名称（无扩展名）
	viper.SetConfigType("yaml")                     // 配置文件类型
	viper.AddConfigPath(getCurrentAbPathByCaller()) // 查找配置文件的路径
	err := viper.ReadInConfig()                     // 读取配置数据
	if err != nil {
		panic(err)
	}

	var c Config
	c.Mysql.Host = viper.GetString("mysql.host")
	c.Mysql.Port = viper.GetString("mysql.port")
	c.Mysql.Username = viper.GetString("mysql.username")
	c.Mysql.Password = viper.GetString("mysql.password")
	c.Mysql.DBname = viper.GetString("mysql.dbname")

	c.Openai.Apikey = viper.GetString("openai.apikey")

	return &c
}

func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
