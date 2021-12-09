package es

import (
	"test/utils"
)

type UserES struct {
	ID           int      `json:"id"`
	Username     string   `json:"username"`
	Profession   string   `json:"profession"`
	Email        string   `json:"email"`
	TextInfo     string   `json:"text_info"`
	RegisterTime int64    `json:"register_time"`
	Attribute    []string `json:"attribute"`
	Pets         []*Pet   `json:"pets"`
}

type Pet struct {
	PetName   string   `json:"pet_name"`
	PetSex    int8     `json:"pet_sex"`
	PetAttack int      `json:"pet_attack"`
	PetTag    []string `json:"pet_tag"`
}

var petNames []string = []string{"小拳石", "喷火龙", "水箭龟", "裂空座", "古拉顿", "海皇牙", "黑", "白"}
var attribute []string = []string{"坚毅", "隐忍", "勇敢", "懦弱", "果敢"}
var petTag []string = []string{"愚笨", "智慧", "无畏", "坚韧"}
var profession []string = []string{"剑客", "刀客", "弓箭手", "突击手", "忍者"}

func RandomUserEs(id int) *UserES {
	e := &UserES{
		ID:           id,
		Profession:   profession[utils.RangeRand(0, len(profession)-1)],
		Username:     utils.GetRandomString(5),
		Email:        utils.GetRandomString(7) + "@" + utils.GetRandomString(3) + ".com",
		TextInfo:     utils.GetRandomString(5),
		RegisterTime: int64(utils.RangeRand(1546272000, 1640966399)),
		Attribute: func() []string {
			max := utils.RangeRand(1, 5)
			s := make([]string, max)
			for i := 0; i < max; i++ {
				s[i] = attribute[utils.RangeRand(0, len(attribute)-1)]
			}
			s2 := make([]string, 0)
			_ = utils.ArrayUnique(&s2, s)
			return s2
		}(),
		Pets: func() []*Pet {
			max := utils.RangeRand(1, 5)
			sp := make([]*Pet, max)
			for i := 0; i < max; i++ {
				sp[i] = RandomPet()
			}
			return sp
		}(),
	}

	return e
}

func RandomPet() *Pet {
	return &Pet{
		PetName:   petNames[utils.RangeRand(0, len(petNames)-1)],
		PetSex:    2,
		PetAttack: utils.RangeRand(10, 10000),
		PetTag: func() []string {
			max := utils.RangeRand(1, 5)

			s := make([]string, max)
			for i := 0; i < max; i++ {
				s[i] = petTag[utils.RangeRand(0, len(petTag)-1)]
			}
			s2 := make([]string, 0)
			_ = utils.ArrayUnique(&s2, s)
			return s2
		}(),
	}
}

func TmpSearch()  {
	cli := utils.NewClient()
	cli.SetIndex("user")
}