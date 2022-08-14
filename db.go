package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
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

type ErrorInfo struct{
	Func string `json:"Func"`
	Scene string `json:"Scene"`
}

type TaskMessage struct {
	RetType string `json:"RetType"`
	ErrInfo ErrorInfo  `json:"ErrInfo"`
	Task Device `json:"Task"`
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

func RecordTaskLog(data *TaskMessage){
	s:= fmt.Sprintf("Type: %s Did: %d Task %s SubIndex: %d State: %s UseTime: %d  HotJob: %s " +
						"CurIndex: %d FinishSum: %d IncMoney: %d IncGold:%d IncEx: %d",
					data.RetType, data.Task.Did, data.Task.Tasks[data.Task.Level1Step.CurIndex], data.Task.Level1Step.SubIndex,
					data.Task.State, data.Task.Level1Step.UsedTime, data.Task.HotJob,
					data.Task.Level2Step.CurIndex, data.Task.Level2Step.FinishSum, data.Task.Level2Step.IncMoney,
					data.Task.Level2Step.IncGold, data.Task.Level2Step.IncEx)
	if data.RetType == "Done" || data.RetType == "HotDone"{
		SLogger.Println(s)
		//记录数据库
	}else {
		FLogger.Println(s)
	}
}

var base_db *gorm.DB
var roles []RoleBaseInfoXc
var m_role sync.Mutex

var (
	SLogger *log.Logger
	FLogger *log.Logger
	ELogger *log.Logger
)

func initLog() {
	sfile, err := os.OpenFile("./log/success.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	ffile, err := os.OpenFile("./log/false.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	efile, err := os.OpenFile("./log/event.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	SLogger = log.New(sfile, "", log.Ldate|log.Ltime)
	FLogger = log.New(ffile, "", log.Ldate|log.Ltime)
	ELogger = log.New(efile, "", log.Ldate|log.Ltime)
}

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
	initLog()
	//db.AutoMigrate(&RawAccountUsing{})
	InitRoleInfo()

}
