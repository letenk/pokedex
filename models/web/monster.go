package web

import (
	"mime/multipart"

	"github.com/letenk/pokedex/models/domain"
)

type MonsterQueryRequest struct {
	Name    string   `form:"name"`
	Types   []string `form:"types"`
	Catched string   `form:"catched"`
	Sort    string   `form:"sort"`
	Order   string   `form:"order"`
}

type MonsterCreateRequest struct {
	Name        string                `json:"name" form:"name" binding:"required"`
	CategoryID  string                `json:"category_id" form:"category_id" binding:"required"`
	Description string                `json:"description" form:"description" binding:"required"`
	Length      float32               `json:"length" form:"length" binding:"required"`
	Weight      uint16                `json:"weight" form:"weight" binding:"required"`
	Hp          uint16                `json:"hp" form:"hp" binding:"required"`
	Attack      uint16                `json:"attack" form:"attack" binding:"required"`
	Defends     uint16                `json:"defends" form:"defends" binding:"required"`
	Speed       uint16                `json:"speed" form:"speed" binding:"required"`
	Image       *multipart.FileHeader `form:"image" binding:"required"`
	TypeID      []string              `json:"type_id" form:"type_id" binding:"required"`
}

type MonstersResponseList struct {
	ID         string                `json:"id"`
	Name       string                `json:"name"`
	CategoryID string                `json:"category_name"`
	Catched    bool                  `json:"catched"`
	Image      string                `json:"image"`
	Types      []MonsterTypeResponse `json:"types"`
}

type MonsterResponseDetail struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	CategoryID  string                `json:"category_name"`
	Description string                `json:"description"`
	Length      float32               `json:"length"`
	Weight      uint16                `json:"weight"`
	Hp          uint16                `json:"hp"`
	Attack      uint16                `json:"attack"`
	Defends     uint16                `json:"defends"`
	Speed       uint16                `json:"speed"`
	Catched     bool                  `json:"catched"`
	Image       string                `json:"image"`
	Types       []MonsterTypeResponse `json:"types"`
}

type MonsterTypeResponse struct {
	Name string `json:"name"`
}

func FormatMonsterResponseList(monsters []domain.Monster) []MonstersResponseList {
	// If monster is empty, return empty slice
	if len(monsters) == 0 {
		return []MonstersResponseList{}
	}

	var formatters []MonstersResponseList

	for _, data := range monsters {
		formatter := MonstersResponseList{}
		formatter.ID = data.ID
		formatter.Name = data.Name
		formatter.CategoryID = data.Category.Name
		formatter.Catched = data.Catched
		formatter.Image = data.Image

		monsterTypes := []MonsterTypeResponse{}
		for _, t := range data.Types {
			typeResponse := MonsterTypeResponse{}
			typeResponse.Name = t.Name
			monsterTypes = append(monsterTypes, typeResponse)
		}
		formatter.Types = monsterTypes

		formatters = append(formatters, formatter)
	}

	return formatters
}

// Format for handle single response monster detail
func FormatMonsterResponseDetail(monster domain.Monster) MonsterResponseDetail {

	formatter := MonsterResponseDetail{}
	formatter.ID = monster.ID
	formatter.Name = monster.Name
	formatter.CategoryID = monster.Category.Name
	formatter.Description = monster.Description
	formatter.Length = monster.Length
	formatter.Weight = monster.Weight
	formatter.Hp = monster.Hp
	formatter.Attack = monster.Attack
	formatter.Defends = monster.Defends
	formatter.Speed = monster.Speed
	formatter.Catched = monster.Catched
	formatter.Image = monster.Image

	monsterTypes := []MonsterTypeResponse{}
	for _, t := range monster.Types {
		typeResponse := MonsterTypeResponse{}
		typeResponse.Name = t.Name
		monsterTypes = append(monsterTypes, typeResponse)
	}
	formatter.Types = monsterTypes

	return formatter
}
