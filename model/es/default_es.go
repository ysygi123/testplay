package es

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/olivere/elastic/v7"
	"reflect"
	"testplay/utils"
)

type QueryBase struct {
	Page int
	Size int
	Sort []*SortBase
}

type SortBase struct {
	SortField string
	Ascending bool
}

func GetDefaultEsJsonByte(s []*elastic.SearchHit) (returnByte []byte) {
	returnByte = make([]byte, 0)
	tmpByte := make([]byte, 0)
	if len(s) == 0 {
		return
	}
	for _, v := range s {
		tmpByte = append(tmpByte, v.Source...)
		tmpByte = append(tmpByte, ',')
	}
	tmpByte = tmpByte[:len(tmpByte)-1]
	returnByte = append(returnByte, '[')
	returnByte = append(returnByte, tmpByte...)
	returnByte = append(returnByte, ']')
	return
}

func DecodeDefaultEsModel(s []*elastic.SearchHit, target interface{}) (err error) {
	targetType := reflect.TypeOf(target)
	if targetType.Kind() != reflect.Ptr {
		err = errors.New("传入类型错误")
		return
	}

	targetValue := reflect.ValueOf(target)
	normalTarget := reflect.MakeSlice(targetValue.Elem().Type(), 0, 0)
	for _, v := range s {
		tmpX := reflect.New(targetValue.Elem().Type().Elem().Elem())
		fmt.Println(tmpX.Type())
		err = jsoniter.Unmarshal(v.Source, &tmpX)
		fmt.Println(utils.Data2json(tmpX), err)
		normalTarget = reflect.Append(normalTarget, tmpX)
	}
	//targetValue.Set(normalTarget)
	return
}
