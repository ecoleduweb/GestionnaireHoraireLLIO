package tests

import (
	"encoding/json"
	"llio-api/database"
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateTimeBank_Surplus(t *testing.T) {
	// 1. Préparation des dates (Lundi de la semaine passée)
	now := time.Now()
	offset := int(now.Weekday())
	if offset == 0 {
		offset = 7
	} // Dimanche = 7
	// Lundi de la semaine courante
	mondayCurrentWeek := now.AddDate(0, 0, -(offset - 1))
	// Lundi de la semaine D'AVANT (début de la période de calcul)
	mondayLastWeek := mondayCurrentWeek.AddDate(0, 0, -7)
	startDateStr := mondayLastWeek.Format("2006-01-02")

	// 2. Création manuelle des activités en BD
	// On crée 3 activités de 4h = 12h totales pour la semaine passée
	// On les lie à 'doNotDeleteUser' qui est l'utilisateur utilisé par défaut par sendRequest
	activities := []DAOs.Activity{
		{
			Name:        "Test TimeBank 1",
			Description: "Test Auto",
			StartDate:   createDate(mondayLastWeek, 9, 0),
			EndDate:     createDate(mondayLastWeek, 13, 0), // 4 heures
			UserId:      doNotDeleteUser.Id,
			ProjectId:   doNotDeleteProject.Id,
			CategoryId:  doNotDeleteCategory.Id,
		},
		{
			Name:        "Test TimeBank 2",
			Description: "Test Auto",
			StartDate:   createDate(mondayLastWeek.AddDate(0, 0, 1), 9, 0), // Mardi
			EndDate:     createDate(mondayLastWeek.AddDate(0, 0, 1), 13, 0),
			UserId:      doNotDeleteUser.Id,
			ProjectId:   doNotDeleteProject.Id,
			CategoryId:  doNotDeleteCategory.Id,
		},
		{
			Name:        "Test TimeBank 3",
			Description: "Test Auto",
			StartDate:   createDate(mondayLastWeek.AddDate(0, 0, 2), 9, 0), // Mercredi
			EndDate:     createDate(mondayLastWeek.AddDate(0, 0, 2), 13, 0),
			UserId:      doNotDeleteUser.Id,
			ProjectId:   doNotDeleteProject.Id,
			CategoryId:  doNotDeleteCategory.Id,
		},
	}

	// Insertion en base
	for _, act := range activities {
		err := database.DB.Create(&act).Error
		if err != nil {
			t.Fatalf("Erreur création activité de test: %v", err)
		}
	}

	// Nettoyage des données après le test
	defer func() {
		database.DB.Where("name LIKE ?", "Test TimeBank%").Delete(&DAOs.Activity{})
	}()

	// 3. Préparation de la requête
	// Scénario : 10h attendues vs 12h faites = Solde de +2
	payload := map[string]interface{}{
		"startDate":          startDateStr,
		"hoursPerWeek":       10,
		"offset":             0,
		"excludeCurrentWeek": true,
	}

	// 4. Envoi de la requête via le helper global 'sendRequest'
	// Note: Assure-toi que la route correspond bien à celle définie dans routes.go
	w := sendRequest(router, "POST", "/user/calculate", payload)

	// 5. Vérifications
	assertResponse(t, w, http.StatusOK, nil)

	var response DTOs.TimeBankResponseDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// On attend +2 (12h faites - 10h attendues)
	assert.Equal(t, 2, response.Balance, "Le solde devrait être de 2 heures")
}

// Helper local pour simplifier la création de dates
func createDate(baseDate time.Time, hour, min int) time.Time {
	return time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), hour, min, 0, 0, baseDate.Location())
}
