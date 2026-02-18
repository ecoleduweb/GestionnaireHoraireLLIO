package tests

import (
	"net/http"
	"strconv"
	"testing"

	"llio-api/models/enums"

	"github.com/stretchr/testify/assert"
)

func TestDeleteUserBadId(t *testing.T) {
	w := sendRequest(router, "DELETE", "/users/spin", nil, doNotDeleteUser.Id, enums.Administrator)
	assertResponse(t, w, http.StatusNotFound, nil)
	assert.NotNil(t, w.Body)
}

func TestDeleteUserNoId(t *testing.T) {
	w := sendRequest(router, "DELETE", "/users/", nil, doNotDeleteUser.Id, enums.Administrator)
	assertResponse(t, w, http.StatusNotFound, nil)
	assert.NotNil(t, w.Body)
}

func TestDeleteUserNotEmpty(t *testing.T) {
	w := sendRequest(router, "DELETE", "/user/"+strconv.Itoa(doNotDeleteUser.Id), nil, doNotDeleteUser.Id, enums.Administrator)
	assertResponse(t, w, http.StatusForbidden, nil)
	assert.NotNil(t, w.Body)
}

func TestDeleteUserWithoutPermission(t *testing.T) {
	w := sendRequest(router, "DELETE", "/user/"+strconv.Itoa(pleaseDeleteUser.Id), nil, doNotDeleteUser.Id, enums.Employee)
	assertResponse(t, w, http.StatusForbidden, nil)
	assert.NotNil(t, w.Body)
}

func TestDeleteUser(t *testing.T) {
	w := sendRequest(router, "DELETE", "/user/"+strconv.Itoa(pleaseDeleteUser.Id), nil, doNotDeleteUser.Id, enums.Administrator)
	assertResponse(t, w, http.StatusOK, nil)
	assert.NotNil(t, w.Body)
}
