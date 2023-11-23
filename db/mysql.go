package db

import (
	"chatbot/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBInstance *gorm.DB

func InitMysql(user, password, host, port, dbname string) error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)
	//dsn := "root:123456@tcp(127.0.0.1:3306)/mytestdb?charset=utf8mb4&parseTime=True&loc=Local"
	dialector := mysql.Open(dsn)
	DBInstance, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

func CreateCharacter(character *model.Character) error {
	newDB := DBInstance.Session(&gorm.Session{}).Table("characters")
	return newDB.Create(character).Error
}

// 通过 ID 查询
func GetCharacterByID(cid uint) (*model.Character, error) {
	var character model.Character

	newDB := DBInstance.Session(&gorm.Session{}).Table("characters")
	err := newDB.Where("character_id = ?", cid).First(&character).Error

	//err := db.Raw("SELECT * FROM characters WHERE character_id = ?", cid).Scan(&character).Error
	if err != nil {
		return nil, err
	}
	return &character, nil
}

// 查询所有记录
func GetAllCharacters() ([]model.Character, error) {
	var characters []model.Character

	//err := db.Raw("SELECT * FROM characters").Scan(&characters).Error
	//if err != nil {
	//	return nil, err
	//}

	newDB := DBInstance.Session(&gorm.Session{}).Table("characters")
	err := newDB.Find(&characters).Error

	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出到哪里）
	//	logger.Config{
	//		LogLevel: logger.Info, // 日志级别
	//		// 其他配置选项...
	//	},
	//)
	//err := db.Session(&gorm.Session{Logger: newLogger}).Find(&characters).Error
	//
	if err != nil {
		return nil, err
	}
	return characters, nil
}

// 更新character，若没有对应的character，则新增
func UpdateCharacter(character *model.Character) error {
	newDB := DBInstance.Session(&gorm.Session{}).Table("characters")
	return newDB.Save(character).Error
}

// 根据character_id删除一条记录，如果character_id不存在也不会返回错误
func DeleteCharacter(cid uint) error {
	newDB := DBInstance.Session(&gorm.Session{}).Table("characters")
	return newDB.Where("character_id = ?", cid).Delete(&model.Character{}).Error
}
