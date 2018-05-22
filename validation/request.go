package validation

import (
	"github.com/devicehive/devicehive-frontend-go/pg"
	"github.com/devicehive/devicehive-frontend-go/dh-jwt"
)

func checkAction(claims *dhjwt.JwtCustomClaims, action Action)  {

}

func HasPermission(device *pg.Device, claims *dhjwt.JwtCustomClaims, action Action) bool  {
	return true
}
