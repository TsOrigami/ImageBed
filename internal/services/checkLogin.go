package services

import (
	usr "ImageV2/internal/db/user"
	"errors"
	"net/http"
	"strings"
)

func CheckLogin(r *http.Request) error {
	token := r.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}
	if token == "" {
		return errors.New("unauthorized")
	} else {
		if !usr.CheckToken(token) {
			return errors.New("unauthorized")
		}
	}
	return nil
}
