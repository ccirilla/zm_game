package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	DbInit()
	InitALlDevice()
}



func main() {

	router := gin.Default()
	router.GET("/GetAccount", GetAccount)
	router.GET("/ReportHeart", ReportHeart)
	router.GET("/GetHealth", GetHealth)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run(":8001") // 监听并在 0.0.0.0:8080 上启动服务
}

func GetAccount(c *gin.Context) {
	did_str := c.Query("Did")
	did_int, _ := strconv.Atoi(did_str)
	data := GetAccountInfo(did_int)

	c.JSON(200, gin.H{
		"Status": "OK",
		"Data":   data,
		"Count":  len(data),
	})
}

func ReportHeart(c *gin.Context) {
	param := c.Query("data")
	var val = Device{}
	json.Unmarshal([]byte(param), &val)
	val.LastActive = int(time.Now().Unix())
	data := UpdateDeviceInfo(&val)
	if data == nil {
		c.JSON(200, gin.H{
			"Status": "OK",
			"Data":   data,
		})
		return
	}
	c.JSON(200, gin.H{
		"Status": "Update",
		"Data":   data,
	})
}

func GetHealth(c *gin.Context) {
	data := GetHealthInfo()

	c.JSON(200, gin.H{
		"Status": "OK",
		"Data":   data,
	})
}

func GetDevice(c *gin.Context) {
	param := c.Query("param")
	params := strings.Split(param, "-")
	var data interface{}

	if params[0] == "getall"{
		data = GetAllDevices(params[1])
	}else{
		did, _ := strconv.Atoi(params[1])
		data = GetDidDevices(did)
	}
	c.JSON(200, gin.H{
		"Status": "OK",
		"Data":   data,
	})
}