package jwt

import (
	"fmt"
	"metaLand/app/api/internal/config"
	"strconv"
	"time"

	JWT "github.com/dgrijalva/jwt-go"
)

// Sign jwt token for current uin.
func Sign(uin uint64, c config.Config) (token string) {
	jwt := JWT.New(JWT.SigningMethodHS256)
	claims := make(JWT.MapClaims)
	claims["exp"] = time.Now().Add(time.Duration(c.JWT.Expired) * time.Second).Unix()
	claims["iat"] = time.Now().Unix()
	claims["comer_uin"] = fmt.Sprintf("%d", uin)
	jwt.Claims = claims
	token, _ = jwt.SignedString([]byte(c.JWT.Secret))
	return
}

// Verify jwt token if success then return the uin.
func Verify(token string, c config.Config) (uin uint64, err error) {
	auth, err := JWT.Parse(token, func(t *JWT.Token) (interface{}, error) {
		return []byte(c.JWT.Secret), nil
	})
	if err != nil {
		return
	}
	claims, _ := auth.Claims.(JWT.MapClaims)
	uinStr, _ := claims["comer_uin"].(string)
	uin, err = strconv.ParseUint(uinStr, 10, 64)
	return
}
