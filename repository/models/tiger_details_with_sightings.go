package models

import "time"

type TigerDetailWithSightings struct {
	TigerId      int64     `db:"tiger_id" json:"-"`
	TigerName    string    `db:"name" json:"tigerName"`
	DateOfBirth  time.Time `db:"dob" json:"dob"`
	LastSeen     time.Time `db:"last_seen" json:"lastSeen"`
	LastPosition Location  `db:"last_position" json:"lastPosition"`
	ReportedBy   int64     `db:"user_id" json:"-"`
	Image        []byte    `db:"image" json:"-"`
}
