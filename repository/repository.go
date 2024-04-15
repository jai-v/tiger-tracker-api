package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"tiger-tracker-api/repository/models"
)

//go:generate mockgen -source=./repository.go -destination=./mocks/mock_repository.go -package=mocks

type AppRepository interface {
	GetRecentTigerSightings(ctx *gin.Context, pageNumber, pageSize int) ([]models.TigerDetailWithSightings, error)
}

func NewAppRepository(db *sqlx.DB) AppRepository {
	return appRepository{db: db}
}

type appRepository struct {
	db *sqlx.DB
}

func (ar appRepository) GetRecentTigerSightings(ctx *gin.Context, pageNumber, pageSize int) ([]models.TigerDetailWithSightings, error) {
	offSet := (pageNumber - 1) * pageSize
	query := `select td.name, td.dob, rs.tiger_id, rs.last_seen, rs.last_position, rs.image, rs.user_id from tiger_details td right join
				(select tiger_sightings.* from tiger_sightings,
					(select tiger_id,max(last_seen) as last_seen from tiger_sightings ts2 group by ts2.tiger_id) as recent_sightings
				where tiger_sightings.tiger_id=recent_sightings.tiger_id and tiger_sightings.last_seen=recent_sightings.last_seen order by recent_sightings.last_seen desc limit ? offset ?) as rs
				on td.id=rs.tiger_id;`

	var tigerDetailsWithSightings []models.TigerDetailWithSightings
	err := ar.db.SelectContext(ctx, &tigerDetailsWithSightings, query, pageSize, offSet)
	return tigerDetailsWithSightings, err
}
