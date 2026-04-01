package tests

import (
	"encoding/json"
	"errors"
	"llio-api/database"
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteCategory(t *testing.T) {

	// Création d'une catégorie
	category := DTOs.CategoryDTO{
		Name:        "Test Catégorie",
		Description: "Description de TestTest",
		ProjectId:   doNotDeleteProject.Id,
	}

	w := sendRequest(router, "POST", "/category", category, &doNotDeleteUser.Id)
	assertResponse(t, w, http.StatusCreated, nil)

	// Vérification du corps de la réponse
	var responseBody struct {
		Reponse  string        `json:"reponse"`
		Category DAOs.Category `json:"category"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "La catégorie a bien été ajoutée à la base de données.", responseBody.Reponse)
	assert.Equal(t, category.Name, responseBody.Category.Name)
	assert.Equal(t, category.Description, responseBody.Category.Description)

	idToDelete := strconv.Itoa(responseBody.Category.Id)

	//Suppression d'une catégorie
	g := sendRequest(router, "DELETE", "/category/"+idToDelete, nil, &doNotDeleteUser.Id)
	assertResponse(t, g, http.StatusOK, nil)

	var deletedCategory DAOs.Category
	errDB := database.DB.Where("id = ?", idToDelete).First(&deletedCategory).Error
	assert.Error(t, errDB)
	assert.True(t, errors.Is(errDB, gorm.ErrRecordNotFound))
}

func TestDeleteCategoryAlreadyLinked(t *testing.T) {

	// Création d'une catégorie
	category := DTOs.CategoryDTO{
		Name:        "Test Catégorie 2",
		Description: "Description de TestTest",
		ProjectId:   doNotDeleteProject.Id,
	}

	w := sendRequest(router, "POST", "/category", category, &doNotDeleteUser.Id)
	assertResponse(t, w, http.StatusCreated, nil)

	// Vérification du corps de la réponse
	var responseBody struct {
		Reponse  string        `json:"reponse"`
		Category DAOs.Category `json:"category"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "La catégorie a bien été ajoutée à la base de données.", responseBody.Reponse)
	assert.Equal(t, category.Name, responseBody.Category.Name)
	assert.Equal(t, category.Description, responseBody.Category.Description)

	idToDelete := strconv.Itoa(responseBody.Category.Id)

	// Création d'une activité
	activity := DTOs.ActivityDTO{
		Name:        "Test tâche",
		Description: "Test automatique de la création de tâche",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(24 * time.Hour),
		UserId:      doNotDeleteUser.Id,
		ProjectId:   doNotDeleteProject.Id,
		CategoryId:  responseBody.Category.Id,
	}

	w = sendRequest(router, "POST", "/activity", activity, nil)
	assertResponse(t, w, http.StatusCreated, nil)

	// Vérification du corps de la réponse
	var responseBody2 struct {
		Reponse  string                   `json:"reponse"`
		Activity DTOs.DetailedActivityDTO `json:"activity"`
	}
	err2 := json.Unmarshal(w.Body.Bytes(), &responseBody2)
	assert.NoError(t, err2)
	log.Println("Response Body:", string(w.Body.Bytes())) // Pour débogage
	log.Println("Response Body Struct:", responseBody2.Reponse)
	assert.Equal(t, "L'activité a bien été ajoutée à la base de données.", responseBody2.Reponse)
	assert.Equal(t, activity.Name, responseBody2.Activity.Name)
	assert.Equal(t, activity.Description, responseBody2.Activity.Description)
	assert.Equal(t, doNotDeleteProject.Name, responseBody2.Activity.ProjectName)

	// Vérification que l'activité est bien ajoutée dans la base de données
	var createdActivity DAOs.Activity
	errDB2 := database.DB.Where("name = ?", activity.Name).First(&createdActivity).Error
	assert.NoError(t, errDB2)
	assert.Equal(t, activity.Name, createdActivity.Name)

	//Supression de la catégorie
	g := sendRequest(router, "DELETE", "/category/"+idToDelete, nil, &doNotDeleteUser4.Id)
	assertResponse(t, g, http.StatusBadRequest, nil)

	var createdCategory DAOs.Category
	errDB := database.DB.Where("id = ?", idToDelete).First(&createdCategory).Error
	assert.NoError(t, errDB)
	assert.Equal(t, category.Name, createdCategory.Name)
}

func TestDeleteNonExistentCategory(t *testing.T) {
	nonExistentId := "99999"

	w := sendRequest(router, "DELETE", "/category/"+nonExistentId, nil, &doNotDeleteUser2.Id)
	assertResponse(t, w, http.StatusNotFound, nil)
}
