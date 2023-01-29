package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var BillDb *gorm.DB

func init(){
	var err error
	dsn := "crdu:YwCo#I74m@tcp(10.96.140.69:3306)/billsystem?charset=utf8&parseTime=True&loc=Local"
	BillDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("mysql connect error %v", err)
	}

	if BillDb.Error != nil {
		fmt.Printf("database error %v", BillDb.Error)
	}
}