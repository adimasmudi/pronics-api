package helper

import (
	"errors"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/api/oauth2/v2"
)

func FormatValidationError(err error) []string {
	var errors []string

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}

func ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return token, err
	}

	return token, nil
}

var httpClient = &http.Client{}

func VerifyIdToken(idToken string) (*oauth2.Tokeninfo, error) {
    oauth2Service, err := oauth2.New(httpClient)
    tokenInfoCall := oauth2Service.Tokeninfo()
    tokenInfoCall.IdToken(idToken)
    tokenInfo, err := tokenInfoCall.Do()
    if err != nil {
        return nil, err
    }
    return tokenInfo, nil
}