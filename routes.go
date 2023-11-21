package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Character struct {
	CharacterID      uint   `json:"character_id" gorm:"column:character_id;primaryKey;autoIncrement"`
	Nickname         string `json:"nickname" gorm:"column:nickname;size:255"`
	Gender           string `json:"gender" gorm:"column:gender;type:enum('man', 'woman', 'other')"`
	CharacterSetting string `json:"character_setting" gorm:"column:character_setting;type:text"`
	Prologue         string `json:"prologue,omitempty" gorm:"column:prologue;type:text"` // 开场白
	DialogueExamples string `json:"dialogue_examples,omitempty" gorm:"column:dialogue_examples;type:text"`
}

// 首页
func IndexHandler(c *gin.Context) {
	//
}

// 获取全部角色信息
func GetAllCharacterInfoHandler(c *gin.Context) {
	characters, err := GetAllCharacters(TableCharactersHandler)
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
	var character Character
	err := c.ShouldBind(&character)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("request err: %s", err.Error()))
		return
	}

	if !isLegalCharacterForCreate(character) {
		fmt.Println(character)
		c.JSON(http.StatusBadRequest, "Illegal request data")
		return
	}

	character.CharacterID = 0 // 确保使用数据库的自增ID
	err = CreateCharacter(TableCharactersHandler, &character)
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

	err = DeleteCharacter(TableCharactersHandler, uint(cid))
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
	_, err = GetCharacterByID(TableCharactersHandler, uint(cid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var character Character
	err = c.ShouldBind(&character)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("request err: %s", err.Error()))
		return
	}

	if !isLegalCharacterForCreate(character) {
		c.JSON(http.StatusBadRequest, "Illegal request data")
		return
	}

	character.CharacterID = uint(cid)

	err = UpdateCharacter(TableCharactersHandler, &character)
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

	character, err := GetCharacterByID(TableCharactersHandler, uint(cid))
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

	character, err := GetCharacterByID(TableCharactersHandler, uint(cid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var chatReq MessagesRequest
	err = c.ShouldBind(&chatReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	systemPrompt := OrganizeSystemPrompt(character)

	resp, err := ChatWithGPT3_5(systemPrompt, chatReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": resp})
	return

}

// 用于在创建和更新角色的时候检查角色是否合法
// 后面这里可以做例如敏感名称过滤等
func isLegalCharacterForCreate(character Character) bool {
	if len(character.Gender) == 0 || len(character.Nickname) == 0 || len(character.CharacterSetting) == 0 {
		return false
	}

	return true
}

func OrganizeSystemPrompt(character *Character) string {
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
