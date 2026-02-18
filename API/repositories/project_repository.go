package repositories

import (
	"llio-api/database"
	"llio-api/models/DAOs"
)

// CreateProject ajoute un projet à la base de données
func CreateProject(project *DAOs.Project) (*DAOs.Project, error) {
	err := database.DB.Create(project).Error
	return project, DBErrorManager(err)
}

// GetProjects récupère les projets triés par activité récente
func GetProjects() ([]*DAOs.Project, error) {
	var projects []*DAOs.Project
	err := database.DB.
		Table("projects").
		Select("DISTINCT projects.*").
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
		`).
		Find(&projects).Error

	return projects, DBErrorManager(err)
}

// GetProjectActivities récupère les activités d'un projet avec le temps passé
func GetProjectActivities(projectId int) ([]DAOs.ActivityWithTimeSpent, error) {
	var tempResults []DAOs.ActivityWithTimeSpent
	err := database.DB.
		Select(`
			activities.user_id, 
			activities.category_id, 
			activities.project_id, 
			CAST(SUM(TIMESTAMPDIFF(SECOND, activities.start_date, activities.end_date))/3600.0 AS DECIMAL(10,2)) AS time_spent
		`).
		Table("activities").
		Where("project_id = ?", projectId).
		Group("user_id, category_id, project_id").
		Scan(&tempResults).Error

	if err != nil {
		return nil, DBErrorManager(err)
	}
	return tempResults, nil
}

// GetProjectActivitiesFromUser récupère les activités d'un projet pour un utilisateur spécifique
func GetProjectActivitiesFromUser(projectId int, userId *int) ([]DAOs.ActivityWithTimeSpent, error) {
	var tempResults []DAOs.ActivityWithTimeSpent
	err := database.DB.
		Select(`
			activities.user_id, 
			activities.category_id, 
			activities.project_id, 
			CAST(SUM(TIMESTAMPDIFF(SECOND, activities.start_date, activities.end_date))/3600.0 AS DECIMAL(10,2)) AS time_spent
		`).
		Table("activities").
		Where("project_id = ?", projectId).
		Where("user_id = ?", userId).
		Group("user_id, category_id, project_id").
		Scan(&tempResults).Error

	if err != nil {
		return nil, DBErrorManager(err)
	}
	return tempResults, nil
}

// GetProjectsByManagerId récupère les projets par manager
func GetProjectsByManagerId(id int) ([]*DAOs.Project, error) {
	var projects []*DAOs.Project
	err := database.DB.Find(&projects, "manager_id = ?", id).Error
	return projects, DBErrorManager(err)
}

// GetProjectsByUserId récupère les projets d'un utilisateur
func GetProjectsByUserId(id int) ([]*DAOs.Project, error) {
	var projects []*DAOs.Project
	err := database.DB.
		Distinct("projects.*").
		Joins("JOIN activities ON activities.project_id = projects.id").
		Where("activities.user_id = ?", id).
		Find(&projects).Error
	return projects, DBErrorManager(err)
}

// GetProjectById récupère un projet par son ID
func GetProjectById(id string) (*DAOs.Project, error) {
	var project DAOs.Project
	err := database.DB.First(&project, id).Error
	return &project, DBErrorManager(err)
}

// UpdateProject met à jour un projet
func UpdateProject(projectDAO *DAOs.Project) (*DAOs.Project, error) {
	err := database.DB.Model(projectDAO).
		Select("ManagerId", "UniqueId", "Name", "Status", "Billable", "UpdatedAt", "EndAt", "EstimatedHours").
		Updates(projectDAO).Error
	return projectDAO, DBErrorManager(err)
}

// GetProjectsByActivityPerUser récupère les projets d'un utilisateur triés par activité récente
func GetProjectsByActivityPerUser(userId int) ([]*DAOs.Project, error) {
	var projects []*DAOs.Project
	err := database.DB.
		Select("DISTINCT projects.*").
		Joins("JOIN activities ON activities.project_id = projects.id").
		Joins(`
			LEFT JOIN (
				SELECT project_id, MAX(start_date) AS last_activity
				FROM activities
				GROUP BY project_id
			) AS last_activities ON last_activities.project_id = projects.id
		`).
		Where("activities.user_id = ?", userId).
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
		`).
		Find(&projects).Error

	return projects, DBErrorManager(err)
}
