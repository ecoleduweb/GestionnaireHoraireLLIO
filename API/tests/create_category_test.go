package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"llio-api/models/DAOs"
	"llio-api/models/DTOs"

	"github.com/stretchr/testify/assert"
)

func TestCreateCategory(t *testing.T) {
	// Création d'une catégorie
	category := DTOs.CategoryDTO{
		Name:        "Nouvelle Catégorie",
		Description: "Description de test",
		ProjectId:   doNotDeleteProject.Id,
	}

	w := sendRequest(router, "POST", "/category", category, nil)
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
}

func TestDoNotCreateCategoryWithSameNameInSameProject(t *testing.T) {
	// Création d'une première catégorie
	initialCategory := DTOs.CategoryDTO{
		Name:        "Catégorie Unique",
		Description: "Description de test",
		ProjectId:   doNotDeleteProject.Id,
	}

	w := sendRequest(router, "POST", "/category", initialCategory, nil)
	assertResponse(t, w, http.StatusCreated, nil)

	// Tentative de création d'une catégorie avec le même nom dans le même projet
	duplicateCategory := DTOs.CategoryDTO{
		Name:        "Catégorie Unique",
		Description: "Autre description",
		ProjectId:   doNotDeleteProject.Id,
	}

	w = sendRequest(router, "POST", "/category", duplicateCategory, nil)
	assertResponse(t, w, http.StatusConflict, nil)

}

func TestCreateCategoryWithSameNameDifferentProjects(t *testing.T) {
	// Création d'une catégorie dans le premier projet
	firstCategory := DTOs.CategoryDTO{
		Name:        "Catégorie Multiprojet",
		Description: "Description de test",
		ProjectId:   doNotDeleteProject.Id,
	}

	w = sendRequest(router, "POST", "/category", firstCategory, nil)
	assertResponse(t, w, http.StatusCreated, nil)

	// Création d'une catégorie avec le même nom dans le nouveau projet
	secondCategory := DTOs.CategoryDTO{
		Name:        "Catégorie Multiprojet",
		Description: "Autre description",
		ProjectId:   doNotDeleteProject2.Id,
	}

	w = sendRequest(router, "POST", "/category", secondCategory, nil)
	assertResponse(t, w, http.StatusCreated, nil)

	// Vérification du corps de la réponse
	var responseBody struct {
		Reponse  string        `json:"reponse"`
		Category DAOs.Category `json:"category"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "La catégorie a bien été ajoutée à la base de données.", responseBody.Reponse)
	assert.Equal(t, secondCategory.Name, responseBody.Category.Name)
	assert.Equal(t, secondCategory.Description, responseBody.Category.Description)
}

func TestDoNotCreateCategoryWithInvalidName(t *testing.T) {
	// Création d'une catégorie avec un nom trop long
	category := DTOs.CategoryDTO{
		Name:        "Ceci est un très long nom de catégorie qui dépasse largement la limite de 50 caractères et qui devrait donc être rejeté par la validation",
		Description: "Description de test",
		ProjectId:   doNotDeleteProject.Id,
	}

	w := sendRequest(router, "POST", "/category", category, nil)
	expectedErrors := []DTOs.FieldErrorDTO{
		{Field: "name", Message: "Le champ name est invalide ou manquant"},
	}
	assertResponse(t, w, http.StatusBadRequest, expectedErrors)
}

func TestDoNotCreateCategoryWithoutName(t *testing.T) {
	// Création d'une catégorie sans nom
	category := DTOs.CategoryDTO{
		Name:        "",
		Description: "Description de test",
		ProjectId:   doNotDeleteProject.Id,
	}

	w := sendRequest(router, "POST", "/category", category, nil)
	expectedErrors := []DTOs.FieldErrorDTO{
		{Field: "name", Message: "Le champ name est invalide ou manquant"},
	}
	assertResponse(t, w, http.StatusBadRequest, expectedErrors)
}

func TestDoNotCreateCategoryWithoutDescription(t *testing.T) {
	// Création d'une catégorie sans nom
	category := DTOs.CategoryDTO{
		Name:        "Test",
		Description: "",
		ProjectId:   doNotDeleteProject.Id,
	}

	w := sendRequest(router, "POST", "/category", category, nil)
	expectedErrors := []DTOs.FieldErrorDTO{
		{Field: "description", Message: "Le champ description est invalide ou manquant"},
	}
	assertResponse(t, w, http.StatusBadRequest, expectedErrors)
}

func TestDoNotCreateCategoryWithoutProject(t *testing.T) {
	// Création d'une catégorie sans nom
	category := DTOs.CategoryDTO{
		Name:        "Test",
		Description: "Description",
	}

	w := sendRequest(router, "POST", "/category", category, nil)
	expectedErrors := []DTOs.FieldErrorDTO{
		{Field: "projectId", Message: "Le champ projectId est invalide ou manquant"},
	}
	assertResponse(t, w, http.StatusBadRequest, expectedErrors)
}

// Pour plus tard quand les projets seront en place
// func TestDoNotCreateCategoryWithNonExestingProject(t *testing.T) {
// 	// Création d'une catégorie sans nom
// 	category := DTOs.CategoryDTO{
// 		Name:        "Test",
// 		Description: "Description",
// 		ProjectId:   doNotDeleteProject.Id,
// 	}

// 	w := sendRequest(router, "POST", "/category", category)
// 	assertResponse(t, w, http.StatusBadRequest, nil)
// }
