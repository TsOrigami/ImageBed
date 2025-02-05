package user

import (
	dbInfo "ImageV2/internal/db/sql"
)

func GetAmount(account string) (int, error) {
	var (
		count int
		err   error
	)
	db, err := dbInfo.GetDB()
	if err != nil {
		return 0, err
	}
	err = db.Connect.QueryRow(`SELECT count FROM user_info WHERE account = ?`, account).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
