package models

type UserDetails struct {
	Id           int64  `db:"Id"`
	UserName     string `db:"user_name"`
	PasswordHash string `db:"password_hash"`
	Email        string `db:"email"`
}
