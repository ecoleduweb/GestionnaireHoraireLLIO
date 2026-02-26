package repositories

import (
	"llio-api/database"
	"llio-api/models/DAOs"
)

func GetCoManagersByProjectIds(projectIds []int) ([]*DAOs.CoManager, error) {
	var coManagers []*DAOs.CoManager
	if len(projectIds) == 0 {
		return coManagers, nil
	}
	err := database.DB.
		Where("project_id IN (?)", projectIds).
		Find(&coManagers).Error
	return coManagers, DBErrorManager(err)
}

func AddCoManagerToProject(coManager *DAOs.CoManager) (*DAOs.CoManager, error) {
	err := database.DB.Create(coManager).Error

	if err != nil {
		return nil, DBErrorManager(err)
	}

	return coManager, nil
}

func IsUserCoManager(projectId int, userId int) (bool, error) {
	var count int64
	err := database.DB.
		Model(&DAOs.CoManager{}).
		Where("project_id = ? AND user_id = ?", projectId, userId).
		Count(&count).Error

	if err != nil {
		return false, DBErrorManager(err)
	}
	return count > 0, nil
}
