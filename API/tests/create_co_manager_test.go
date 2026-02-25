package tests

import (
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
}
