package models

import (
	"time"
)

type TigerDetailWithSightings struct {
	TigerId      int       `db:"tiger_id"`
	TigerName    string    `db:"name"`
	DateOfBirth  time.Time `db:"dob"`
	LastSeen     time.Time `db:"last_seen"`
	LastPosition Location  `db:"last_position"`
	ReportedBy   int       `db:"user_id"`
	Image        []byte    `db:"image"`
}
