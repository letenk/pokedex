package web

import "mime/multipart"

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
