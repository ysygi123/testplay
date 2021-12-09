package utils

import (
	"fmt"
	"reflect"
)

// ArrayColumn array_column函数
func ArrayColumn(desk, input interface{}, columnKey, indexKey string) (err error) {
	if desk == nil || input == nil || (len(columnKey) == 0 && len(indexKey) == 0) {
		err = fmt.Errorf("错误")
		return
	}

	if len(columnKey) == 0 && len(indexKey) != 0{
		err = indexChange(desk, input,indexKey)
	}

	if len(columnKey) != 0 && len(indexKey) != 0 {
		err = indexColumnChange(desk, input, columnKey, indexKey)
	}

	if len(columnKey) != 0 && len(indexKey) == 0 {
		err = columnChange(desk, input, columnKey)
	}
	return
}

// commonValidate 公共验证
func commonValidate(deskType, inputType reflect.Type, deskValue, inputValue reflect.Value) (err error) {
	if deskType.Kind() != reflect.Ptr {
		err = fmt.Errorf("不是指针")
		return
	}

	if deskValue.Elem().IsNil() {
		err = fmt.Errorf("请初始化目标")
		return
	}

	if inputType.Kind() != reflect.Slice && inputType.Kind() != reflect.Array {
		err = fmt.Errorf("input类型错误 不是 array或者slice")
		return
	}

	if inputValue.Len() == 0 {
		err = fmt.Errorf("传入数组为长度为0")
		return
	}

	return
}

// commonIndexValidate 公共验证
func commonIndexValidate(deskType reflect.Type) (err error)  {
	if deskType.Elem().Kind() != reflect.Map {
		err = fmt.Errorf("不是map")
		return
	}

	return
}

// indexChange 只传index 返回是 map[Type]struct
func indexChange(desk, input interface{}, indexKey string) (err error) {
	deskType := reflect.TypeOf(desk)
	deskValue := reflect.ValueOf(desk)
	inputValue := reflect.ValueOf(input)
	inputType := reflect.TypeOf(input)
	if err = commonValidate(deskType, inputType, deskValue, inputValue); err != nil {
		return
	}

	if err = commonIndexValidate(deskType); err != nil {
		return
	}

	if deskType.Elem().Elem().String() != inputType.Elem().String() {
		err = fmt.Errorf("map值和数组值类型不对等")
		return
	}

	structFieldKey, ok := inputType.Elem().Elem().FieldByName(indexKey)
	if !ok {
		err = fmt.Errorf("无 key 字段")
		return
	}

	if deskType.Elem().Key().String() != structFieldKey.Type.Kind().String() {
		err = fmt.Errorf("key类型不对等")
		return
	}

	for i := 0; i < inputValue.Len(); i++ {
		deskValue.Elem().SetMapIndex(inputValue.Index(i).Elem().FieldByName(indexKey), inputValue.Index(i))
	}

	return
}

// indexColumnChange 传了key和column，那么返回必定是 map[Type]Type
func indexColumnChange(desk, input interface{}, columnKey, indexKey string) (err error) {
	deskType := reflect.TypeOf(desk)
	deskValue := reflect.ValueOf(desk)
	inputValue := reflect.ValueOf(input)
	inputType := reflect.TypeOf(input)
	if err = commonValidate(deskType, inputType, deskValue, inputValue); err != nil {
		return
	}

	if err = commonIndexValidate(deskType); err != nil {
		return
	}

	//如果input是个[]map或者[]*map
	if inputType.Elem().Kind() == reflect.Map{
		for i := 0; i < inputValue.Len(); i++ {
			findKey := false
			findCol := false
			key := reflect.Value{}
			value := reflect.Value{}
			for _, v := range inputValue.Index(i).MapKeys(){
				err = checkKeyAndColumn(indexKey, columnKey, v, inputValue.Index(i), &findKey, &findCol, deskType, &key, &value)
				if err != nil {
					return
				}
			}
			if !findCol {
				err = fmt.Errorf("查无column")
				return
			}
			if !findKey {
				err = fmt.Errorf("查无key")
				return
			}
			deskValue.Elem().SetMapIndex(key, value)
		}
	}else if inputType.Elem().Kind() == reflect.Ptr && inputType.Elem().Elem().Kind() == reflect.Map {
		for i := 0; i < inputValue.Len(); i++ {
			findKey := false
			findCol := false
			key := reflect.Value{}
			value := reflect.Value{}
			for _, v := range inputValue.Index(i).Elem().MapKeys(){
				//这边多一个elem 多一个指针
				err = checkKeyAndColumn(indexKey, columnKey, v, inputValue.Index(i).Elem(), &findKey, &findCol, deskType, &key, &value)
				if err != nil {
					return
				}
			}
			if !findCol {
				err = fmt.Errorf("查无column")
				return
			}
			if !findKey {
				err = fmt.Errorf("查无key")
				return
			}
			deskValue.Elem().SetMapIndex(key, value)

		}
	}else if inputType.Elem().Kind() == reflect.Ptr && inputType.Elem().Elem().Kind() == reflect.Struct {
		structFieldKey, ok := inputType.Elem().Elem().FieldByName(indexKey)
		if !ok {
			err = fmt.Errorf("无 key 字段")
			return
		}

		if deskType.Elem().Key().String() != structFieldKey.Type.Kind().String() {
			err = fmt.Errorf("key类型不对等")
			return
		}

		structFieldColumn, ok := inputType.Elem().Elem().FieldByName(columnKey)
		if !ok {
			err = fmt.Errorf("无 column 字段")
			return
		}

		if deskType.Elem().Elem().String() !=  structFieldColumn.Type.Kind().String(){
			err = fmt.Errorf("value类型不对等")
			return
		}

		for i := 0; i < inputValue.Len(); i++ {
			deskValue.Elem().SetMapIndex(inputValue.Index(i).Elem().FieldByName(indexKey), inputValue.Index(i).Elem().FieldByName(columnKey))
		}
	}else {
		err = fmt.Errorf("传入数组类型错误")
	}

	return
}

// checkKeyAndColumn 检验key和value
func checkKeyAndColumn(indexKey, columnKey string, curMapValue, indexValue reflect.Value, findKey, findCol *bool, deskType reflect.Type, key, value *reflect.Value) (err error) {
	if !(*findKey) {
		*findKey, *key, err = checkIsThisKey(indexKey, curMapValue, indexValue, deskType)
		if err != nil {
			return
		}
	}

	if !(*findCol) {
		*findCol, *value, err = checkIsThisKey(columnKey, curMapValue, indexValue, deskType)
		if err != nil {
			return
		}
	}

	return
}

// checkIsThisKey 判断 map情况下 要的key或者column是否能找得到，类型是否一样
func checkIsThisKey(key string, v, indexValue reflect.Value, deskType reflect.Type) (hasThis bool, returnValue reflect.Value, err error) {
	if v.String() != key {
		return
	}
	hasThis = true

	if indexValue.MapIndex(v).Kind().String() != deskType.Elem().Key().String() {
		err = fmt.Errorf("key 类型不对等")
	}else {
		returnValue = indexValue.MapIndex(v)
	}
	return
}

// columnChange 只传了column  那么返回必定是slice
func columnChange(desk, input interface{}, columnKey string) (err error) {
	deskType := reflect.TypeOf(desk)
	deskValue := reflect.ValueOf(desk)
	inputValue := reflect.ValueOf(input)
	inputType := reflect.TypeOf(input)
	err = commonValidate(deskType, inputType, deskValue, inputValue)
	if err != nil {
		return
	}

	columnStruct, ok := inputType.Elem().Elem().FieldByName(columnKey)
	if !ok {
		err = fmt.Errorf("column不出租")
		return
	}

	if columnStruct.Type.Kind().String() != deskType.Elem().Elem().String() {
		err = fmt.Errorf("值类型不对")
		return
	}
	dv := reflect.MakeSlice(deskValue.Elem().Type(), 0, 0)
	for i := 0; i < inputValue.Len(); i++ {
		dv = reflect.Append(dv, inputValue.Index(i).Elem().FieldByName(columnKey))
	}
	deskValue.Elem().Set(dv)
	return
}

// ArrayUnique 数组去重
func ArrayUnique(desk, input interface{}) (err error) {
	deskType := reflect.TypeOf(desk)
	deskValue := reflect.ValueOf(desk)
	inputValue := reflect.ValueOf(input)
	inputType := reflect.TypeOf(input)
	if err = commonValidate(deskType, inputType, deskValue, inputValue); err != nil {
		return
	}
	if deskType.Elem() != inputType {
		err = fmt.Errorf("类型错误")
		return
	}
	tmpEmptyStructType := reflect.TypeOf(struct {}{})
	tmpEmptyStructValue := reflect.ValueOf(struct {}{})
	tmpMap := reflect.MakeMap(reflect.MapOf(inputType.Elem(), tmpEmptyStructType))
	dv := reflect.MakeSlice(deskValue.Elem().Type(), 0, 0)
	for i := 0; i < inputValue.Len(); i++  {
		x := tmpMap.MapIndex(inputValue.Index(i))
		if !x.IsValid() {
			tmpMap.SetMapIndex(inputValue.Index(i), tmpEmptyStructValue)
			dv = reflect.Append(dv, inputValue.Index(i))
		}
	}
	deskValue.Elem().Set(dv)
	return
}