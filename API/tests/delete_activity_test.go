package tests

import (
	"encoding/json"
	"errors"
	"llio-api/database"
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteActivity(t *testing.T) {

	activityToDelete := DTOs.ActivityDTO{
		Name:        "Activity to Delete",
		Description: "This activity will be deleted",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(24 * time.Hour),
		UserId:      doNotDeleteUser.Id,
		ProjectId:   doNotDeleteProject.Id,
		CategoryId:  doNotDeleteCategory.Id,
	}

	createW := sendRequest(router, "POST", "/activity", activityToDelete, doNotDeleteUser.Id)
	assertResponse(t, createW, http.StatusCreated, nil)

	var createResponseBody struct {
		Reponse  string        `json:"reponse"`
		Activity DAOs.Activity `json:"activity"`
	}
	err := json.Unmarshal(createW.Body.Bytes(), &createResponseBody)
	assert.NoError(t, err)

	idToDelete := strconv.Itoa(createResponseBody.Activity.Id)

	w := sendRequest(router, "DELETE", "/activity/"+idToDelete, nil, doNotDeleteUser.Id)
	assertResponse(t, w, http.StatusOK, nil)

	var deletedActivity DAOs.Activity
	errDB := database.DB.Where("id = ?", idToDelete).First(&deletedActivity).Error
	assert.Error(t, errDB)
	assert.True(t, errors.Is(errDB, gorm.ErrRecordNotFound))
}

func TestDeleteNonExistentActivity(t *testing.T) {
	nonExistentId := "99999"

	w := sendRequest(router, "DELETE", "/activity/"+nonExistentId, nil, doNotDeleteUser.Id)
	assertResponse(t, w, http.StatusNotFound, nil)
}
