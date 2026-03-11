package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"llio-api/models/enums"

	"github.com/stretchr/testify/assert"
)

func TestCreateProject(t *testing.T) {

	project := DTOs.ProjectDTO{
		UniqueId:  "Inter-1234",
		ManagerId: doNotDeleteUser.Id,
		Name:      "Nom de test",
		Status:    enums.ProjectStatus(enums.InProgress),
	}

	w := sendRequest(router, "POST", "/project", project, nil, enums.Administrator)
	assertResponse(t, w, http.StatusCreated, nil)

	var responseBody struct {
		Response string       `json:"response"`
		Project  DAOs.Project `json:"project"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "Le projet a bien été ajouté à la base de données", responseBody.Response)
	assert.Equal(t, project.Name, responseBody.Project.Name)
	assert.Equal(t, project.Status, responseBody.Project.Status)
}

func TestDoNotCreateProjectWithoutUniqueId(t *testing.T) {
	project := DTOs.ProjectDTO{
		UniqueId:  "",
		ManagerId: doNotDeleteUser.Id,
		Name:      "",
		Status:    enums.ProjectStatus(enums.InProgress),
	}

	w := sendRequest(router, "POST", "/project", project, nil, enums.Administrator)
	expectedErrors := []DTOs.FieldErrorDTO{
		{Field: "uniqueId", Message: "Le champ uniqueId est invalide ou manquant"},
	}
	assertResponse(t, w, http.StatusBadRequest, expectedErrors)
}

func TestDoNotCreateProjectWithInvalidUniqueId(t *testing.T) {
	project := DTOs.ProjectDTO{
		UniqueId:  "Ceci est un très long nom de projet qui dépasse largement la limite de 50 caractères et qui devrait donc être rejeté par la validation",
		ManagerId: doNotDeleteUser.Id,
		Name:      "Name de test",
		Status:    enums.ProjectStatus(enums.InProgress),
	}

	w := sendRequest(router, "POST", "/project", project, nil, enums.Administrator)
	expectedErrors := []DTOs.FieldErrorDTO{
		{Field: "uniqueId", Message: "Le champ uniqueId est invalide ou manquant"},
	}
	assertResponse(t, w, http.StatusBadRequest, expectedErrors)
}

func TestDoNotCreateProjectWithInconsistentDates(t *testing.T) {
	now := time.Now()

	project := DTOs.ProjectDTO{
		UniqueId:  "Test Project-1234",
		ManagerId: doNotDeleteUser.Id,
		Name:      "Name de test",
		Status:    enums.ProjectStatus(enums.InProgress),
		CreatedAt: now,
		EndAt:     now.Add(-24 * time.Hour), // Un jour avant CreatedAt
	}

	w := sendRequest(router, "POST", "/project", project, nil, enums.Administrator)
	expectedErrors := []DTOs.FieldErrorDTO{
		{Field: "endAt", Message: "La date de fin doit être après la date de création"},
	}
	assertResponse(t, w, http.StatusBadRequest, expectedErrors)
}

func TestCreateProjectWithConsistentDates(t *testing.T) {
	// Temps actuel
	now := time.Now()

	// Création d'un projet avec des dates cohérentes
	project := DTOs.ProjectDTO{
		UniqueId:  "Test Project Dates",
		ManagerId: doNotDeleteUser.Id,
		Name:      "Description avec dates",
		Status:    enums.ProjectStatus(enums.InProgress),
		CreatedAt: now,
		EndAt:     now.Add(7 * 24 * time.Hour), // Une semaine plus tard
	}

	w := sendRequest(router, "POST", "/project", project, nil, enums.Administrator)
	assertResponse(t, w, http.StatusCreated, nil)

	var responseBody struct {
		Response string       `json:"response"`
		Project  DAOs.Project `json:"project"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "Le projet a bien été ajouté à la base de données", responseBody.Response)
}
