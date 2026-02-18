package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"llio-api/database"
	"llio-api/models/DAOs"
	"llio-api/models/enums"

	"github.com/stretchr/testify/assert"
)
func TestGetProjectsSortedByRecentActivity(t *testing.T) {
	now := time.Now()

	// Créer un user dédié directement en BD
	testUser := DAOs.User{
		FirstName: "Filter",
		LastName:  "Test",
		Email:     "filter.test@example.com",
	}
	database.DB.Where("email = ?", testUser.Email).FirstOrCreate(&testUser)
	database.DB.Where("email = ?", testUser.Email).First(&testUser)
	assert.NotZero(t, testUser.Id, "L'Id du user de test doit être non nul")

	// Créer le projet "récent" directement en BD
	recentProject := DAOs.Project{
		UniqueId:  "Recent-Sort-001",
		ManagerId: testUser.Id,
		Name:      "Projet activité récente",
		Status:    enums.ProjectStatus(enums.InProgress),
	}
	database.DB.Where("unique_id = ?", recentProject.UniqueId).FirstOrCreate(&recentProject)
	database.DB.Where("unique_id = ?", recentProject.UniqueId).First(&recentProject)
	assert.NotZero(t, recentProject.Id, "L'Id du projet récent doit être non nul")

	// Créer le projet "ancien" directement en BD
	oldProject := DAOs.Project{
		UniqueId:  "Old-Sort-001",
		ManagerId: testUser.Id,
		Name:      "Projet activité ancienne",
		Status:    enums.ProjectStatus(enums.InProgress),
	}
	database.DB.Where("unique_id = ?", oldProject.UniqueId).FirstOrCreate(&oldProject)
	database.DB.Where("unique_id = ?", oldProject.UniqueId).First(&oldProject)
	assert.NotZero(t, oldProject.Id, "L'Id du projet ancien doit être non nul")

	// Catégorie pour le projet récent
	testCategory := DAOs.Category{
		Name:        "Categorie Sort Test Recent",
		Description: "Catégorie pour test de tri",
		ProjectId:   recentProject.Id,
	}
	database.DB.Where("name = ? AND project_id = ?", testCategory.Name, testCategory.ProjectId).FirstOrCreate(&testCategory)
	database.DB.Where("name = ? AND project_id = ?", testCategory.Name, testCategory.ProjectId).First(&testCategory)
	assert.NotZero(t, testCategory.Id, "L'Id de la catégorie récente doit être non nul")

	// Catégorie pour le projet ancien
	testCategory2 := DAOs.Category{
		Name:        "Categorie Sort Test Old",
		Description: "Catégorie pour test de tri",
		ProjectId:   oldProject.Id,
	}
	database.DB.Where("name = ? AND project_id = ?", testCategory2.Name, testCategory2.ProjectId).FirstOrCreate(&testCategory2)
	database.DB.Where("name = ? AND project_id = ?", testCategory2.Name, testCategory2.ProjectId).First(&testCategory2)
	assert.NotZero(t, testCategory2.Id, "L'Id de la catégorie ancienne doit être non nul")

	// Activité récente (3 jours) directement en BD
	recentActivity := DAOs.Activity{
		Name:        "Activite recente tri",
		Description: "Activité il y a 3 jours",
		StartDate:   now.Add(-3 * 24 * time.Hour),
		EndDate:     now.Add(-3*24*time.Hour + 2*time.Hour),
		UserId:      testUser.Id,
		ProjectId:   recentProject.Id,
		CategoryId:  testCategory.Id,
	}
	database.DB.Where("name = ? AND project_id = ?", recentActivity.Name, recentActivity.ProjectId).FirstOrCreate(&recentActivity)

	// Activité ancienne (10 jours) directement en BD
	oldActivity := DAOs.Activity{
		Name:        "Activite ancienne tri",
		Description: "Activité il y a 10 jours",
		StartDate:   now.Add(-10 * 24 * time.Hour),
		EndDate:     now.Add(-10*24*time.Hour + 2*time.Hour),
		UserId:      testUser.Id,
		ProjectId:   oldProject.Id,
		CategoryId:  testCategory2.Id,
	}
	database.DB.Where("name = ? AND project_id = ?", oldActivity.Name, oldActivity.ProjectId).FirstOrCreate(&oldActivity)

	// Récupérer la liste triée par activité récente
	w := sendRequest(router, "GET", "/projects?sortBy=recentActivity", nil, enums.Administrator)
	assertResponse(t, w, http.StatusOK, nil)

	var projectsBody struct {
		Projects []DAOs.Project `json:"projects"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &projectsBody)
	assert.NoError(t, err)
	assert.True(t, len(projectsBody.Projects) >= 2, "La liste doit contenir au moins 2 projets")

	recentProjectIndex := -1
	oldProjectIndex := -1
	for i, p := range projectsBody.Projects {
		if p.Id == recentProject.Id {
			recentProjectIndex = i
		}
		if p.Id == oldProject.Id {
			oldProjectIndex = i
		}
	}

	assert.NotEqual(t, -1, recentProjectIndex, "Le projet récent doit être présent dans la liste")
	assert.NotEqual(t, -1, oldProjectIndex, "Le projet ancien doit être présent dans la liste")
	assert.Less(t, recentProjectIndex, oldProjectIndex,
		"Le projet avec une activité récente (3 jours) doit apparaître avant celui avec une activité ancienne (10 jours)")
}