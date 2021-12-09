package utils

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
)

var Engine *xorm.Engine

func init() {
	var err error
	Engine, err = xorm.NewEngine("mysql", "root:1qazxsw2@tcp(127.0.0.1:3306)/mytest1?charset=utf8mb4")
	if err != nil {
		panic(err)
	}
	Engine.SetMaxIdleConns(100)
	Engine.SetMaxIdleConns(500)

}

