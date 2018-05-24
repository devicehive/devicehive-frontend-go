package validation

import (
	"github.com/devicehive/devicehive-frontend-go/dhjwt"
	"github.com/devicehive/devicehive-frontend-go/pg"
)

func inArrayInt8(arr []int8, el int8) bool {
	for _, v := range arr {
		if v == el || v == 0 {
			return true
		}
	}
	return false
}

func inArrayString(arr []string, el string) bool {
	for _, v := range arr {
		if v == el || v == "*" {
			return true
		}
	}
	return false
}

func HasPermission(device *pg.Device, claims *dhjwt.JwtCustomClaims, action Action) bool {
	return inArrayInt8(claims.Payload.Actions, int8(action)) &&
		inArrayString(claims.Payload.DeviceType, string(device.DeviceTypeId)) &&
		inArrayString(claims.Payload.Networks, string(device.NetworkId))
}
