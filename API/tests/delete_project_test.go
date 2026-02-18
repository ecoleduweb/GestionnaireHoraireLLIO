package tests

import (
	"net/http"
	"strconv"
	"testing"

	"llio-api/models/enums"

	"github.com/stretchr/testify/assert"
)

func TestDeleteProjectBadId(t *testing.T) {
	w := sendRequest(router, "DELETE", "/project/spin", nil, doNotDeleteUser.Id, enums.Administrator)
	assertResponse(t, w, http.StatusNotFound, nil)
	assert.NotNil(t, w.Body)
}

func TestDeleteProjectNoId(t *testing.T) {
	w := sendRequest(router, "DELETE", "/project/", nil, doNotDeleteUser.Id, enums.Administrator)
	assertResponse(t, w, http.StatusNotFound, nil)
	assert.NotNil(t, w.Body)
}

func TestDeleteProjectNotEmpty(t *testing.T) {
	w := sendRequest(router, "DELETE", "/project/"+strconv.Itoa(doNotDeleteProject.Id), nil, doNotDeleteUser.Id, enums.Administrator)
	assertResponse(t, w, http.StatusForbidden, nil)
	assert.NotNil(t, w.Body)
}

func TestDeleteProjectWithoutPermission(t *testing.T) {
	w := sendRequest(router, "DELETE", "/project/"+strconv.Itoa(pleaseDeleteProject.Id), nil, doNotDeleteUser.Id, enums.Employee)
	assertResponse(t, w, http.StatusForbidden, nil)
	assert.NotNil(t, w.Body)
}

func TestDeleteProject(t *testing.T) {
	w := sendRequest(router, "DELETE", "/project/"+strconv.Itoa(pleaseDeleteProject.Id), nil, doNotDeleteUser.Id, enums.Administrator)
	assertResponse(t, w, http.StatusOK, nil)
	assert.NotNil(t, w.Body)
}
