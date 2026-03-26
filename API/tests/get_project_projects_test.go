package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"llio-api/models/DAOs"
	"llio-api/models/enums"

	"github.com/stretchr/testify/assert"
)

func TestGetProjectById(t *testing.T) {
	url := fmt.Sprintf("/project/%d", doNotDeleteProject.Id)
	w := sendRequest(router, "GET", url, nil, nil)
	assertResponse(t, w, http.StatusOK, nil)

	var responseBody struct {
		Project DAOs.Project `json:"project"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, doNotDeleteProject.Id, responseBody.Project.Id)
	assert.Equal(t, doNotDeleteProject.UniqueId, responseBody.Project.UniqueId)
	assert.Equal(t, doNotDeleteProject.Name, responseBody.Project.Name)
}

func TestGetProjectWithInvalidId(t *testing.T) {
	w := sendRequest(router, "GET", "/project/999999", nil, nil)
	assertResponse(t, w, http.StatusNotFound, nil)
}

func TestGetProjectWithNonNumericId(t *testing.T) {
	w := sendRequest(router, "GET", "/project/invalid", nil, nil)
	assertResponse(t, w, http.StatusInternalServerError, nil)
}

func TestGetAllProjects(t *testing.T) {
	w := sendRequest(router, "GET", "/projects", nil, nil, enums.Administrator)
	assertResponse(t, w, http.StatusOK, nil)

	// Vérification du corps de la réponse
	var responseBody struct {
		Projects []DAOs.Project `json:"projects"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	// Vérifie qu'on a au moins un projet
	assert.Greater(t, len(responseBody.Projects), 0)

	// Vérifie que doNotDeleteProject est présent
	var found bool
	for _, project := range responseBody.Projects {
		if project.Id == doNotDeleteProject.Id {
			found = true
			assert.Equal(t, doNotDeleteProject.UniqueId, project.UniqueId)
			assert.Equal(t, doNotDeleteProject.Name, project.Name)
			break
		}
	}
	assert.True(t, found, "Le projet de test n'a pas été trouvé dans la liste des projets")
}

func TestGetAllProjects_AsEmployee(t *testing.T) {
	w := sendRequest(router, "GET", "/projects", nil, nil, enums.Employee)
	assertResponse(t, w, http.StatusOK, nil)

	// Vérification du corps de la réponse
	var responseBody struct {
		Projects []DAOs.Project `json:"projects"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	// Vérifie qu'on a au moins un projet
	assert.Greater(t, len(responseBody.Projects), 0)

	// Vérifie que doNotDeleteProject est présent
	var found bool
	for _, project := range responseBody.Projects {
		if project.Id == doNotDeleteProject.Id {
			found = true
			assert.Equal(t, doNotDeleteProject.UniqueId, project.UniqueId)
			assert.Equal(t, doNotDeleteProject.Name, project.Name)
			break
		}
	}
	assert.True(t, found, "Le projet de test n'a pas été trouvé dans la liste des projets")
}
