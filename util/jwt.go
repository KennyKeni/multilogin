package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GetTokenExpiration(tokenString string) (time.Time, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return time.Time{}, err
	}

	// Get the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"]; ok {
			switch expValue := exp.(type) {
			case float64:
				expTime := time.Unix(int64(expValue), 0)
				return expTime, nil
			case int64:
				expTime := time.Unix(expValue, 0)
				return expTime, nil
			}
		}
	}

	return time.Time{}, fmt.Errorf("expiration claim not found or invalid")
}
