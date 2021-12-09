package mysql

import (
	"fmt"
	"testplay/utils"
)

type User struct {
	Id           int      `json:"id" xorm:"not null pk autoincr INT(11)"`
	Username     string   `json:"username" xorm:""`
	Profession   string   `json:"profession" xorm:""`
	Email        string   `json:"email" xorm:""`
	TextInfo     string   `json:"text_info" xorm:""`
	RegisterTime int64    `json:"register_time" xorm:""`
	Attribute    []string `json:"attribute" xorm:""`
	Pets         []*Pet   `json:"pets" xorm:""`
}

func (u *User) TableName() string {
	return "user"
}

type Pet struct {
	PetName   string   `json:"pet_name" xorm:""`
	PetSex    int8     `json:"pet_sex" xorm:""`
	PetAttack int      `json:"pet_attack" xorm:""`
	PetTag    []string `json:"pet_tag" xorm:""`
}

func TmpAddMany(us []*User) (err error) {
	_, err = utils.Engine.Table("user").InsertMulti(us)
	return
}

func TmpAdd(u *User) (id int64, err error) {
	id, err = utils.Engine.Insert(u)
	fmt.Println(id, err, u)
	return
}
