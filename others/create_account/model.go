package main

import (
	"gorm.io/driver/mysql"
	"time"

	//"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RawAccountCb struct {
	Id       int    `json:"id,omitempty"`
	Account  string `json:"account,omitempty"`
	Password string `json:"password,omitempty"`
	Status   int    `gorm:"default:0"`
	Did      int
}

type AccountUnit struct {
	Account  string
	Password string
}

type CreateRoleSuccessRecordSecond struct {
	Account    string `json:"Account,omitempty"`
	StartTime  string
	UseTime    float32
	Did        int
	KillSum    int
	Up15Time   float32
	AttackTime float32
	Up20Time   float32
	Up24Time   float32
	Up28Time   float32
	CreateTime time.Time
}

type CreateRoleFalseRecordSecond struct {
	Account    string `json:"account,omitempty"`
	StartTime  string
	UseTime    float32
	Did        int
	KillSum    int
	ErrorType  int
	MainTask   string
	SmallTask  string
	Scene      string
	Func       string
	Level      int
	Atk        int
	Def        int
	CreateTime time.Time
}

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

type RoleBaseInfoCb struct {
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

type CommonData struct {
	Account   string  `json:"Account,omitempty"`
	Password  string  `json:"Password,omitempty"`
	Did       int     `json:"Did,omitempty"`
	StartTime string  `json:"StartTime,omitempty"`
	UseTime   float32 `json:"UseTime,omitempty"`
	KillSum   int     `json:"KillSum,omitempty"`
	Id        string  `json:"Id,omitempty"`
	Level     int     `json:"Level,omitempty"`
	RestEx    int     `json:"RestEx,omitempty"`
	Gold      int     `json:"Gold,omitempty"`
	Money     int     `json:"Money,omitempty"`
	Atk       int     `json:"Atk,omitempty"`
	Def       int     `json:"Def,omitempty"`
}

type TaskSuccess struct {
	Up15   float32 `json:"Up15,omitempty"`
	Attack float32 `json:"Attack,omitempty"`
	Up20   float32 `json:"Up20,omitempty"`
	Up24   float32 `json:"Up24,omitempty"`
	Up28   float32 `json:"Up28,omitempty"`
}

type TaskFalse struct {
	Type      int    `json:"Type,omitempty"`
	MainTask  string `json:"MainTask,omitempty"`
	SmallTask string `json:"SmallTask,omitempty"`
	Func      string `json:"Func,omitempty"`
	Scene     string `json:"Scene,omitempty"`
}

type ReportData struct {
	Common  CommonData  `json:"Common,omitempty"`
	Success TaskSuccess `json:"Success,omitempty"`
	False   TaskFalse   `json:"False,omitempty"`
	Did     int         `json:"Did,omitempty"`
	Ret     int         `json:"Ret,omitempty"`
}

func FindAccount(db *gorm.DB, account *RawAccountCb) {
	db.Where(&RawAccountCb{Status: 1}).First(account)
	//println(account.Account)
}

func UpdateAccount(db *gorm.DB, id, did int) {
	var account = RawAccountCb{Id: id}
	db.Model(&account).Updates(map[string]interface{}{"status": 0, "did": did})
}

func UpdateTaskResult(data *ReportData) {
	if data.Ret == 0 {
		var success_data = CreateRoleSuccessRecordSecond{}
		success_data.Account = data.Common.Account
		success_data.StartTime = data.Common.StartTime
		success_data.UseTime = data.Common.UseTime
		success_data.Did = data.Common.Did
		success_data.KillSum = data.Common.KillSum
		success_data.Up15Time = data.Success.Up15
		success_data.AttackTime = data.Success.Attack
		success_data.Up20Time = data.Success.Up20
		success_data.Up24Time = data.Success.Up24
		success_data.Up28Time = data.Success.Up28
		success_data.CreateTime = time.Now()
		db.Create(&success_data)

		//var role_data = RoleBaseInfoXc{
		var role_data = RoleBaseInfoCb{
			Account:    data.Common.Account,
			Password:   data.Common.Password,
			RoleId:     data.Common.Id,
			Level:      data.Common.Level,
			RestEx:     data.Common.RestEx,
			Gold:       data.Common.Gold,
			Money:      data.Common.Money,
			Atk:        data.Common.Atk,
			Def:        data.Common.Def,
			CreateTime: time.Now(),
		}
		db.Create(&role_data)
	} else {
		var false_data = CreateRoleFalseRecordSecond{
			Account:    data.Common.Account,
			StartTime:  data.Common.StartTime,
			UseTime:    data.Common.UseTime,
			Did:        data.Common.Did,
			KillSum:    data.Common.KillSum,
			ErrorType:  data.False.Type,
			MainTask:   data.False.MainTask,
			SmallTask:  data.False.SmallTask,
			Scene:      data.False.Scene,
			Func:       data.False.Func,
			Level:      data.Common.Level,
			Atk:        data.Common.Atk,
			Def:        data.Common.Def,
			CreateTime: time.Now(),
		}
		db.Create(&false_data)
	}
}

var db *gorm.DB

func init() {
	var err error
	var dsn string = "root:Qq950711.@(sh-cynosdbmysql-grp-jih4m75e.sql.tencentcdb.com:23017)/zm_game?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		println(err)
		panic(err)
	}
	//db.AutoMigrate(&RawAccountUsing{})
}
