package authentication

import (
	"blogger/config"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// GenerateToken returns a signed token with permissions.
func GenerateToken(userID uint64) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	permissions["user_id"] = userID

	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)

	return userToken.SignedString(config.SECRET_KEY)
}

// TokenValidation verifies if the given token is valid.
func TokenValidation(r *http.Request) error {
	userToken := tokenExtraction(r)
	t, err := jwt.Parse(userToken, verificationKey)
	if err != nil {
		return err
	}

	// fmt.Println(t)

	if _, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return nil
	}

	return errors.New("invalid token")
}

func tokenExtraction(r *http.Request) string {
	userToken := r.Header.Get("Authorization")
	tmp := strings.Split(userToken, " ")

	if len(tmp) == 2 {
		userToken = tmp[1]
	}

	return userToken
}

func verificationKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return config.SECRET_KEY, nil
}

func UserIDExtraction(r *http.Request) (uint64, error) {
	userToken := tokenExtraction(r)
	t, err := jwt.Parse(userToken, verificationKey)
	if err != nil {
		return 0, err
	}

	if permissions, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		userID, err := strconv.ParseUint(fmt.Sprintf("%.0f", permissions["user_id"]), 10, 64)
		if err != nil {
			return 0, err
		}

		return userID, nil
	}

	return 0, errors.New("invalid token")
}
