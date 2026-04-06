package tests

import (
	"llio-api/database"
	"llio-api/models/DAOs"
	"net/http"
	"strconv"
	"testing"

	"llio-api/models/enums"

	"github.com/stretchr/testify/assert"
)

func TestDeleteCoManagerBadProjectId(t *testing.T) {
	w := sendRequest(router, "DELETE",
		"/project/999999/coManager/"+strconv.Itoa(doNotDeleteUser2.Id),
		nil,
		&doNotDeleteUser.Id,
		enums.ProjectManager,
	)

	assertResponse(t, w, http.StatusNotFound, nil)
	assert.NotNil(t, w.Body)
}
func TestDeleteCoManagerBadUserId(t *testing.T) {
	w := sendRequest(router, "DELETE",
		"/project/"+strconv.Itoa(doNotDeleteProject.Id)+"/coManager/999999",
		nil,
		&doNotDeleteUser.Id,
		enums.ProjectManager,
	)

	assertResponse(t, w, http.StatusNotFound, nil)
	assert.NotNil(t, w.Body)
}
func TestDeleteCoManagerAuthUserRoleTooLow(t *testing.T) {
	w := sendRequest(router, "DELETE",
		"/project/"+strconv.Itoa(doNotDeleteProject.Id)+"/coManager/"+strconv.Itoa(doNotDeleteUser2.Id),
		nil,
		&doNotDeleteUser.Id,
		enums.Employee,
	)

	assertResponse(t, w, http.StatusForbidden, nil)
	assert.NotNil(t, w.Body)
}
func TestDeleteCoManagerAuthUserNotInProject(t *testing.T) {
	w := sendRequest(router, "DELETE",
		"/project/"+strconv.Itoa(doNotDeleteProject.Id)+"/coManager/"+strconv.Itoa(doNotDeleteCoManager.UserId),
		nil,
		&doNotDeleteUser3.Id,
		enums.ProjectManager,
	)

	assertResponse(t, w, http.StatusForbidden, nil)
	assert.NotNil(t, w.Body)
}
func TestDeleteCoManagerUserNotCoManager(t *testing.T) {
	w := sendRequest(router, "DELETE",
		"/project/"+strconv.Itoa(doNotDeleteProject.Id)+"/coManager/"+strconv.Itoa(doNotDeleteUser3.Id),
		nil,
		&doNotDeleteUser.Id,
		enums.Administrator,
	)

	assertResponse(t, w, http.StatusNotFound, nil)
	assert.NotNil(t, w.Body)
}

func TestDeleteCoManager(t *testing.T) {

	var existing DAOs.CoManager
	err := database.DB.
		Where("user_id = ?", doNotDeleteCoManager.UserId).
		Where("project_id = ?", doNotDeleteCoManager.ProjectId).
		First(&existing).Error
	assert.NoError(t, err)

	w := sendRequest(router, "DELETE",
		"/project/"+strconv.Itoa(doNotDeleteCoManager.ProjectId)+"/coManager/"+strconv.Itoa(doNotDeleteCoManager.UserId),
		nil,
		&doNotDeleteUser.Id,
		enums.Administrator,
	)

	assertResponse(t, w, http.StatusOK, nil)
	assert.NotNil(t, w.Body)

	var deleted DAOs.CoManager
	errDB := database.DB.
		Where("user_id = ?", doNotDeleteCoManager.UserId).
		Where("project_id = ?", doNotDeleteCoManager.ProjectId).
		First(&deleted).Error

	assert.Error(t, errDB)
}
