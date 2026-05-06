package tests

import (
	"net/http"
	"strconv"
	"testing"

	"llio-api/database"
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"llio-api/models/enums"

	"github.com/stretchr/testify/assert"
)

func TestArchiveAndDeArchiveProject(t *testing.T) {

	var dbProject1 DAOs.Project
	errDB1 := database.DB.Model(&DAOs.Project{}).Where("id = ?", doNotDeleteProject.Id).First(&dbProject1).Error
	assert.NoError(t, errDB1)

	assert.Equal(t, int(dbProject1.Status), 0)

	archiveProject := DTOs.ArchiveProjectDTO{
		Archived: true,
	}

	w := sendRequest(router, "POST", "/project/toggleArchive/"+strconv.Itoa(doNotDeleteProject.Id), archiveProject, &doNotDeleteUser.Id, enums.ProjectManager)
	assertResponse(t, w, http.StatusOK, nil)

	var dbProject2 DAOs.Project
	errDB2 := database.DB.Model(&DAOs.Project{}).Where("id = ?", doNotDeleteProject.Id).First(&dbProject2).Error
	assert.NoError(t, errDB2)

	assert.Equal(t, int(dbProject2.Status), 2)

	unarchiveProject := DTOs.ArchiveProjectDTO{
		Archived: false,
	}

	w2 := sendRequest(router, "POST", "/project/toggleArchive/"+strconv.Itoa(doNotDeleteProject.Id), unarchiveProject, &doNotDeleteUser.Id, enums.ProjectManager)
	assertResponse(t, w2, http.StatusOK, nil)

	var dbProject3 DAOs.Project
	errDB3 := database.DB.Model(&DAOs.Project{}).Where("id = ?", doNotDeleteProject.Id).First(&dbProject3).Error
	assert.NoError(t, errDB3)

	assert.Equal(t, int(dbProject3.Status), 0)
}

func TestArchiveNoPermissionUser(t *testing.T) {
	var dbProject1 DAOs.Project
	errDB1 := database.DB.Model(&DAOs.Project{}).Where("id = ?", doNotDeleteProject.Id).First(&dbProject1).Error
	assert.NoError(t, errDB1)

	assert.Equal(t, int(dbProject1.Status), 0)

	w := sendRequest(router, "POST", "/project/toggleArchive/"+strconv.Itoa(doNotDeleteProject.Id), "{\"Archived\":true}", &doNotDeleteUser2.Id, enums.ProjectManager)
	assertResponse(t, w, http.StatusBadRequest, nil)

	var dbProject2 DAOs.Project
	errDB2 := database.DB.Model(&DAOs.Project{}).Where("id = ?", doNotDeleteProject.Id).First(&dbProject2).Error
	assert.NoError(t, errDB2)

	assert.Equal(t, int(dbProject2.Status), 0)
}
