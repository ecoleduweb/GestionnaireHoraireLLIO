package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"llio-api/models/DAOs"
	"llio-api/models/DTOs"

	"github.com/stretchr/testify/assert"
)

func TestUpdateActivity(t *testing.T) {
	// activité à modifier
	initialActivity := DTOs.ActivityDTO{
		Name:        "Original Activity",
		Description: "Original Description",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(24 * time.Hour),
		UserId:      doNotDeleteUser.Id,
		ProjectId:   doNotDeleteProject.Id,
		CategoryId:  doNotDeleteCategory.Id,
	}

	// Création de l'acitivité
	createW := sendRequest(router, "POST", "/activity", initialActivity, nil)
	assertResponse(t, createW, http.StatusCreated, nil)

	var createResponseBody struct {
		Reponse  string        `json:"reponse"`
		Activity DAOs.Activity `json:"activity"`
	}
	err := json.Unmarshal(createW.Body.Bytes(), &createResponseBody)
	assert.NoError(t, err)

	// Activité après maj
	updateActivity := DTOs.ActivityDTO{
		Id:          createResponseBody.Activity.Id,
		Name:        "Updated Activity",
		Description: "Updated Description",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(48 * time.Hour),
		UserId:      doNotDeleteUser.Id,
		ProjectId:   doNotDeleteProject.Id,
		CategoryId:  doNotDeleteCategory.Id,
	}

	w := sendRequest(router, "PUT", "/activity", updateActivity, nil)
	assertResponse(t, w, http.StatusOK, nil)

	var updateResponseBody struct {
		UpdatedActivity DTOs.DetailedActivityDTO `json:"updated_activity"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &updateResponseBody)
	assert.NoError(t, err)
	assert.Equal(t, updateActivity.Name, updateResponseBody.UpdatedActivity.Name)
	assert.Equal(t, updateActivity.Description, updateResponseBody.UpdatedActivity.Description)
	assert.Equal(t, doNotDeleteProject.Name, updateResponseBody.UpdatedActivity.ProjectName)
}

func TestUpdateActivityWithInvalidId(t *testing.T) {
	updateActivity := DTOs.ActivityDTO{
		Id:          99999,
		Name:        "Updated Activity",
		Description: "Updated Description",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(48 * time.Hour),
		UserId:      doNotDeleteUser.Id,
		ProjectId:   doNotDeleteProject.Id,
		CategoryId:  doNotDeleteCategory.Id,
	}

	w := sendRequest(router, "PUT", "/activity", updateActivity, nil)
	assertResponse(t, w, http.StatusNotFound, nil)
}
