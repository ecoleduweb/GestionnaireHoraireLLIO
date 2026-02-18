package repositories

import (
	"llio-api/database"
	"llio-api/models/DAOs"

	"gorm.io/gorm"
)

func CreateProject(project *DAOs.Project) (*DAOs.Project, error) {
	err := database.DB.Create(project).Error
	return project, DBErrorManager(err)
}
func applyRecentActivityOrdering(db *gorm.DB) *gorm.DB {
	return db.
		Joins(`
			LEFT JOIN (
				SELECT project_id, MAX(start_date) AS last_activity
				FROM activities
				GROUP BY project_id
			) AS last_activities ON last_activities.project_id = projects.id
		`).
		Order(`
			CASE 
				WHEN last_activities.last_activity IS NOT NULL 
				AND last_activities.last_activity > DATE_SUB(NOW(), INTERVAL 7 DAY) 
				THEN 0 ELSE 1 
			END
		`).
		Order(`
			CASE 
				WHEN last_activities.last_activity IS NOT NULL 
				AND last_activities.last_activity > DATE_SUB(NOW(), INTERVAL 7 DAY) 
				THEN last_activities.last_activity 
			END DESC
		`).
		Order(`
			CASE 
				WHEN last_activities.last_activity IS NULL 
				OR last_activities.last_activity <= DATE_SUB(NOW(), INTERVAL 7 DAY) 
				THEN LOWER(projects.name) 
			END ASC
		`)
}

func GetProjects() ([]*DAOs.Project, error) {
	var projects []*DAOs.Project

	query := database.DB.
		Table("projects").
		Select("DISTINCT projects.*")

	query = applyRecentActivityOrdering(query)

	err := query.Find(&projects).Error
	return projects, DBErrorManager(err)
}

func GetProjectActivities(projectId int) ([]DAOs.ActivityWithTimeSpent, error) {
	var tempResults []DAOs.ActivityWithTimeSpent

	err := database.DB.
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

	return tempResults, err
}

func GetProjectActivitiesFromUser(projectId int, userId *int) ([]DAOs.ActivityWithTimeSpent, error) {
	var tempResults []DAOs.ActivityWithTimeSpent

	err := database.DB.
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
	err := database.DB.
		Distinct("projects.*").
		Joins("JOIN activities ON activities.project_id = projects.id").
		Where("activities.user_id = ?", id).
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

	query := database.DB.
		Select("DISTINCT projects.*").
		Joins("JOIN activities ON activities.project_id = projects.id").
		Where("activities.user_id = ?", userId)

	query = applyRecentActivityOrdering(query)

	err := query.Find(&projects).Error
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


