package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
	"sync"
)

var gm sync.Mutex

func main() {

	router := gin.Default()
	router.GET("/GetAccount", GetAccount)
	router.GET("/ReportInfo", ReportInfo)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func GetAccount(c *gin.Context) {
	did_str := c.Query("Did")
	did_int, _ := strconv.Atoi(did_str)
	var data = RawAccountCb{}

	gm.Lock()
	FindAccount(db, &data)
	UpdateAccount(db, data.Id, did_int)
	gm.Unlock()

	c.JSON(200, gin.H{
		"Account":  data.Account,
		"Password": data.Password,
	})
}

func ReportInfo(c *gin.Context) {
	var val = ReportData{}
	param := c.Query("data")
	//println(param)
	json.Unmarshal([]byte(param), &val)

	UpdateTaskResult(&val)

	c.JSON(200, gin.H{
		"message": "OK",
	})
}
