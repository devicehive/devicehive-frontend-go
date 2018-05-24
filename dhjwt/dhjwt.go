package dhjwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtCustomClaims struct {
	Payload struct {
		Actions    []int8   `json:"a"`
		ExpiresAt  int64    `json:"e"`
		TypeUser   int      `json:"t"`
		User       int64    `json:"u"`
		Networks   []string `json:"n"`
		DeviceType []string `json:"dt"`
	} `json:"payload"`
}

func (c JwtCustomClaims) Valid() error {
	vErr := new(jwt.ValidationError)
	now := jwt.TimeFunc().Unix()

	// The claims below are optional, by default, so if they are set to the
	// default value in Go, let's not fail the verification for them.
	if c.VerifyExpiresAt(now, false) == false {
		delta := time.Unix(now, 0).Sub(time.Unix(c.Payload.ExpiresAt, 0))
		vErr.Inner = fmt.Errorf("token is expired by %v", delta)
		vErr.Errors |= jwt.ValidationErrorExpired
	}

	//if c.VerifyIssuedAt(now, false) == false {
	//	vErr.Inner = fmt.Errorf("Token used before issued")
	//	vErr.Errors |= jwt.ValidationErrorIssuedAt
	//}
	//
	//if c.VerifyNotBefore(now, false) == false {
	//	vErr.Inner = fmt.Errorf("token is not valid yet")
	//	vErr.Errors |= jwt.ValidationErrorNotValidYet
	//}

	if vErr.Errors == 0 {
		return nil
	}

	return vErr
}

// Compares the aud claim against cmp.
// If required is false, this method will return true if the value matches or is unset
//func (c *JwtCustomClaims) VerifyAudience(cmp string, req bool) bool {
//	return verifyAud(c.Audience, cmp, req)
//}

// Compares the exp claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *JwtCustomClaims) VerifyExpiresAt(cmp int64, req bool) bool {
	return verifyExp(c.Payload.ExpiresAt, cmp, req)
}

// Compares the iat claim against cmp.
// If required is false, this method will return true if the value matches or is unset
//func (c *JwtCustomClaims) VerifyIssuedAt(cmp int64, req bool) bool {
//	return verifyIat(c.IssuedAt, cmp, req)
//}

// Compares the iss claim against cmp.
// If required is false, this method will return true if the value matches or is unset
//func (c *JwtCustomClaims) VerifyIssuer(cmp string, req bool) bool {
//	return verifyIss(c.Issuer, cmp, req)
//}

// Compares the nbf claim against cmp.
// If required is false, this method will return true if the value matches or is unset
//func (c *JwtCustomClaims) VerifyNotBefore(cmp int64, req bool) bool {
//	return verifyNbf(c.NotBefore, cmp, req)
//}

// ----- helpers

//func verifyAud(aud string, cmp string, required bool) bool {
//	if aud == "" {
//		return !required
//	}
//	if subtle.ConstantTimeCompare([]byte(aud), []byte(cmp)) != 0 {
//		return true
//	} else {
//		return false
//	}
//}

func verifyExp(exp int64, now int64, required bool) bool {
	if exp == 0 {
		return !required
	}
	return now <= exp
}

//func verifyIat(iat int64, now int64, required bool) bool {
//	if iat == 0 {
//		return !required
//	}
//	return now >= iat
//}
//
//func verifyIss(iss string, cmp string, required bool) bool {
//	if iss == "" {
//		return !required
//	}
//	if subtle.ConstantTimeCompare([]byte(iss), []byte(cmp)) != 0 {
//		return true
//	} else {
//		return false
//	}
//}
//
//func verifyNbf(nbf int64, now int64, required bool) bool {
//	if nbf == 0 {
//		return !required
//	}
//	return now >= nbf
//}
