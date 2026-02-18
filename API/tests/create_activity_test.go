package tests

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"

	"llio-api/database"
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"

	"github.com/stretchr/testify/assert"
)

func TestCreateActivity(t *testing.T) {

	activity := DTOs.ActivityDTO{
		Name:        "Test tâche",
		Description: "Test automatique de la création de tâche",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(24 * time.Hour),
		UserId:      doNotDeleteUser.Id,
		ProjectId:   doNotDeleteProject.Id,
		CategoryId:  doNotDeleteCategory.Id,
	}

	w = sendRequest(router, "POST", "/activity", activity, doNotDeleteUser.Id)
	assertResponse(t, w, http.StatusCreated, nil)

	// Vérification du corps de la réponse
	var responseBody struct {
		Reponse  string                   `json:"reponse"`
		Activity DTOs.DetailedActivityDTO `json:"activity"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	log.Println("Response Body:", string(w.Body.Bytes())) // Pour débogage
	log.Println("Response Body Struct:", responseBody.Reponse)
	assert.Equal(t, "L'activité a bien été ajoutée à la base de données.", responseBody.Reponse)
	assert.Equal(t, activity.Name, responseBody.Activity.Name)
	assert.Equal(t, activity.Description, responseBody.Activity.Description)
	assert.Equal(t, doNotDeleteProject.Name, responseBody.Activity.ProjectName)

	// Vérification que l'activité est bien ajoutée dans la base de données
	var createdActivity DAOs.Activity
	errDB := database.DB.Where("name = ?", activity.Name).First(&createdActivity).Error
	assert.NoError(t, errDB)
	assert.Equal(t, activity.Name, createdActivity.Name)
}

func TestDoNotCreateActivityWithEndDateBeforeStartDate(t *testing.T) {

	// Création d'une activités à envoyer dans la requête
	activity := DTOs.ActivityDTO{
		Name:        "Test tâche",
		Description: "Test automatique de la création de tâche",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(-24 * time.Hour),
		UserId:      doNotDeleteUser.Id,
		ProjectId:   doNotDeleteProject.Id,
		CategoryId:  doNotDeleteCategory.Id,
	}

	w = sendRequest(router, "POST", "/activity", activity, doNotDeleteUser.Id)

	expectedErrors := []DTOs.FieldErrorDTO{
		{Field: "startDate", Message: "La date de début doit être avant la date de fin"},
	}
	assertResponse(t, w, http.StatusBadRequest, expectedErrors)
}

func TestCreateActivityWithoutNameAndDescription(t *testing.T) {

	activity := DTOs.ActivityDTO{
		Name:        "",
		Description: "",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(24 * time.Hour),
		UserId:      doNotDeleteUser.Id,
		ProjectId:   doNotDeleteProject.Id,
		CategoryId:  doNotDeleteCategory.Id,
	}

	w = sendRequest(router, "POST", "/activity", activity, doNotDeleteUser.Id)

	assertResponse(t, w, http.StatusCreated, nil)
}

func TestDoNotCreateActivityWithLenghtNameOver50(t *testing.T) {

	activity := DTOs.ActivityDTO{
		//Sitation dans le film Astérix et Obélix : mission Cléopâtre
		Name:        "Vous savez, moi je ne crois pas qu’il y ait de bonne ou de mauvaise situation. Moi, si je devais résumer ma vie aujourd’hui avec vous, je dirais que c’est d’abord des rencontres...",
		Description: "Test automatique de la création de tâche",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(24 * time.Hour),
		UserId:      doNotDeleteUser.Id,
		ProjectId:   doNotDeleteProject.Id,
		CategoryId:  doNotDeleteCategory.Id,
	}
	w = sendRequest(router, "POST", "/activity", activity, doNotDeleteUser.Id)

	expectedErrors := []DTOs.FieldErrorDTO{
		{Field: "name", Message: "Le champ name est invalide ou manquant"},
	}
	assertResponse(t, w, http.StatusBadRequest, expectedErrors)
}

func TestDoNotCreateActivityWithoutDates(t *testing.T) {

	activity := DTOs.ActivityDTO{
		Name:        "Test tâche",
		Description: "Test automatique de la création de tâche",
		UserId:      doNotDeleteUser.Id,
		ProjectId:   doNotDeleteProject.Id,
		CategoryId:  doNotDeleteCategory.Id,
	}

	w = sendRequest(router, "POST", "/activity", activity, doNotDeleteUser.Id)

	expectedErrors := []DTOs.FieldErrorDTO{
		{Field: "startDate", Message: "Le champ startDate est invalide ou manquant"},
		{Field: "endDate", Message: "Le champ endDate est invalide ou manquant"},
	}
	assertResponse(t, w, http.StatusBadRequest, expectedErrors)
}

func TestDoNotCreateActivityWithInvalidStartDate(t *testing.T) {

	activity := DTOs.ActivityDTO{
		Name:        "Test tâche",
		Description: "Test automatique de la création de tâche",
		StartDate:   time.Time{},
		EndDate:     time.Now().Add(24 * time.Hour),
		UserId:      doNotDeleteUser.Id,
		ProjectId:   doNotDeleteProject.Id,
		CategoryId:  doNotDeleteCategory.Id,
	}

	w = sendRequest(router, "POST", "/activity", activity, doNotDeleteUser.Id)

	expectedErrors := []DTOs.FieldErrorDTO{
		{Field: "startDate", Message: "Le champ startDate est invalide ou manquant"},
	}
	assertResponse(t, w, http.StatusBadRequest, expectedErrors)
}

func TestDoNotCreateActivityWithInvalidEndDate(t *testing.T) {

	activity := DTOs.ActivityDTO{
		Name:        "Test tâche",
		Description: "Test automatique de la création de tâche",
		StartDate:   time.Now(),
		EndDate:     time.Time{},
		UserId:      doNotDeleteUser.Id,
		ProjectId:   doNotDeleteProject.Id,
		CategoryId:  doNotDeleteCategory.Id,
	}

	w = sendRequest(router, "POST", "/activity", activity, doNotDeleteUser.Id)

	expectedErrors := []DTOs.FieldErrorDTO{
		{Field: "endDate", Message: "Le champ endDate est invalide ou manquant"},
	}
	assertResponse(t, w, http.StatusBadRequest, expectedErrors)
}
