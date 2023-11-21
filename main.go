package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var TableCharactersHandler *gorm.DB

func main() {
	// 初始化数据库
	user := "root"
	password := "123456"
	host := "localhost"
	port := "3306"
	dbname := "chatbot"
	db, err := initMysql(user, password, host, port, dbname)
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

	err = r.Run(":8080")
	if err != nil {
		fmt.Println("run err:", err.Error())
	}
}
