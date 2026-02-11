package services

import (
	"fmt"
	"llio-api/customs_errors"
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"llio-api/models/enums"
	"llio-api/repositories"
	"log"

	"github.com/jinzhu/copier"
)

func VerifyProjectJSON(projectDTO *DTOs.ProjectDTO) []DTOs.FieldErrorDTO {
	var errors []DTOs.FieldErrorDTO

	if !projectDTO.CreatedAt.IsZero() {
		// Vérifier que CreatedAt est avant EndAt
		if !projectDTO.EndAt.IsZero() && projectDTO.CreatedAt.After(projectDTO.EndAt) {
			errors = append(errors, DTOs.FieldErrorDTO{
				Field:   "endAt",
				Message: "La date de fin doit être après la date de création",
			})
		}

		// Vérifier que CreatedAt est avant UpdatedAt
		if !projectDTO.UpdatedAt.IsZero() && projectDTO.CreatedAt.After(projectDTO.UpdatedAt) {
			errors = append(errors, DTOs.FieldErrorDTO{
				Field:   "updatedAt",
				Message: "La date de mise à jour doit être après la date de création",
			})
		}
	}

	return errors
}

func CreateProject(projectDTO *DTOs.ProjectDTO) (*DTOs.ProjectDTO, error) {

	project := &DAOs.Project{}
	err := copier.Copy(project, projectDTO)
	if err != nil {
		return nil, err
	}

	projectDAOAdded, err := repositories.CreateProject(project)
	if err != nil {
		return nil, err
	}

	_, err = repositories.CreateCategory(&DAOs.Category{
		Name:        "Par défaut",
		Description: "Catégorie par défaut",
		ProjectId:   projectDAOAdded.Id,
		CreatedAt:   projectDAOAdded.CreatedAt,
		UpdatedAt:   projectDAOAdded.UpdatedAt,
		Activities:  []DAOs.Activity{},
	})
	if err != nil {
		return nil, err
	}
	projectDTOResponse := &DTOs.ProjectDTO{}
	err = copier.Copy(projectDTOResponse, projectDAOAdded)
	return projectDTOResponse, err
}

func GetProjects() ([]*DTOs.ProjectDTO, error) {
	projects, err := repositories.GetProjects()
	if err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		return make([]*DTOs.ProjectDTO, 0), nil
	}

	var projectsDTO []*DTOs.ProjectDTO
	if err := copier.Copy(&projectsDTO, &projects); err != nil {
		return nil, err
	}

	return projectsDTO, nil
}

func GetDetailedProjects() ([]map[string]any, error) {
	projects, err := repositories.GetProjects()
	if err != nil {
		return nil, err
	}

	return formatProjects(projects, nil)
}

func GetDetailedProjectsByManagerId(id int) ([]map[string]any, error) {
	projects, err := repositories.GetProjectsByManagerId(id)
	if err != nil {
		return nil, err
	}

	return formatProjects(projects, nil)
}

func GetDetailedProjectsByUserId(id int) ([]map[string]any, error) {
	projects, err := repositories.GetProjectsByActivityPerUser(id)
	if err != nil {
		return nil, err
	}
	return formatProjects(projects, &id)
}

func GetProjectById(id string) (*DTOs.ProjectDTO, error) {
	project, err := repositories.GetProjectById(id)
	if err != nil {
		return nil, err
	}

	projectDTO := &DTOs.ProjectDTO{}
	err = copier.Copy(projectDTO, project)

	return projectDTO, err
}

func DeleteProjectById(id int) error {
	// Check if the project exists
	projectDAO, err := repositories.GetProjectById(fmt.Sprintf("%d", id))

	if err != nil {
		return err
	}
	if projectDAO == nil {
		return customs_errors.ErrProjectNotFound
	}

	projectHasActivities, err := repositories.ProjectHasActivities(id)
	if err != nil {
		return err
	}
	if projectHasActivities {
		log.Printf("Project %d has activities: %v", id, projectHasActivities)
		return customs_errors.ErrProjectHasActivities
	}

	projectCategoriesDeleted, err := repositories.ProjectDeleteCategories(id)
	if err != nil {
		return err
	}

	if projectCategoriesDeleted {
		log.Printf("Deleted categories: %v from project: %d", projectCategoriesDeleted, id)
	}

	//Delete the project
	errDelete := repositories.DeleteProjectById(id)
	if errDelete != nil {
		return errDelete
	}

	return nil
}

func UpdateProject(projectDTO *DTOs.ProjectDTO) (*DTOs.ProjectDTO, error) {

	projectDAO := &DAOs.Project{}
	err := copier.Copy(projectDAO, projectDTO)
	if err != nil {
		return nil, err
	}

	projectDAOUpdated, err := repositories.UpdateProject(projectDAO)
	if err != nil {
		return nil, err
	}

	projectDTOResponse := &DTOs.ProjectDTO{}
	err = copier.Copy(projectDTOResponse, projectDAOUpdated)
	return projectDTOResponse, err
}

func formatProjects(projects []*DAOs.Project, userId *int) ([]map[string]any, error) {
	users, err := repositories.GetAllUsers()
	if err != nil {
		return nil, err
	}
	userMap := make(map[int]*DAOs.User)
	for _, user := range users {
		userMap[user.Id] = user
	}

	categories, err := repositories.GetCategories()
	if err != nil {
		return nil, err
	}
	categoryMap := make(map[int]*DAOs.Category)
	for _, cat := range categories {
		categoryMap[cat.Id] = cat
	}

	var result []map[string]any
	for _, project := range projects {
		var tempActivities []DAOs.ActivityWithTimeSpent
		if userId != nil {
			tempActivities, err = repositories.GetProjectActivitiesFromUser(project.Id, userId)
		} else {
			tempActivities, err = repositories.GetProjectActivities(project.Id)
		}
		if err != nil {
			return nil, err
		}

		activities := make([]DAOs.Activity, len(tempActivities))
		for i, result := range tempActivities {
			activities[i] = DAOs.Activity{
				UserId:     result.UserID,
				CategoryId: result.CategoryID,
				ProjectId:  result.ProjectID,
				TimeSpent:  result.TimeSpent,
			}
		}

		formattedProject := formatProjectWithActivities(project, activities, userMap, categoryMap)
		result = append(result, formattedProject)
	}

	return result, nil
}

func formatProjectWithActivities(project *DAOs.Project, activities []DAOs.Activity, userMap map[int]*DAOs.User, categoryMap map[int]*DAOs.Category) map[string]any {
	employeesMap := make(map[int]map[string]any)

	for _, activity := range activities {
		user, exists := userMap[activity.UserId]
		if !exists {
			continue
		}

		category, exists := categoryMap[activity.CategoryId]
		if !exists {
			continue
		}

		if _, ok := employeesMap[user.Id]; !ok {
			employeesMap[user.Id] = map[string]any{
				"name":       user.FirstName + " " + user.LastName,
				"categories": make([]map[string]any, 0),
			}
		}

		categories := employeesMap[user.Id]["categories"].([]map[string]any)
		found := false
		for i, cat := range categories {
			if cat["name"] == category.Name {
				categories[i]["timeSpent"] = activity.TimeSpent
				found = true
				break
			}
		}

		if !found {
			categories = append(categories, map[string]any{
				"name":          category.Name,
				"timeSpent":     activity.TimeSpent,
				"timeEstimated": 0,
			})
			employeesMap[user.Id]["categories"] = categories
		}
	}

	// Convert to array
	var employees []map[string]any
	for _, emp := range employeesMap {
		employees = append(employees, emp)
	}

	// Calculate total time
	var totalTimeSpent float64
	for _, emp := range employees {
		for _, cat := range emp["categories"].([]map[string]any) {
			totalTimeSpent += cat["timeSpent"].(float64)
		}
	}

	// Get manager info
	manager, exists := userMap[project.ManagerId]
	lead := ""
	if exists {
		lead = manager.FirstName + " " + manager.LastName
	}

	return map[string]any{
		"id":                 project.Id,
		"uniqueId":           project.UniqueId,
		"name":               project.Name,
		"lead":               lead,
		"coLeads":            []string{},
		"employees":          employees,
		"totalTimeEstimated": project.EstimatedHours,
		"totalTimeRemaining": float64(project.EstimatedHours) - totalTimeSpent,
		"totalTimeSpent":     totalTimeSpent,
		"isArchived":         project.Status == enums.ProjectStatus(enums.Finish),
		"managerId":          project.ManagerId,
		"createdAt":          project.CreatedAt,
		"updatedAt":          project.UpdatedAt,
		"billable":           project.Billable,
	}
}

// GetProjectsByActivityPerUser retrieves projects based on activities for the current user.
func GetProjectsByActivityPerUser(userId int) ([]*DTOs.ProjectDTO, error) {
	projects, err := repositories.GetProjectsByActivityPerUser(userId)
	if err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		return make([]*DTOs.ProjectDTO, 0), nil
	}

	var projectsDTO []*DTOs.ProjectDTO
	if err := copier.Copy(&projectsDTO, &projects); err != nil {
		return nil, err
	}

	return projectsDTO, nil
}
