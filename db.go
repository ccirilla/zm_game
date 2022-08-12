package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

type RoleBaseInfoXc struct {
	//Id int
	Account    string `json:"account,omitempty"`
	Password   string
	RoleId     string
	Level      int
	RestEx     int
	Gold       int
	Money      int
	Atk        int
	Def        int
	CreateTime time.Time
}

func GetAccountInfo(did int) []map[string]string {
	data := make([]map[string]string, EveryDidNum)
	index := did * EveryDidNum
	m_role.Lock()
	for i := 0; i < EveryDidNum && index+i < AccountsNum; i++ {
		account := make(map[string]string)
		account["Account"] = roles[index+i].Account
		account["PassWord"] = roles[index+i].Password
		account["RoleId"] = roles[index+i].RoleId
		data[i] = account
	}
	m_role.Unlock()
	return data
}

var base_db *gorm.DB
var roles []RoleBaseInfoXc
var m_role sync.Mutex

func InitRoleInfo() {
	base_db.Limit(AccountsNum).Find(&roles)
}

func DbInit() {
	var err error
	var dsn string = "root:Qq950711.@(sh-cynosdbmysql-grp-jih4m75e.sql.tencentcdb.com:23017)/zm_game?charset=utf8mb4&parseTime=True&loc=Local"
	base_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		println(err)
		panic(err)
	}
	//db.AutoMigrate(&RawAccountUsing{})
	InitRoleInfo()
}
