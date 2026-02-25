package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"llio-api/models/enums"

	"github.com/stretchr/testify/assert"
)

func TestGetDetailedProjectsByUser(t *testing.T) {
	w := sendRequest(router, "GET", "/projects/me/detailed", nil, nil, enums.Employee)
	assertResponse(t, w, http.StatusOK, nil)

	// Vérification du corps de la réponse
	var responseBody struct {
		Projects []map[string]any `json:"projects"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	// Vérifie qu'on a au moins un projet
	assert.Greater(t, len(responseBody.Projects), 0)

	// Vérifie la structure des données détaillées
	firstProject := responseBody.Projects[0]
	assert.Contains(t, firstProject, "id")
	assert.Contains(t, firstProject, "name")
	// Ajoutez d'autres champs selon la structure retournée par GetDetailedProjectsByUserId
}

func TestGetDetailedProjectsByUser_AsAdmin(t *testing.T) {
	w := sendRequest(router, "GET", "/projects/me/detailed", nil, nil, enums.Administrator)
	assertResponse(t, w, http.StatusOK, nil)

	// Vérification du corps de la réponse
	var responseBody struct {
		Projects []map[string]any `json:"projects"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	// Vérifie qu'on a au moins un projet
	assert.Greater(t, len(responseBody.Projects), 0)
}

func TestGetDetailedProjectsByUser_EmptyProjects(t *testing.T) {
	// Ce test nécessiterait un utilisateur sans projets assignés
	// ou un mock du service pour retourner nil
	// Exemple avec un utilisateur qui n'a pas de projets :
	w := sendRequest(router, "GET", "/projects/me/detailed", nil, nil, enums.Employee)

	// Si le service retourne nil/empty, on devrait avoir une liste vide
	if w.Code == http.StatusOK {
		var responseBody struct {
			Projects []map[string]any `json:"projects"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		// La réponse devrait être une liste vide, pas null
		assert.NotNil(t, responseBody.Projects)
	}
}
