package repositories

import (
	"llio-api/database"
	"llio-api/models/DAOs"
)

func CreateProject(project *DAOs.Project) (*DAOs.Project, error) {
	err := database.DB.Create(project).Error
	return project, DBErrorManager(err)
}

func GetProjects() ([]*DAOs.Project, error) {
	var projects []*DAOs.Project
	err := database.DB.Find(&projects).Error
	return projects, DBErrorManager(err)
}

func GetProjectActivities(projectId int, from string, to string) ([]DAOs.ActivityWithTimeSpent, error) {
	var tempResults []DAOs.ActivityWithTimeSpent
	var err error

	if from == "" || to == "" {
		err = database.DB.
			Select(`
            activities.user_id, 
            activities.category_id, 
            activities.project_id, 
            CAST(SUM(TIMESTAMPDIFF(SECOND, activities.start_date, activities.end_date))/3600.0 AS DECIMAL(10,2)) as time_spent
        `).
			Table("activities").
			Where("project_id = ?", projectId).
			Group("user_id, category_id, project_id").
			Scan(&tempResults).Error
		if err != nil {
			return nil, DBErrorManager(err)
		}
	} else if from != "" && to != "" {
		fromDate := from
		toDate := to

		if from == to {
			toDate = to + " 23:59:59"
			fromDate = from + " 00:00:00"
		}

		fromDate = from + "T00:00:00"
		toDate = to + "T23:59:59"

		err = database.DB.
			Select(`
				activities.user_id, 
				activities.category_id, 
				activities.project_id, 
				CAST(SUM(TIMESTAMPDIFF(SECOND, activities.start_date, activities.end_date))/3600.0 AS DECIMAL(10,2)) as time_spent
			`).
			Table("activities").
			Where("project_id = ? AND Start_Date >= ? AND End_Date <= ?", projectId, fromDate, toDate).
			Group("user_id, category_id, project_id").
			Scan(&tempResults).Error
	}

	if err != nil {
		return nil, DBErrorManager(err)
	}

	return tempResults, err
}

func GetProjectActivitiesFromUser(projectId int, userId *int, from string, to string) ([]DAOs.ActivityWithTimeSpent, error) {
	var tempResults []DAOs.ActivityWithTimeSpent
	var err error

	if from == "" || to == "" {
		err = database.DB.
			Select(`
				activities.user_id, 
				activities.category_id, 
				activities.project_id, 
				CAST(SUM(TIMESTAMPDIFF(SECOND, activities.start_date, activities.end_date))/3600.0 AS DECIMAL(10,2)) as time_spent
			`).
			Table("activities").
			Where("project_id = ?", projectId).
			Where("user_id = ?", userId).
			Group("user_id, category_id, project_id").
			Scan(&tempResults).Error
	} else if from != "" && to != "" {
		fromDate := from
		toDate := to

		if from == to {
			toDate = to + " 23:59:59"
			fromDate = from + " 00:00:00"
		}

		fromDate = from + "T00:00:00"
		toDate = to + "T23:59:59"

		err = database.DB.
			Select(`
				activities.user_id, 
				activities.category_id, 
				activities.project_id, 
				CAST(SUM(TIMESTAMPDIFF(SECOND, activities.start_date, activities.end_date))/3600.0 AS DECIMAL(10,2)) as time_spent
			`).
			Table("activities").
			Where("project_id = ? AND Start_Date >= ? AND End_Date <= ?", projectId, fromDate, toDate).
			Where("user_id = ?", userId).
			Group("user_id, category_id, project_id").
			Scan(&tempResults).Error
	}

	if err != nil {
		return nil, DBErrorManager(err)
	}

	return tempResults, err
}

func GetProjectsByManagerId(id int) ([]*DAOs.Project, error) {
	var projects []*DAOs.Project
	err := database.DB.Find(&projects, "manager_id = ?", id).Error
	return projects, DBErrorManager(err)
}

func GetProjectsByUserId(id int) ([]*DAOs.Project, error) {
	var projects []*DAOs.Project
	err := database.DB.Distinct("projects.*").Joins("JOIN activities ON activities.project_id = projects.id").
		Where("activities.user_id = ?", 1).
		Find(&projects).Error
	return projects, DBErrorManager(err)
}

func GetProjectById(id string) (*DAOs.Project, error) {
	var project DAOs.Project

	err := database.DB.First(&project, id).Error
	return &project, DBErrorManager(err)
}

// Les categories sont supprimés en casade delete par la DB
func DeleteProjectById(id int) error {
	err := database.DB.Delete(&DAOs.Project{}, id).Error
	return DBErrorManager(err)
}

func UpdateProject(projectDAO *DAOs.Project) (*DAOs.Project, error) {
	// AJOUT : .Select(...) force la mise à jour de ces colonnes, même si la valeur est 0 ou false
	err := database.DB.Model(projectDAO).
		Select("ManagerId", "UniqueId", "Name", "Status", "Billable", "UpdatedAt", "EndAt", "EstimatedHours").
		Updates(projectDAO).Error

	return projectDAO, DBErrorManager(err)
}

func GetProjectsByActivityPerUser(userId int) ([]*DAOs.Project, error) {
	var projects []*DAOs.Project
	err := database.DB.Distinct("projects.*").Joins("JOIN activities ON activities.project_id = projects.id").
		Where("activities.user_id = ?", userId).
		Find(&projects).Error
	return projects, DBErrorManager(err)
}

func ProjectHasActivities(id int) (bool, error) {
	var count int64
	err := database.DB.Model(&DAOs.Activity{}).Where("project_id = ?", id).Count(&count).Error
	if err != nil {
		return false, DBErrorManager(err)
	}
	return count > 0, nil
}
