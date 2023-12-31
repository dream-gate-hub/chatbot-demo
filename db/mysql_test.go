package db

import (
	"chatbot"
	"chatbot/model"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testInitDB(t *testing.T) {
	user := "root"
	password := "123456"
	host := "localhost"
	port := "3306"
	dbname := "chatbot"
	db, err := InitMysql(user, password, host, port, dbname)
	assert.Nil(t, err, "init mysql err")

	main.TableCharactersHandler = db.Table("characters")
	assert.NotNil(t, main.TableCharactersHandler, "table [characters] is nil")
}

func TestCreateCharacter(t *testing.T) {
	testInitDB(t)

	var c = &model.Character{
		Nickname:         "altman",
		Gender:           "man",
		CharacterSetting: "a cat",
		Prologue:         "hi",
		DialogueExamples: "how are you?",
	}
	err := CreateCharacter(main.TableCharactersHandler, c)
	assert.Nil(t, err, "CreateCharacter err")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func TestGetCharacterByID(t *testing.T) {
	testInitDB(t)
	c, err := GetCharacterByID(main.TableCharactersHandler, 2)
	assert.Nil(t, err, "GetCharacterByID err")

	fmt.Println(c)

}

func TestUpdateCharacter(t *testing.T) {
	testInitDB(t)
	var c = &model.Character{
		CharacterID:      6,
		Nickname:         "alter man",
		Gender:           "man",
		CharacterSetting: "---",
		Prologue:         "---",
		DialogueExamples: "---",
	}
	err := UpdateCharacter(main.TableCharactersHandler, c)
	assert.Nil(t, err, "UpdateCharacter err")
}

func TestDeleteCharacter(t *testing.T) {
	testInitDB(t)
	err := DeleteCharacter(main.TableCharactersHandler, 11)
	assert.Nil(t, err, "DeleteCharacter err")
}

func TestGetAllCharacters(t *testing.T) {
	testInitDB(t)
	allCharacters, err := GetAllCharacters(main.TableCharactersHandler)
	assert.Nil(t, err, "GetAllCharacters err")

	fmt.Println(allCharacters)

}
