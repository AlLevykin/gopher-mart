package rest

import (
	"compress/gzip"
	"io"
	"net/http"
	"time"
)

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

func SetCookie(w http.ResponseWriter, l string, expire time.Duration) {
	cookie := &http.Cookie{
		Name:    "GOPHERMART-SESSION",
		Value:   l,
		Path:    "/",
		Expires: time.Now().Add(expire),
	}
	http.SetCookie(w, cookie)
}
