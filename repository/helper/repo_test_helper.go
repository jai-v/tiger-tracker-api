package helper

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"tiger-tracker-api/repository/models"
)

type RepoTestHelper struct {
	db *sqlx.DB
}

func NewRepoTestHelper(db *sqlx.DB) RepoTestHelper {
	return RepoTestHelper{db: db}
}

func (repoTestHelper RepoTestHelper) CreateTiger(tiger models.TigerDetails) (int64, error) {
	insertQuery := `insert into tiger_details(name, dob) values(?,?)`
	res, err := repoTestHelper.db.Exec(insertQuery, tiger.Name, tiger.DateOfBirth)
	if err != nil {
		return 0, err
	}
	lastInsertId, _ := res.LastInsertId()
	return lastInsertId, nil
}

func (repoTestHelper RepoTestHelper) CreateUser(user models.UserDetails) (int64, error) {
	insertQuery := `insert into user_details(user_name, password_hash, email) values(?,?,?)`
	res, err := repoTestHelper.db.Exec(insertQuery, user.UserName, user.PasswordHash, user.Email)
	if err != nil {
		return 0, err
	}
	lastInsertId, _ := res.LastInsertId()
	return lastInsertId, nil
}

func (repoTestHelper RepoTestHelper) CreateTigerSighting(sighting models.TigerSightings) (int64, error) {
	position := fmt.Sprintf("POINT(%f %f)", sighting.LastPosition.Latitude, sighting.LastPosition.Longitude)
	insertQuery := `insert into tiger_sightings(tiger_id, last_seen, last_position, user_id, image) values(?,?,ST_GeomFromText(?),?,?)`
	res, err := repoTestHelper.db.Exec(insertQuery, sighting.TigerId, sighting.LastSeen, position, sighting.ReportedBy, sighting.Image)
	if err != nil {
		return 0, err
	}
	lastInsertId, _ := res.LastInsertId()
	return lastInsertId, nil
}

func (repoTestHelper RepoTestHelper) ClearAllData() error {
	queries := []string{
		`delete from tiger_sightings`,
		`delete from user_details`,
		`delete from tiger_details`,
		`alter table tiger_sightings AUTO_INCREMENT = 1`,
		`alter table user_details AUTO_INCREMENT = 1`,
		`alter table tiger_details AUTO_INCREMENT = 1`,
	}

	for _, query := range queries {
		_, err := repoTestHelper.db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}
