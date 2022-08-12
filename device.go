package main

import (
	"strconv"
	"sync"
	"time"
)

var m_device sync.Mutex
var devices []Device

type Level1Step struct {
	CurIndex int `json:"CurIndex"`
	SubIndex int `json:"SubIndex"`
	UsedTime int `json:"UsedTime"`
}

type Level2Step struct {
	CurIndex   int `json:"CurIndex"`
	FinishSum int    `json:"FinishSum"`
	RoleLevel int    `json:"RoleLevel"`
	IncGold   int    `json:"IncGold"`
	IncMoney  int    `json:"IncMoney"`
	IncEx     int    `json:"IncEx"`
}

type Device struct {
	Did         int        `json:"Did"`
	Tasks       []string   `json:"Tasks"`
	Mark        string     `json:"Mark"`
	State       string     `json:"State"`
	HotJob      string     `json:"HotJob"`
	Level1Step  Level1Step `json:"Level1Step"`
	Level2Step  Level2Step `json:"Level2Step"`
	IncAllGold  int        `json:"IncAllGold"`
	IncAllMoney int        `json:"IncAllMoney"`
	Host        string     `json:"Host"`
	LastActive  int        `json:"LastActive"`
}

func CopyData(old *Device, new *Device){
	old.IncAllMoney = new.IncAllMoney
	old.IncAllGold = new.IncAllGold
	old.HotJob = new.HotJob
	old.Level1Step = new.Level1Step
	old.Level2Step = new.Level2Step
}

func UpdateDeviceInfo(device *Device) *Device {
	if device.Did >= DeviceNum{
		return nil
	}
	m_device.Lock()
	sys_device := devices[device.Did]
	m_device.Unlock()

	if device.Mark != sys_device.Mark{
		return &sys_device
	}
	m_device.Lock()
	CopyData(&devices[device.Did], device)
	m_device.Unlock()
	return nil
}

func GetHealthInfo() map[string]interface{}{
	var up_cnt = 0
	var down_cnt = 0
	var t = time.Now().Unix()
	farr := []map[string]int{}
	m_device.Lock()
	for _, val := range devices{
		ELogger.Printf("%+v", val)
		ELogger.Println(t - int64(val.LastActive))
		if val.LastActive == 0 || t - int64(val.LastActive) < int64(ActiveOverTime){
			up_cnt++;
		}else{
			down_cnt++;
			tmp := map[string]int{
				"Did" : val.Did,
				"OverTime" : (int(t) - val.LastActive) / 60,
			}
			farr = append(farr, tmp)
		}
	}
	m_device.Unlock()
	ret := make(map[string]interface{})
	ret["UP"] = up_cnt
	ret["Down"] = down_cnt
	ret["DownDids"] = farr
	return ret
}

func InitDevice(device *Device) {
	device.Tasks = []string{"挂机1", "挂机2", "金砖", "Over"}
	device.Mark = strconv.FormatInt(time.Now().Unix(), 10)
	device.HotJob = "NULL"
	device.IncAllGold = 0
	device.IncAllMoney = 0
	//level1
	device.Level1Step.CurIndex = 0
	device.Level1Step.SubIndex = 0
	device.Level1Step.UsedTime = 0
	//level2
	device.Level2Step.CurIndex = 0
	device.Level2Step.FinishSum = 0
	device.Level2Step.IncEx = 0
	device.Level2Step.IncGold = 0
	device.Level2Step.IncMoney = 0
	device.Level2Step.RoleLevel = 0
	device.LastActive = 0
}

func InitDevices(devices []Device) {
	for i := 0; i < DeviceNum; i++{
		InitDevice(&devices[i])
		devices[i].Did = i
	}
}

func InitALlDevice() {
	devices = make([]Device, DeviceNum+1)
	InitDevices(devices)
}
