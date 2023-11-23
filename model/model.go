package model

type Character struct {
	CharacterID      uint   `json:"character_id" gorm:"column:character_id;primaryKey;autoIncrement"`
	Nickname         string `json:"nickname" gorm:"column:nickname;size:255"`
	Gender           string `json:"gender" gorm:"column:gender;type:enum('man', 'woman', 'other')"`
	CharacterSetting string `json:"character_setting" gorm:"column:character_setting;type:text"`
	Prologue         string `json:"prologue,omitempty" gorm:"column:prologue;type:text"` // 开场白
	DialogueExamples string `json:"dialogue_examples,omitempty" gorm:"column:dialogue_examples;type:text"`
}
