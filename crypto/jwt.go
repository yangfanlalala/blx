package crypto

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

func JwtToken(claims jwt.Claims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func JwtParse(claims jwt.Claims, signed, secret string) error {
	token, err := jwt.ParseWithClaims(signed, claims, func(tkn *jwt.Token) (interface{}, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的加密算法: %v", tkn.Header["alg"])
		}
		jwtSecret := secret
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("无效的Token")
	}
	return nil
}
