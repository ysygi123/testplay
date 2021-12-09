package service

import (
	"fmt"
	"reflect"
)

type FUCK int

func T1() {
	const Zero FUCK = 0
	type cat struct {
		Name string
		Type int `json:"type" id:"100"`
	}

	cgitta := &cat{
		Name: "使命召唤",
		Type: 3,
	}

	typeOfCat := reflect.TypeOf(cgitta)
	fmt.Println(typeOfCat.Elem().Name(), typeOfCat.Kind())
	//for i := 0; i < typeOfCat.NumField(); i++ {
	//	filedType := typeOfCat.Field(i)
	//	fmt.Printf("name: %v  tag: '%v'\n", filedType.Name, filedType.Tag)
	//}
	//
	//if catType, ok := typeOfCat.FieldByName("Type"); ok {
	//	fmt.Println(catType.Tag.Get("json"), catType.Tag.Get("id"))
	//}
}
