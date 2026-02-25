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
				SELECT project_id, MAX(start_date) AS latest_started_activity
				FROM activities
				GROUP BY project_id
			) AS last_activities ON last_activities.project_id = projects.id
		`).
		// permet de faire un premier tri des activitées selon celles qui datent de moins de 1mois
		Order(`
			CASE 
				WHEN last_activities.latest_started_activity IS NOT NULL 
				AND last_activities.latest_started_activity > DATE_SUB(NOW(), INTERVAL 1 MONTH) 
				THEN 0 ELSE 1 
			END
		`).
		// Prend les activitées de moins de 1 mois et les trient par date 
		Order(`
			CASE 
				WHEN last_activities.latest_started_activity IS NOT NULL 
				AND last_activities.latest_started_activity > DATE_SUB(NOW(), INTERVAL 1 MONTH) 
				THEN last_activities.latest_started_activity 
			END DESC
		`).
		// Trie le reste des activitées par odre alphabétique 
		Order(`
			CASE 
				WHEN last_activities.latest_started_activity IS NULL 
				OR last_activities.latest_started_activity <= DATE_SUB(NOW(), INTERVAL 1 MONTH) 
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


