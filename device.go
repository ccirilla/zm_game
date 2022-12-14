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
	CurIndex  int `json:"CurIndex"`
	FinishSum int `json:"FinishSum"`
	RoleLevel int `json:"RoleLevel"`
	IncGold   int `json:"IncGold"`
	IncMoney  int `json:"IncMoney"`
	IncEx     int `json:"IncEx"`
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

func CopyData(old *Device, new *Device) {
	*old = *new
	/*
		old.IncAllMoney = new.IncAllMoney
		old.IncAllGold = new.IncAllGold
		old.HotJob = new.HotJob
		old.Level1Step = new.Level1Step
		old.Level2Step = new.Level2Step
		*
	*/
}

func UpdateDeviceInfo(device *Device) *Device {
	if device.Did >= DeviceNum {
		return nil
	}
	m_device.Lock()
	devices[device.Did].Host = device.Host
	sys_device := devices[device.Did]
	m_device.Unlock()

	if device.Mark != sys_device.Mark {
		return &sys_device
	}
	m_device.Lock()
	CopyData(&devices[device.Did], device)
	m_device.Unlock()
	return nil
}

func GetHealthInfo() map[string]interface{} {
	var up_cnt = 0
	var down_cnt = 0
	var t = time.Now().Unix()
	farr := []map[string]int{}
	m_device.Lock()
	for _, val := range devices {
		if t - int64(val.LastActive) < int64(ActiveOverTime) {
			up_cnt++
		} else {
			down_cnt++
			tmp := map[string]int{
				"Did":      val.Did,
				"OverTime": (int(t) - val.LastActive) / 60,
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

func GetAllDevices(host string) []Device {
	ret := []Device{}
	m_device.Lock()
	for _, val := range devices {
		if val.Host == host {
			ret = append(ret, val)
		}
	}
	m_device.Unlock()
	return ret
}

func GetDidDevices(did int) Device {
	if did > DeviceNum {
		return Device{}
	}
	m_device.Lock()
	ret := devices[did]
	m_device.Unlock()
	return ret
}

func SetAllDevices(op string, param []string) {
	m_device.Lock()
	for i := 0; i < len(devices); i++ {
		if op == "task" {
			devices[i].Tasks = param
			devices[i].Level1Step.CurIndex = 0
			devices[i].Level1Step.SubIndex = 0
			devices[i].Mark = strconv.FormatInt(time.Now().Unix(), 10)
		} else if op == "hot" {
			devices[i].HotJob = param[0]
			devices[i].Mark = strconv.FormatInt(time.Now().Unix(), 10)
		}
	}
	m_device.Unlock()
}

func SetHostDevices(host string, op string, param []string) {
	m_device.Lock()
	for i := 0; i < len(devices); i++ {
		if devices[i].Host != host {
			continue
		}
		if op == "task" {
			devices[i].Tasks = param
			devices[i].Level1Step.CurIndex = 0
			devices[i].Level1Step.SubIndex = 0
			devices[i].Mark = strconv.FormatInt(time.Now().Unix(), 10)
		} else if op == "hot" {
			devices[i].HotJob = param[0]
			devices[i].Mark = strconv.FormatInt(time.Now().Unix(), 10)
		}
	}
	m_device.Unlock()
}

func SetDidDevices(did int, op string, param []string) {
	if did > DeviceNum {
		return
	}
	m_device.Lock()

	if op == "task" {
		devices[did].Tasks = param
		devices[did].Level1Step.CurIndex = 0
		devices[did].Level1Step.SubIndex = 0
		devices[did].Mark = strconv.FormatInt(time.Now().Unix(), 10)
	} else if op == "hot" {
		devices[did].HotJob = param[0]
		devices[did].Mark = strconv.FormatInt(time.Now().Unix(), 10)
	}
	m_device.Unlock()
}

func InitDevice(device *Device) {
	device.Tasks = []string{"??????", "??????", "??????", "Over"}
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
	//device.LastActive = 0
	//device.Host = "A"
}

func InitDevices(devices []Device) {
	for i := 0; i < DeviceNum; i++ {
		InitDevice(&devices[i])
		devices[i].Did = i
		devices[i].Host = "A"
		devices[i].LastActive = 0
	}
}

func InitALlDevice() {
	devices = make([]Device, DeviceNum)
	InitDevices(devices)
}

func SyncAllDevice()  {
	m_device.Lock()
	for i := 0; i < DeviceNum; i++ {
		InitDevice(&devices[i])
		devices[i].Tasks = []string{"??????", "??????", "??????", "Over"}
	}
	m_device.Unlock()
}