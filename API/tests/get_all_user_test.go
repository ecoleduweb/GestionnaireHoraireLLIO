package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"llio-api/models/enums"

	"github.com/stretchr/testify/assert"
)

func TestGetAllUsers(t *testing.T) {
	w := sendRequest(router, "GET", "/users", nil, nil, enums.ProjectManager)
	assertResponse(t, w, http.StatusOK, nil)
	assert.NotNil(t, w.Body)
}

func TestGetAllUsersBySingleRole(t *testing.T) {
	w := sendRequest(router, "GET", "/users?role=1", nil, nil, enums.ProjectManager)
	assertResponse(t, w, http.StatusOK, nil)
	assert.NotNil(t, w.Body)
}

func TestGetAllUsersByMultipleRoles(t *testing.T) {
	w := sendRequest(router, "GET", "/users?role=1&role=2", nil, nil, enums.ProjectManager)
	assertResponse(t, w, http.StatusOK, nil)
	assert.NotNil(t, w.Body)
}

func TestGetAllUsersWithInvalidToken(t *testing.T) {
	// Create a request without using the sendRequest helper to avoid token generation
	req, _ := http.NewRequest("GET", "/users", nil)
	req.Header.Set("Content-Type", "application/json")

	// Add an invalid token cookie
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    "invalid_token",
		HttpOnly: true,
		Path:     "/",
	}
	req.AddCookie(cookie)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assertResponse(t, w, http.StatusUnauthorized, nil)
}

func TestGetAllUsersWithoutPermission(t *testing.T) {
	w := sendRequest(router, "GET", "/users", nil, nil, enums.Employee)
	assertResponse(t, w, http.StatusForbidden, nil)
}
