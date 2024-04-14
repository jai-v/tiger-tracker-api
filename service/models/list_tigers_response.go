package models

import "tiger-tracker-api/repository/models"

type ListTigersResponse struct {
	Tigers []models.TigerDetailWithSightings `json:"tigers"`
}
