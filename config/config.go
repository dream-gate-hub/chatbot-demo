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
	RoleSetting struct {
		MaxTextLength int
	}
}

var ConfigInstance *Config

func InitConfig() {
	viper.SetConfigName("config")                   // 配置文件名称（无扩展名）
	viper.SetConfigType("yaml")                     // 配置文件类型
	viper.AddConfigPath(getCurrentAbPathByCaller()) // 查找配置文件的路径
	err := viper.ReadInConfig()                     // 读取配置数据
	if err != nil {
		panic(err)
	}

	ConfigInstance.Mysql.Host = viper.GetString("mysql.host")
	ConfigInstance.Mysql.Port = viper.GetString("mysql.port")
	ConfigInstance.Mysql.Username = viper.GetString("mysql.username")
	ConfigInstance.Mysql.Password = viper.GetString("mysql.password")
	ConfigInstance.Mysql.DBname = viper.GetString("mysql.dbname")

	ConfigInstance.Openai.Apikey = viper.GetString("openai.apikey")

	ConfigInstance.RoleSetting.MaxTextLength = viper.GetInt("role_setting.max_text_length")
}

func GetConfigInstance() *Config {
	return ConfigInstance
}

func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
