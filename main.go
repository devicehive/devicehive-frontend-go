package main

import (
	"database/sql"
	"fmt"
	"github.com/devicehive/devicehive-frontend-go/dh-jwt"
	"github.com/devicehive/devicehive-frontend-go/pg"
	"github.com/devicehive/devicehive-frontend-go/validation"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

const JWT_SECRET = "magic"
const DB_CONN_STR = "postgres://postgres:mysecretpassword@localhost/postgres?sslmode=disable"



func main() {
	e := echo.New()
	db := pg.New(DB_CONN_STR)
	defer db.Close()

	config := middleware.JWTConfig{
		Claims:     &dhjwt.JwtCustomClaims{},
		SigningKey: []byte(JWT_SECRET),
	}
	e.Use(middleware.JWTWithConfig(config))

	e.GET("/", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*dhjwt.JwtCustomClaims)
		fmt.Println(claims)
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
		validation.HasPermission(device, claims, validation.GET_DEVICE_NOTIFICATION)
		fmt.Println(claims)
		return c.String(http.StatusOK, deviceId)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
