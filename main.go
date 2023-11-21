package main

import (
	"chatbot/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var TableCharactersHandler *gorm.DB
var AllConfig *config.Config

func main() {
	AllConfig = config.GetAllConfig()

	// 初始化数据库
	db, err := initMysql(
		AllConfig.Mysql.Username,
		AllConfig.Mysql.Password,
		AllConfig.Mysql.Host,
		AllConfig.Mysql.Port,
		AllConfig.Mysql.DBname)
	if err != nil {
		fmt.Println("init mysql err:", err.Error())
		return
	}
	TableCharactersHandler = db.Table("characters")
	if TableCharactersHandler == nil {
		fmt.Println("table [characters] is nil")
		return
	}

	r := gin.Default()
	r.GET("/", IndexHandler)
	//r.POST("/create_character", CreateCharacterHandler)
	characterRouterGroup := r.Group("/characters")
	characterRouterGroup.GET("/all", GetAllCharacterInfoHandler) // 获取全部角色信息
	characterRouterGroup.POST("/create", CreateCharacterHandler) // 创建角色
	characterRouterGroup.DELETE("/:cid", DeleteCharacterHandler) // 删除角色
	characterRouterGroup.PUT("/:cid", UpdateCharacterHandler)    // 更新角色
	characterRouterGroup.GET("/:cid", GetCharacterInfoHandler)   // 获取角色信息
	characterRouterGroup.POST("/chat/:cid", ChatHandler)         // 与角色聊天

	// 自签证书 仅用于测试
	//go func() {
	//	err := r.RunTLS(":443", `./self_signed_cert/cert.pem`, `./self_signed_cert/key.pem`)
	//	if err != nil {
	//		panic(err)
	//	}
	//}()

	err = r.Run(":8080")
	if err != nil {
		fmt.Println("run err:", err.Error())
	}
}
