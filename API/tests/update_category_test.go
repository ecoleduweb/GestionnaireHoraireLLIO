package tests

import (
	"encoding/json"
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateCategory(t *testing.T) {
	// Création d'une catégorie à modifier
	category := DTOs.CategoryDTO{
		Id:          99999,
		Name:        "Catégorie à modifier",
		Description: "Description de test",
		ProjectId:   doNotDeleteProject.Id,
	}

	w := sendRequest(router, "POST", "/category", category, nil)
	assertResponse(t, w, http.StatusCreated, nil)

	// Vérification du corps de la réponse
	var responseBody struct {
		Reponse  string        `json:"reponse"`
		Category DAOs.Category `json:"activity"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	// Création d'une catégorie à modifier
	categoryToUpdate := DTOs.CategoryDTO{
		Id:          category.Id,
		Name:        "Catégorie à modifier",
		Description: "Description de test",
	}

	w = sendRequest(router, "PUT", "/category", categoryToUpdate, nil)
	assertResponse(t, w, http.StatusOK, nil)

	var updateResponseBody struct {
		UpdatedCategory DTOs.CategoryDTO `json:"updatedCategory"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &updateResponseBody)
	assert.NoError(t, err)
	assert.Equal(t, categoryToUpdate.Name, updateResponseBody.UpdatedCategory.Name)
	assert.Equal(t, categoryToUpdate.Description, updateResponseBody.UpdatedCategory.Description)
}

func TestUpdateCategoryWithInvalidId(t *testing.T) {
	category := DTOs.CategoryDTO{
		Id:          9999901,
		Name:        "Catégorie à modifier",
		Description: "Description de test",
		ProjectId:   doNotDeleteProject.Id,
	}

	w := sendRequest(router, "PUT", "/category", category, nil)
	assertResponse(t, w, http.StatusNotFound, nil)
}
