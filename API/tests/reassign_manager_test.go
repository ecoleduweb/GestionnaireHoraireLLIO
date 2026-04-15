package tests

import (
	"encoding/json"
	"llio-api/database"
	"llio-api/models/DAOs"
	"log"
	"net/http"
	"strconv"
	"testing"

	"llio-api/models/enums"

	"github.com/stretchr/testify/assert"
)

func TestReassignManager(t *testing.T) {
	w := sendRequest(
		router,
		"PUT",
		"/project/"+strconv.Itoa(doNotDeleteProject.Id)+"/reassignManager/"+strconv.Itoa(doNotDeleteUser2.Id),
		nil,
		&doNotDeleteUser.Id,
		enums.Administrator,
	)

	assertResponse(t, w, http.StatusOK, nil)
	assert.NotNil(t, w.Body)

	var responseBody struct {
		Response string `json:"response"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	log.Println("Response:", responseBody.Response)

	assert.Equal(t, "Le manager a bien été réattribué", responseBody.Response)

	// Vérification DB
	var project DAOs.Project
	errDB := database.DB.First(&project, doNotDeleteProject.Id).Error
	assert.NoError(t, errDB)
	assert.Equal(t, doNotDeleteUser2.Id, project.ManagerId)
}