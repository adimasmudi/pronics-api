package helper

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(ID string) (string, error) {
	claim := jwt.MapClaims{}
	
	claim["ID"] = ID
	claim["generateAt"] = time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}