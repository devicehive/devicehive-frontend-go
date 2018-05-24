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
)

const (
	JWT_SECRET  = "magic"
	DB_CONN_STR = "postgres://postgres:mysecretpassword@localhost/postgres?sslmode=disable"
)

func main() {
	e := echo.New()
	db := pg.New(DB_CONN_STR)
	kafka := dhkafka.New([]string{"localhost:9092"}, e.Logger)
	defer db.Close()
	defer kafka.Close()
	pConsumer := kafka.GetMessages()

	go func() {
		for {
			select {
			case msg := <-pConsumer.Messages():
				fmt.Println(msg)
			}
		}

	}()

	config := middleware.JWTConfig{
		Claims:     &dhjwt.JwtCustomClaims{},
		SigningKey: []byte(JWT_SECRET),
	}
	e.Use(middleware.JWTWithConfig(config))

	e.GET("/", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*dhjwt.JwtCustomClaims)
		//fmt.Println(claims)
		return c.JSON(http.StatusOK, claims)
	})

	e.POST("/device/:deviceId/notification", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*dhjwt.JwtCustomClaims)
		deviceId := c.Param("deviceId")
		device, err := db.GetDeviceById(deviceId)
		switch {
		case err == sql.ErrNoRows:
			e.Logger.Errorf("Device with such deviceId = %s not found", deviceId)
			return c.JSON(http.StatusNotFound,
				ErrorMessage{
					Error:   http.StatusNotFound,
					Message: fmt.Sprintf("Device with such deviceId = %s not found", deviceId),
				})
		case err != nil:
			e.Logger.Errorf("Pg error: %s ", err)
			return c.JSON(http.StatusBadRequest, ErrorMessage{
				Error:   http.StatusBadRequest,
				Message: "Internal Error",
			})
		}
		if !validation.HasPermission(device, claims, validation.GET_DEVICE_NOTIFICATION) {
			return c.JSON(http.StatusBadRequest, ErrorMessage{
				Error:   http.StatusBadRequest,
				Message: "You have not permissions for this action",
			})
		}
		body := new(notifications.DeviceNotificationWrapper)
		c.Bind(body)
		notif := notifications.New(device, body)
		kafka.Send(notif)

		return c.JSON(http.StatusOK, notif)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
