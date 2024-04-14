package models

import "time"

type TigerDetails struct {
	Id          int64     `db:"id"`
	Name        string    `db:"name"`
	DateOfBirth time.Time `db:"dob"`
}
