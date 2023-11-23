package main

import (
	"chatbot/config"
	"chatbot/db"
	"chatbot/routes"
	"fmt"
	"gorm.io/gorm"
)

var TableCharactersHandler *gorm.DB

//var AllConfig *config.Config

func main() {
	config.InitConfig()
	ci := config.GetConfigInstance()

	// 初始化数据库
	err := db.InitMysql(
		ci.Mysql.Username,
		ci.Mysql.Password,
		ci.Mysql.Host,
		ci.Mysql.Port,
		ci.Mysql.DBname)
	if err != nil {
		fmt.Println("init mysql err:", err.Error())
		return
	}

	r := routes.InitRouter()

	// 自签证书 仅用于测试
	go func() {
		err := r.RunTLS(":443", `./self_signed_cert/cert.pem`, `./self_signed_cert/key.pem`)
		if err != nil {
			panic(err)
		}
	}()

	err = r.Run(":8080")
	if err != nil {
		fmt.Println("run err:", err.Error())
	}
}
