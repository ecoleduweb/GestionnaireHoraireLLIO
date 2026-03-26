package repositories

import (
	"fmt"
	"llio-api/database"
	"llio-api/models/DAOs"
)

func GetAllForExport() ([]DAOs.Activity, error) {
	var activities []DAOs.Activity

	err := database.DB.
		Preload("User").
		Preload("Project").
		Preload("Category").
		Order("start_date ASC").
		Find(&activities).Error

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return activities, nil
}