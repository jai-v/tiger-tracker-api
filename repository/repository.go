package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"tiger-tracker-api/repository/models"
)

//go:generate mockgen -source=./repository.go -destination=./mocks/mock_repository.go -package=mocks

type AppRepository interface {
	GetRecentTigerSightings(ctx *gin.Context, pageNumber, pageSize int) ([]models.TigerDetailWithSightings, error)
	CreateUser(ctx *gin.Context, userName, email, passwordHash string) (int64, error)
	CheckIfUsernameOrEmailExists(ctx *gin.Context, userName, email string) (int64, error)
	GetUserDetailsByUsername(ctx *gin.Context, userName string) (models.UserDetails, error)
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

func (ar appRepository) CreateUser(ctx *gin.Context, userName, email, passwordHash string) (int64, error) {
	query := `insert into user_details(user_name, email, password_hash) values(?,?,?)`
	result, err := ar.db.ExecContext(ctx, query, userName, email, passwordHash)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (ar appRepository) CheckIfUsernameOrEmailExists(ctx *gin.Context, userName, email string) (int64, error) {
	query := `select count(1) from user_details where user_name=? or email=?`
	var count int64
	err := ar.db.GetContext(ctx, &count, query, userName, email)
	return count, err
}

func (ar appRepository) GetUserDetailsByUsername(ctx *gin.Context, userName string) (models.UserDetails, error) {
	query := `select id, user_name, email, password_hash from user_details where user_name=?`
	var userDetails models.UserDetails
	err := ar.db.GetContext(ctx, &userDetails, query, userName)
	return userDetails, err
}
