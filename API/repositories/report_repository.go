package repositories

import (
	"fmt"
	"llio-api/database"
	"llio-api/models/DAOs"
)

func GetAllForExport(from string, to string) ([]DAOs.Activity, error) {
	var activities []DAOs.Activity

	err := database.DB.Find(&activities).Error

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return activities, nil
}