package tests

import (
	"fmt"
	"net/http"
	"testing"

	"llio-api/database"
	"llio-api/models/DAOs"
	"llio-api/models/enums"
)

func TestUpdateUserRoleAsProjectManager(t *testing.T) {
	testUser := DAOs.User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test.user@example.com",
		Role:      enums.Employee,
	}
	database.DB.Create(&testUser)

	newRole := enums.ProjectManager
	roleRequest := map[string]interface{}{
		"role": &newRole,
	}

	w := sendRequest(router, "PATCH", fmt.Sprintf("/user/%d/role", testUser.Id), roleRequest, doNotDeleteUser.Id, enums.ProjectManager)

	assertResponse(t, w, http.StatusForbidden, nil)
}

func TestUpdateUserRoleAsEmployee(t *testing.T) {
	testUser := DAOs.User{
		FirstName: "Test",
		LastName:  "User2",
		Email:     "test.user2@example.com",
		Role:      enums.Employee,
	}
	database.DB.Create(&testUser)

	newRole := int(enums.ProjectManager)
	roleRequest := map[string]interface{}{
		"role": &newRole,
	}

	w := sendRequest(router, "PATCH", fmt.Sprintf("/user/%d/role", testUser.Id), roleRequest, doNotDeleteUser.Id, enums.Employee)

	assertResponse(t, w, http.StatusForbidden, nil)
}

func TestUpdateOwnRole(t *testing.T) {
	newRole := enums.Employee
	roleRequest := map[string]interface{}{
		"role": &newRole,
	}

	w := sendRequest(router, "PATCH", fmt.Sprintf("/user/%d/role", doNotDeleteUser.Id), roleRequest, doNotDeleteUser.Id, enums.Administrator)

	assertResponse(t, w, http.StatusForbidden, nil)
}

func TestUpdateUserRoleWithInvalidID(t *testing.T) {
	newRole := int(enums.ProjectManager)
	roleRequest := map[string]interface{}{
		"role": &newRole,
	}

	w := sendRequest(router, "PATCH", "/user/invalid/role", roleRequest, doNotDeleteUser.Id, enums.Administrator)

	assertResponse(t, w, http.StatusBadRequest, nil)
}

func TestUpdateUserRoleWithNonExistentUser(t *testing.T) {
	// Create request body
	newRole := int(enums.ProjectManager)
	roleRequest := map[string]interface{}{
		"role": &newRole,
	}

	// Send request with non-existent user ID
	nonExistentID := 99999
	w := sendRequest(router, "PATCH", fmt.Sprintf("/user/%d/role", nonExistentID), roleRequest, doNotDeleteUser.Id, enums.Administrator)

	// Assert response
	assertResponse(t, w, http.StatusNotFound, nil)
}

func TestUpdateUserRoleWithInvalidRole(t *testing.T) {
	// Create a test user to update
	testUser := DAOs.User{
		FirstName: "Test",
		LastName:  "User3",
		Email:     "test.user3@example.com",
		Role:      enums.Employee,
	}
	database.DB.Create(&testUser)

	// Create request body with invalid role
	invalidRole := 999
	roleRequest := map[string]interface{}{
		"role": &invalidRole,
	}

	// Send request
	w := sendRequest(router, "PATCH", fmt.Sprintf("/user/%d/role", testUser.Id), roleRequest, doNotDeleteUser.Id, enums.Administrator)

	// Assert response
	assertResponse(t, w, http.StatusBadRequest, nil)
}

func TestUpdateUserRoleMissingRoleField(t *testing.T) {
	// Create a test user to update
	testUser := DAOs.User{
		FirstName: "Test",
		LastName:  "User4",
		Email:     "test.user4@example.com",
		Role:      enums.Employee,
	}
	database.DB.Create(&testUser)

	// Create request body with missing role field
	emptyRequest := map[string]interface{}{}

	// Send request
	w := sendRequest(router, "PATCH", fmt.Sprintf("/user/%d/role", testUser.Id), emptyRequest, doNotDeleteUser.Id, enums.Administrator)

	// Assert response
	assertResponse(t, w, http.StatusBadRequest, nil)
}
