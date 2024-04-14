package models

import "time"

type TigerSightings struct {
	Id           int64     `db:"id"`
	TigerId      int64     `db:"tiger_id"`
	LastSeen     time.Time `db:"last_seen"`
	LastPosition Location  `db:"last_position"`
	ReportedBy   int64     `db:"user_id"`
	Image        []byte    `db:"image"`
}
