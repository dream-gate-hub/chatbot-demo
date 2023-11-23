package routes

import (
	"chatbot/chatgpt"
	"chatbot/config"
	"chatbot/db"
	"chatbot/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func InitRouter() *gin.Engine {
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

	return r
}

// 首页
func IndexHandler(c *gin.Context) {
	//
}

// 获取全部角色信息
func GetAllCharacterInfoHandler(c *gin.Context) {
	characters, err := db.GetAllCharacters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error", "error": err.Error()})
		return
	}
	fmt.Println(characters)
	c.JSON(http.StatusOK, gin.H{"characters_list": characters})
	return
}

// 创建角色
func CreateCharacterHandler(c *gin.Context) {
	//
	var character model.Character
	err := c.ShouldBind(&character)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("request err: %s", err.Error()))
		return
	}

	if !isLegalCharacterForCreate(&character) {
		fmt.Println(character)
		c.JSON(http.StatusBadRequest, "Illegal request data")
		return
	}

	character.CharacterID = 0 // 确保使用数据库的自增ID
	err = db.CreateCharacter(&character)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"character_id": character.CharacterID})
	return
}

// 删除角色
func DeleteCharacterHandler(c *gin.Context) {
	cidStr := c.Param("cid")
	cid, err := strconv.Atoi(cidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 查cid是否存在，如果不存在，返回错误
	//_, err = GetCharacterByID(TableCharactersHandler, uint(cid))
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}

	err = db.DeleteCharacter(uint(cid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Character deleted successfully"})
	return

}

// 更新角色
func UpdateCharacterHandler(c *gin.Context) {
	cidStr := c.Param("cid")
	cid, err := strconv.Atoi(cidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查cid是否存在，如果不存在，返回错误
	_, err = db.GetCharacterByID(uint(cid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var character model.Character
	err = c.ShouldBind(&character)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("request err: %s", err.Error()))
		return
	}

	if !isLegalCharacterForCreate(&character) {
		c.JSON(http.StatusBadRequest, "Illegal request data")
		return
	}

	character.CharacterID = uint(cid)

	err = db.UpdateCharacter(&character)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Character updated successfully"})

}

// 获取角色信息
func GetCharacterInfoHandler(c *gin.Context) {
	cidStr := c.Param("cid")
	cid, err := strconv.Atoi(cidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	character, err := db.GetCharacterByID(uint(cid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"character": character})
	return

}

// 与角色聊天
func ChatHandler(c *gin.Context) {
	cidStr := c.Param("cid")
	cid, err := strconv.Atoi(cidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	character, err := db.GetCharacterByID(uint(cid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var chatReq chatgpt.MessagesRequest
	err = c.ShouldBind(&chatReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	systemPrompt := OrganizeSystemPrompt(character)

	resp, err := chatgpt.ChatWithGPT3_5(systemPrompt, chatReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": resp})
	return

}

// 用于在创建和更新角色的时候检查角色是否合法
// 后面这里可以做例如敏感名称过滤等
func isLegalCharacterForCreate(character *model.Character) bool {
	if len(character.Gender) == 0 || len(character.Nickname) == 0 || len(character.CharacterSetting) == 0 {
		return false
	}

	if len(character.Nickname) > 255*4 { // utf8mb4 nickname的mysql字段类型是VARCHAR(255)
		return false
	}

	if len(character.CharacterSetting) > config.GetConfigInstance().RoleSetting.MaxTextLength ||
		len(character.Prologue) > config.GetConfigInstance().RoleSetting.MaxTextLength ||
		len(character.DialogueExamples) > config.GetConfigInstance().RoleSetting.MaxTextLength {
		return false
	}

	return true
}

func OrganizeSystemPrompt(character *model.Character) string {
	return fmt.Sprintf(`
[Your name]
%s

[Your character's setting]
%s

[Your prologue]
%s

[Dialogue examples]
%s
`, character.Nickname, character.CharacterSetting, character.Prologue, character.DialogueExamples)

}
