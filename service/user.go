package service

import (
	"strconv"
	"testplay/model/es"
	"testplay/utils"
)

func AddUser(id int, username, email, textInfo string) error {
	user := &es.UserES{
		ID:       id,
		Username: username,
		Email:    email,
		TextInfo: textInfo,
	}

	ess := utils.NewClient()
	_, err := ess.SetIndex("user").Upsert(strconv.Itoa(user.ID), user)
	if err != nil {
		return err
	}
	return nil
}

func SearchUser() {

}

func DelUser() {
	ess := utils.NewClient()
	ess.SetIndex("user").EsClient.Delete()
}
