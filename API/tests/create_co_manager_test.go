package tests

import (
	"encoding/json"
	"llio-api/database"
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"log"
	"net/http"
	"strconv"
	"testing"

	"llio-api/models/enums"

	"github.com/stretchr/testify/assert"
)

func TestCreateCoManagerBadProjectId(t *testing.T) {
	w := sendRequest(router, "POST", "/project/999999/coManager/"+strconv.Itoa(doNotDeleteUser2.Id), nil, &doNotDeleteUser.Id, enums.ProjectManager)
	assertResponse(t, w, http.StatusNotFound, nil)
	assert.NotNil(t, w.Body)
}

func TestCreateCoManagerBadUserId(t *testing.T) {
	w := sendRequest(router, "POST", "/project/"+strconv.Itoa(doNotDeleteProject.Id)+"/coManager/999999", nil, &doNotDeleteUser.Id, enums.ProjectManager)
	assertResponse(t, w, http.StatusNotFound, nil)
	assert.NotNil(t, w.Body)
}

func TestCreateCoManagerAuthUserRoleTooLow(t *testing.T) {
	w := sendRequest(router, "POST", "/project/"+strconv.Itoa(doNotDeleteProject.Id)+"/coManager/"+strconv.Itoa(doNotDeleteUser2.Id), nil, &doNotDeleteUser.Id, enums.Employee)
	assertResponse(t, w, http.StatusForbidden, nil)
	assert.NotNil(t, w.Body)
}

func TestCreateCoManagerAuthUserProjectManagerButNotOnProjectRequested(t *testing.T) {
	w := sendRequest(router, "POST", "/project/"+strconv.Itoa(doNotDeleteProject.Id)+"/coManager/"+strconv.Itoa(doNotDeleteUser2.Id), nil, &doNotDeleteUser2.Id, enums.ProjectManager)
	assertResponse(t, w, http.StatusForbidden, nil)
	assert.NotNil(t, w.Body)
}

func TestCreateCoManagerUserAlreadyManager(t *testing.T) {
	w := sendRequest(router, "POST", "/project/"+strconv.Itoa(doNotDeleteProject.Id)+"/coManager/"+strconv.Itoa(doNotDeleteUser.Id), nil, &doNotDeleteUser.Id, enums.Administrator)
	assertResponse(t, w, http.StatusBadRequest, nil)
	assert.NotNil(t, w.Body)
}

func TestCreateCoManager(t *testing.T) {
	w := sendRequest(router, "POST", "/project/"+strconv.Itoa(doNotDeleteProject.Id)+"/coManager/"+strconv.Itoa(doNotDeleteUser2.Id), nil, &doNotDeleteUser.Id, enums.Administrator)
	assertResponse(t, w, http.StatusCreated, nil)
	assert.NotNil(t, w.Body)

	// Vérification du corps de la réponse
	var responseBody struct {
		Reponse   string            `json:"response"`
		CoManager DTOs.CoManagerDTO `json:"coManager"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	log.Println("Response Body:", string(w.Body.Bytes())) // Pour débogage
	log.Println("Response Body Struct:", responseBody.Reponse)
	assert.Equal(t, "Le co-chargé de projet a bien été ajouté à la base de données", responseBody.Reponse)
	assert.Equal(t, doNotDeleteUser2.Id, responseBody.CoManager.UserId)
	assert.Equal(t, doNotDeleteProject.Id, responseBody.CoManager.ProjectId)

	// Vérification que l'activité est bien ajoutée dans la base de données
	var createdActivity DAOs.CoManager
	errDB := database.DB.Where("user_id = ?", doNotDeleteUser2.Id).Where("project_id = ?", doNotDeleteProject.Id).First(&createdActivity).Error
	assert.NoError(t, errDB)
	assert.Equal(t, doNotDeleteUser2.Id, createdActivity.UserId)
	assert.Equal(t, doNotDeleteProject.Id, createdActivity.ProjectId)
}

func TestCreateCoManagerIsAlreadyCoManager(t *testing.T) {
	w := sendRequest(router, "POST", "/project/"+strconv.Itoa(doNotDeleteCoManager.ProjectId)+"/coManager/"+strconv.Itoa(doNotDeleteCoManager.UserId), nil, &doNotDeleteUser.Id, enums.Administrator)
	assertResponse(t, w, http.StatusBadRequest, nil)
	assert.NotNil(t, w.Body)
}
