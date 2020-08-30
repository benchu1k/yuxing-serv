package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb = getRedis()
var requestNum uint64 = 0

// DeviceInfo 设备信息的结构体
type DeviceInfo struct {
	Imei      string  `json:"imei"`
	DeviceID  string  `json:"deviceId"`
	ProductID string  `json:"productId"`
	AppData   string  `json:"appData"`
	Timestamp float64 `json:"timestamp"`
	TenantID  string  `json:"tenantId"`
}

func main() {

	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()

	app.POST("/deviceStatus", deviceStatus)

	app.Run(":8989")

}

func deviceStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
	requestNum++
	fmt.Printf("%d Connected! \n", requestNum)

	buf := make([]byte, 1024)
	num, _ := c.Request.Body.Read(buf)

	var f interface{}
	err := json.Unmarshal(buf[0:num], &f)
	if err != nil {
		fmt.Println(err)
	}

	m := f.(map[string]interface{})
	ds := &DeviceInfo{}

	for k, v := range m {
		switch k {
		case "IMEI":
			ds.Imei = v.(string)
		case "deviceId":
			ds.DeviceID = v.(string)
		case "productId":
			ds.ProductID = v.(string)
		case "payload":
			vv := v.(map[string]interface{})
			for k2, v2 := range vv {
				if k2 == "APPdata" {
					ds.AppData = v2.(string)
				}
			}
		case "timestamp":
			ds.Timestamp = v.(float64)
		case "tenantId":
			ds.TenantID = v.(string)
		}
	}

	data, _ := json.Marshal(ds)
	dataStr := string(data)

	setValue(ds.Imei, dataStr)
}

func setValue(key string, value interface{}) {
	err := rdb.Set(ctx, key, value, 0).Err()

	if err != nil {
		fmt.Println("连接失败")
		panic(err)
	}
}

func getRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       4,
	})

	defer func() {
	}()

	_, err := rdb.Ping(ctx).Result()

	if err != nil {
		fmt.Println(" Redis Connect failed !")
		panic(err)
	}

	return rdb
}
