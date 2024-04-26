package models

type UserDetails struct {
	Id           int64  `db:"id"`
	UserName     string `db:"user_name"`
	PasswordHash string `db:"password_hash"`
	Email        string `db:"email"`
}
