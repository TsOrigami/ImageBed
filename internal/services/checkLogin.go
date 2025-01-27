package services

import (
	"ImageV2/internal/db/user"
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
		_, err := user.GetUsername(token)
		if err != nil {
			return err
		}
	}
	return nil
}
