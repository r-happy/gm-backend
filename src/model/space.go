package model

import (
	"time"
)

type Space struct {
	ID         string    `json:"id"`           // プライマリキー
	ParentID   string    `json:"parent_id"`    // 親ID
	SpaceName  string    `json:"space_name"`   // スペース名
	TimeOfBorn time.Time `json:"time_of_born"` // 作成日時、自動生成
}

type Member struct {
	Space string `json:"space_id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
}

type ResponsibleUid struct {
	Space          string `json:"space_name"`
	ResponsibleUid string `json:"responsible_uid"`
}

func CreateSpace(space *Space) {
	db.Create(space)
}

func CreateMember(member *Member) {
	db.Create(member)
}

func FindSpace(space *Space) Space {
	var s Space
	db.Where(space).First(&s)
	return s
}

func FindSpaces(space *Space) []Space {
	var s []Space
	db.Where(space).Find(&s)
	return s
}

func FindMembers(member *Member) []Member {
	var m []Member
	db.Where(member).Find(&m)
	return m
}

func SaveMember(member *Member) Member {
	db.Model(&Member{}).Where("email = ?", member.Email).Updates(map[string]interface{}{
		"admin": member.Admin,
	})
	return *member
}
