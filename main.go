package main

import (
	"database/sql"
	"fmt"
	"github.com/devicehive/devicehive-frontend-go/dhjwt"
	"github.com/devicehive/devicehive-frontend-go/dhkafka"
	"github.com/devicehive/devicehive-frontend-go/notifications"
	"github.com/devicehive/devicehive-frontend-go/pg"
	"github.com/devicehive/devicehive-frontend-go/validation"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"github.com/patrickmn/go-cache"
	"time"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"github.com/hazelcast/hazelcast-go-client"
	"runtime"
)

const (
	JWT_SECRET  = "magic"
	DB_CONN_STR = "postgres://postgres:12345@localhost/devicehive?sslmode=disable"
)

type ResponceKafka struct{
	Cid string `json:"cId"`
	Err int `json:"err"`
	Message string `json:"message"`
}

type ResponceNotification struct {
	Cid string `json:"cId"`
	Message string `json:"message"`
	Err int `json:"err"`
	Body struct{
		DeviceNotification struct{
			timestamp int64 `json:"timestamp"`
		} `json:"deviceNotification"`
	} `json:"b"`
}


func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	e := echo.New()
	db := pg.New(DB_CONN_STR)
	kafka := dhkafka.New([]string{"localhost:9092"}, e.Logger)
	defer db.Close()
	defer kafka.Close()
	pConsumer := kafka.GetMessages()
	c := cache.New(40*time.Second, 1*time.Minute)
	configHazel := hazelcast.NewHazelcastConfig()
	configHazel.ClientNetworkConfig().AddAddress("127.0.0.1:5701")

	clientHazel, errHazel := hazelcast.NewHazelcastClientWithConfig(configHazel)
	if errHazel != nil {
		log.Fatal(errHazel)
	}
	defer clientHazel.Shutdown()
		go func() {
			for {
				select {
				case msg := <-pConsumer.Messages():
					go func() {
						res := ResponceKafka{}
						json.Unmarshal(msg.Value, &res)
						if ch, ok:=c.Get(res.Cid); ok{
							ch.(chan []byte) <- msg.Value
							c.Delete(res.Cid)
						}
					}()
				default:
					continue
				}
			}

		}()



	config := middleware.JWTConfig{
		Claims:     &dhjwt.JwtCustomClaims{},
		SigningKey: []byte(JWT_SECRET),
	}
	e.Use(middleware.JWTWithConfig(config))

	e.GET("/", func(ctx echo.Context) error {
		user := ctx.Get("user").(*jwt.Token)
		claims := user.Claims.(*dhjwt.JwtCustomClaims)
		//fmt.Println(claims)
		return ctx.JSON(http.StatusOK, claims)
	})

	e.POST("/device/:deviceId/notification", func(ctx echo.Context) error {
		var (
			device pg.Device
			err error
		)
		user := ctx.Get("user").(*jwt.Token)
		claims := user.Claims.(*dhjwt.JwtCustomClaims)
		deviceId := ctx.Param("deviceId")
		m, _:=clientHazel.GetMap("Devices_go")
		vH, errH := m.Get(deviceId);
		if  errH == nil && vH != nil {
			device = vH.(pg.Device)
		}else {
			device, err = db.GetDeviceById(deviceId)
			m.Set(deviceId, device)
			switch {
			case err == sql.ErrNoRows:
				e.Logger.Errorf("Device with such deviceId = %s not found", deviceId)
				return ctx.JSON(http.StatusNotFound,
					ErrorMessage{
						Error:   http.StatusNotFound,
						Message: fmt.Sprintf("Device with such deviceId = %s not found", deviceId),
					})
			case err != nil:
				e.Logger.Errorf("Pg error: %s ", err)
				return ctx.JSON(http.StatusBadRequest, ErrorMessage{
					Error:   http.StatusBadRequest,
					Message: "Internal Error",
				})
			}
		}
		if !validation.HasPermission(&device, claims, validation.GET_DEVICE_NOTIFICATION) {
			return ctx.JSON(http.StatusBadRequest, ErrorMessage{
				Error:   http.StatusBadRequest,
				Message: "You have not permissions for this action",
			})
		}
		body := new(notifications.DeviceNotificationWrapper)
		ctx.Bind(body)
		notif := notifications.New(&device, body)
		ch:= make(chan []byte)
		c.Set(kafka.Send(notif), ch, cache.DefaultExpiration)
		res := <-ch
		data := ResponceNotification{}
		json.Unmarshal(res, &data)
		if data.Err != 0 {
			return ctx.JSON(http.StatusBadRequest, struct {
				Err int `json:"err"`
				Message string `json:"message"`
			}{
				Err:data.Err,
				Message:data.Message,
			})
		}else{
			return ctx.JSON(http.StatusOK, struct {
				Id string `json:"id"`
				Timestamp int64 `json:"timestamp"`
			}{
				Id: data.Cid,
				Timestamp: data.Body.DeviceNotification.timestamp,
			})
		}
	})

	e.Logger.Fatal(e.Start(":1323"))
}
