//go:build integration
// +build integration

package repository

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"tiger-tracker-api/constants"
	"tiger-tracker-api/repository/helper"
	"tiger-tracker-api/repository/models"
	"time"
)

type appRepositoryTestSuite struct {
	suite.Suite
	context        *gin.Context
	db             *sqlx.DB
	repoTestHelper helper.RepoTestHelper
	appRepository  AppRepository
}

func TestAppRepositoryIntegrationTest(t *testing.T) {
	suite.Run(t, new(appRepositoryTestSuite))
}

func (suite *appRepositoryTestSuite) SetupSuite() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv(constants.DB_USER),
		os.Getenv(constants.DB_PASSWORD),
		os.Getenv(constants.DB_HOST),
		os.Getenv(constants.DB_PORT),
		os.Getenv(constants.DB_NAME))

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintln("Could not connect to DB", err))
	}
	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintln("Could not ping DB", err))
	}
	suite.db = db
	suite.appRepository = NewAppRepository(suite.db)
	suite.repoTestHelper = helper.NewRepoTestHelper(suite.db)
	suite.context, _ = gin.CreateTestContext(httptest.NewRecorder())
	suite.context.Request, _ = http.NewRequest("GET", "url", nil)
}

func (suite *appRepositoryTestSuite) TearDownTest() {
	_ = suite.repoTestHelper.ClearAllData()
}

func (suite *appRepositoryTestSuite) TestGetRecentTigerSightings_ShouldReturnListOfTigerSightingsInDescendingOrder_WithPagination() {
	tigers := []models.TigerDetails{
		{Name: "raja", DateOfBirth: time.Date(1997, 04, 01, 0, 0, 0, 0, time.UTC)},
		{Name: "shera", DateOfBirth: time.Date(1998, 1, 02, 0, 0, 0, 0, time.UTC)},
		{Name: "timon", DateOfBirth: time.Date(1992, 7, 15, 0, 0, 0, 0, time.UTC)},
		{Name: "pumba", DateOfBirth: time.Date(1995, 9, 30, 0, 0, 0, 0, time.UTC)},
		{Name: "pooh", DateOfBirth: time.Date(2000, 2, 27, 0, 0, 0, 0, time.UTC)},
	}
	for idx, tiger := range tigers {
		id, err := suite.repoTestHelper.CreateTiger(tiger)
		suite.Nil(err)
		tigers[idx].Id = id
	}
	users := []models.UserDetails{
		{UserName: "test-user", PasswordHash: "hash", Email: "test@test.com"},
	}
	for idx, user := range users {
		id, err := suite.repoTestHelper.CreateUser(user)
		suite.Nil(err)
		users[idx].Id = id
	}
	time1, _ := time.Parse(time.RFC3339, "2024-05-14T01:59:59Z")
	time2, _ := time.Parse(time.RFC3339, "2024-05-14T02:59:59Z")
	time3, _ := time.Parse(time.RFC3339, "2024-05-14T03:59:59Z")
	time4, _ := time.Parse(time.RFC3339, "2024-05-14T04:59:59Z")
	time5, _ := time.Parse(time.RFC3339, "2024-05-14T05:59:59Z")
	time6, _ := time.Parse(time.RFC3339, "2024-05-14T06:59:59Z")
	time7, _ := time.Parse(time.RFC3339, "2024-05-14T07:59:59Z")
	time8, _ := time.Parse(time.RFC3339, "2024-05-14T08:59:59Z")
	time9, _ := time.Parse(time.RFC3339, "2024-05-14T09:59:59Z")
	time10, _ := time.Parse(time.RFC3339, "2024-05-14T10:59:59Z")

	sightings := []models.TigerSightings{
		{TigerId: tigers[0].Id, LastSeen: time1, LastPosition: models.Location{Latitude: 20.67890, Longitude: 90.54789}, ReportedBy: users[0].Id, Image: nil},
		{TigerId: tigers[0].Id, LastSeen: time2, LastPosition: models.Location{Latitude: 21.67890, Longitude: 91.54789}, ReportedBy: users[0].Id, Image: nil},
		{TigerId: tigers[1].Id, LastSeen: time3, LastPosition: models.Location{Latitude: 22.67890, Longitude: 92.54789}, ReportedBy: users[0].Id, Image: nil},
		{TigerId: tigers[1].Id, LastSeen: time4, LastPosition: models.Location{Latitude: 23.67890, Longitude: 93.54789}, ReportedBy: users[0].Id, Image: nil},
		{TigerId: tigers[2].Id, LastSeen: time5, LastPosition: models.Location{Latitude: 24.67890, Longitude: 94.54789}, ReportedBy: users[0].Id, Image: nil},
		{TigerId: tigers[2].Id, LastSeen: time6, LastPosition: models.Location{Latitude: 25.67890, Longitude: 95.54789}, ReportedBy: users[0].Id, Image: nil},
		{TigerId: tigers[3].Id, LastSeen: time7, LastPosition: models.Location{Latitude: 26.67890, Longitude: 96.54789}, ReportedBy: users[0].Id, Image: nil},
		{TigerId: tigers[3].Id, LastSeen: time8, LastPosition: models.Location{Latitude: 27.67890, Longitude: 97.54789}, ReportedBy: users[0].Id, Image: nil},
		{TigerId: tigers[4].Id, LastSeen: time9, LastPosition: models.Location{Latitude: 28.67890, Longitude: 98.54789}, ReportedBy: users[0].Id, Image: nil},
		{TigerId: tigers[4].Id, LastSeen: time10, LastPosition: models.Location{Latitude: 29.67890, Longitude: 99.54789}, ReportedBy: users[0].Id, Image: nil},
	}

	for idx, sighting := range sightings {
		id, err := suite.repoTestHelper.CreateTigerSighting(sighting)
		suite.Nil(err)
		sightings[idx].Id = id
	}
	actualSightings, err := suite.appRepository.GetRecentTigerSightings(suite.context, 1, 10)
	suite.Nil(err)
	suite.Equal(len(actualSightings), 5)
	expectedAllSightings := []models.TigerDetailWithSightings{
		{
			TigerId:      tigers[4].Id,
			TigerName:    tigers[4].Name,
			DateOfBirth:  tigers[4].DateOfBirth,
			LastSeen:     sightings[9].LastSeen,
			LastPosition: sightings[9].LastPosition,
			ReportedBy:   sightings[9].ReportedBy,
			Image:        sightings[9].Image,
		},
		{
			TigerId:      tigers[3].Id,
			TigerName:    tigers[3].Name,
			DateOfBirth:  tigers[3].DateOfBirth,
			LastSeen:     sightings[7].LastSeen,
			LastPosition: sightings[7].LastPosition,
			ReportedBy:   sightings[7].ReportedBy,
			Image:        sightings[7].Image,
		},
		{
			TigerId:      tigers[2].Id,
			TigerName:    tigers[2].Name,
			DateOfBirth:  tigers[2].DateOfBirth,
			LastSeen:     sightings[5].LastSeen,
			LastPosition: sightings[5].LastPosition,
			ReportedBy:   sightings[5].ReportedBy,
			Image:        sightings[5].Image,
		},
		{
			TigerId:      tigers[1].Id,
			TigerName:    tigers[1].Name,
			DateOfBirth:  tigers[1].DateOfBirth,
			LastSeen:     sightings[3].LastSeen,
			LastPosition: sightings[3].LastPosition,
			ReportedBy:   sightings[3].ReportedBy,
			Image:        sightings[3].Image,
		},
		{
			TigerId:      tigers[0].Id,
			TigerName:    tigers[0].Name,
			DateOfBirth:  tigers[0].DateOfBirth,
			LastSeen:     sightings[1].LastSeen,
			LastPosition: sightings[1].LastPosition,
			ReportedBy:   sightings[1].ReportedBy,
			Image:        sightings[1].Image,
		},
	}
	suite.Equal(expectedAllSightings, actualSightings)

	actualPaginatedSigtings, err := suite.appRepository.GetRecentTigerSightings(suite.context, 1, 2)
	suite.Equal(expectedAllSightings[0:2], actualPaginatedSigtings)

	actualPaginatedSigtings, err = suite.appRepository.GetRecentTigerSightings(suite.context, 2, 2)
	suite.Equal(expectedAllSightings[2:4], actualPaginatedSigtings)

	actualPaginatedSigtings, err = suite.appRepository.GetRecentTigerSightings(suite.context, 3, 2)
	suite.Equal(expectedAllSightings[4:], actualPaginatedSigtings)

	actualPaginatedSigtings, err = suite.appRepository.GetRecentTigerSightings(suite.context, 4, 2)
	suite.Empty(actualPaginatedSigtings)
}
