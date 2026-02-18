package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"llio-api/database"
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"llio-api/models/enums"

	"github.com/stretchr/testify/assert"
)
func TestGetProjectsSortedByRecentActivity(t *testing.T) {
	now := time.Now()

	
	var existingUser DAOs.User
	err := database.DB.First(&existingUser).Error
	assert.NoError(t, err, "Un user doit exister en BD")
	assert.NotZero(t, existingUser.Id)

	
	var existingCategory DAOs.Category
	err = database.DB.First(&existingCategory).Error
	assert.NoError(t, err, "Une catégorie doit exister en BD")
	assert.NotZero(t, existingCategory.Id)

	
	recentProjectBody := DTOs.ProjectDTO{
		UniqueId:  "Recent-Sort-API-001",
		ManagerId: existingUser.Id,
		Name:      "Projet activité récente API",
		Status:    enums.ProjectStatus(enums.InProgress),
	}
	w := sendRequest(router, "POST", "/project", recentProjectBody, enums.Administrator)
	assertResponse(t, w, http.StatusCreated, nil)
	var recentProjectResp struct {
		Project DAOs.Project `json:"project"`
	}
	json.Unmarshal(w.Body.Bytes(), &recentProjectResp)
	assert.NotZero(t, recentProjectResp.Project.Id)

	
	oldProjectBody := DTOs.ProjectDTO{
		UniqueId:  "Old-Sort-API-001",
		ManagerId: existingUser.Id,
		Name:      "Projet activité ancienne API",
		Status:    enums.ProjectStatus(enums.InProgress),
	}
	w = sendRequest(router, "POST", "/project", oldProjectBody, enums.Administrator)
	assertResponse(t, w, http.StatusCreated, nil)
	var oldProjectResp struct {
		Project DAOs.Project `json:"project"`
	}
	json.Unmarshal(w.Body.Bytes(), &oldProjectResp)
	assert.NotZero(t, oldProjectResp.Project.Id)

	
	recentCatBody := DTOs.CategoryDTO{
		Name:        "Cat Sort Recent API",
		Description: "Categorie test tri recent",
		ProjectId:   recentProjectResp.Project.Id,
	}
	w = sendRequest(router, "POST", "/category", recentCatBody, enums.Administrator)
	assertResponse(t, w, http.StatusCreated, nil)
	var recentCatResp struct {
		Category DAOs.Category `json:"category"`
	}
	json.Unmarshal(w.Body.Bytes(), &recentCatResp)
	assert.NotZero(t, recentCatResp.Category.Id)

	
	oldCatBody := DTOs.CategoryDTO{
		Name:        "Cat Sort Old API",
		Description: "Categorie test tri ancien",
		ProjectId:   oldProjectResp.Project.Id,
	}
	w = sendRequest(router, "POST", "/category", oldCatBody, enums.Administrator)
	assertResponse(t, w, http.StatusCreated, nil)
	var oldCatResp struct {
		Category DAOs.Category `json:"category"`
	}
	json.Unmarshal(w.Body.Bytes(), &oldCatResp)
	assert.NotZero(t, oldCatResp.Category.Id)

	
	recentActivity := DTOs.ActivityDTO{
		Name:        "Activite recente tri API",
		Description: "Activité il y a 3 jours",
		StartDate:   now.Add(-3 * 24 * time.Hour),
		EndDate:     now.Add(-3*24*time.Hour + 2*time.Hour),
		UserId:      existingUser.Id,
		ProjectId:   recentProjectResp.Project.Id,
		CategoryId:  recentCatResp.Category.Id,
	}
	w = sendRequest(router, "POST", "/activity", recentActivity, enums.Administrator)
	assertResponse(t, w, http.StatusCreated, nil)

	
	oldActivity := DTOs.ActivityDTO{
		Name:        "Activite ancienne tri API",
		Description: "Activité il y a 1 mois",
		StartDate:   now.Add(-33 * 24 * time.Hour),
		EndDate:     now.Add(-33*24*time.Hour + 2*time.Hour),
		UserId:      existingUser.Id,
		ProjectId:   oldProjectResp.Project.Id,
		CategoryId:  oldCatResp.Category.Id,
}
	w = sendRequest(router, "POST", "/activity", oldActivity, enums.Administrator)
	assertResponse(t, w, http.StatusCreated, nil)


	w = sendRequest(router, "GET", "/projects?sortBy=recentActivity", nil, enums.Administrator)
	assertResponse(t, w, http.StatusOK, nil)

	var projectsBody struct {
		Projects []DAOs.Project `json:"projects"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &projectsBody)
	assert.NoError(t, err)
	assert.True(t, len(projectsBody.Projects) >= 2)

	recentIndex := -1
	oldIndex := -1
	for i, p := range projectsBody.Projects {
		if p.Id == recentProjectResp.Project.Id {
			recentIndex = i
		}
		if p.Id == oldProjectResp.Project.Id {
			oldIndex = i
		}
	}

	assert.NotEqual(t, -1, recentIndex, "Le projet récent doit être présent")
	assert.NotEqual(t, -1, oldIndex, "Le projet ancien doit être présent")
	assert.Less(t, recentIndex, oldIndex,
		"Le projet avec activité récente (3 jours) doit apparaître avant l'ancien (33 jours)")
}