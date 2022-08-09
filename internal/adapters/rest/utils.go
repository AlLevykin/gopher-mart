package rest

import (
	"compress/gzip"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"net/http"
	"time"
)

var (
	ErrJWTValidation = errors.New("jwt validation error")
)

var JwtSecretKey = []byte("my_secret_key")

type Claims struct {
	Login string
	jwt.RegisteredClaims
}

func NewToken(userClaims *Claims, expirationTime time.Time) (string, error) {
	claims := &Claims{
		Login: userClaims.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ReadBody(req *http.Request) (string, error) {
	var reader io.Reader
	if req.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(req.Body)
		if err != nil {
			return "", err
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = req.Body
	}
	buf, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func Logout(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:    "GOPHERMART-SESSION",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-1 * time.Hour),
	}
	http.SetCookie(w, cookie)
}

func Login(w http.ResponseWriter, l string, expire time.Duration) error {

	e := time.Now().Add(expire)
	token, err := NewToken(&Claims{
		Login: l,
	}, e)

	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:    "GOPHERMART-SESSION",
		Value:   token,
		Path:    "/",
		Expires: e,
	}
	http.SetCookie(w, cookie)

	return nil
}

func Validate(req *http.Request) error {
	cookie, err := req.Cookie("GOPHERMART-SESSION")
	if err != nil {
		return err
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtSecretKey, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return ErrJWTValidation
	}

	return nil
}
