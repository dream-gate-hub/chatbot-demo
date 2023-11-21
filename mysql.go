package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initMysql(user, password, host, port, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)
	//dsn := "root:123456@tcp(127.0.0.1:3306)/mytestdb?charset=utf8mb4&parseTime=True&loc=Local"
	dialector := mysql.Open(dsn)
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateCharacter(db *gorm.DB, character *Character) error {
	newDB := db.Session(&gorm.Session{})
	return newDB.Create(character).Error
}

// 通过 ID 查询
func GetCharacterByID(db *gorm.DB, cid uint) (*Character, error) {
	var character Character

	newDB := db.Session(&gorm.Session{})
	err := newDB.Where("character_id = ?", cid).First(&character).Error

	//err := db.Raw("SELECT * FROM characters WHERE character_id = ?", cid).Scan(&character).Error
	if err != nil {
		return nil, err
	}
	return &character, nil
}

// 查询所有记录
func GetAllCharacters(db *gorm.DB) ([]Character, error) {
	var characters []Character

	//err := db.Raw("SELECT * FROM characters").Scan(&characters).Error
	//if err != nil {
	//	return nil, err
	//}

	newDB := db.Session(&gorm.Session{})
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
func UpdateCharacter(db *gorm.DB, character *Character) error {
	newDB := db.Session(&gorm.Session{})
	return newDB.Save(character).Error
}

// 根据character_id删除一条记录，如果character_id不存在也不会返回错误
func DeleteCharacter(db *gorm.DB, cid uint) error {
	newDB := db.Session(&gorm.Session{})
	return newDB.Where("character_id = ?", cid).Delete(&Character{}).Error
}
