package tests

import (
	"encoding/json"
	"llio-api/database"
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"llio-api/models/enums"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeBankWorkflow(t *testing.T) {
	// --- 1. Préparation des dates (Semaine dernière) ---
	// On se cale sur le fuseau horaire LOCAL pour correspondre à la logique du service
	loc := time.Local
	now := time.Now().In(loc)

	// Trouver le lundi de la semaine courante
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	} // Dimanche = 7
	mondayCurrentWeek := now.AddDate(0, 0, -(weekday - 1))

	// Lundi de la semaine D'AVANT (Date de début du test)
	// On force l'heure à 00:00:00 pour être propre
	mondayLastWeek := time.Date(mondayCurrentWeek.Year(), mondayCurrentWeek.Month(), mondayCurrentWeek.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -7)
	startDateStr := mondayLastWeek.Format("2006-01-02")

	// --- 2. Nettoyage initial (Au cas où) ---
	// On s'assure que l'utilisateur n'a pas de config résiduelle
	database.DB.Model(&doNotDeleteUser).Updates(map[string]interface{}{
		"time_bank_start_date":     nil,
		"time_bank_hours_per_week": nil,
		"time_bank_balance_offset": 0,
	})

	// --- TEST A : Vérifier que c'est "Non Configuré" au départ ---
	t.Run("CheckUnconfigured", func(t *testing.T) {
		w := sendRequest(router, "GET", "/user/time-bank", nil, enums.Employee)
		assertResponse(t, w, http.StatusOK, nil)

		var resp DTOs.TimeBankBalanceDTO
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.False(t, resp.IsConfigured, "Devrait être non configuré au départ")
		assert.Nil(t, resp.TimeInBank)
	})

	// --- TEST B : Sauvegarder la configuration ---
	t.Run("SaveConfig", func(t *testing.T) {
		configPayload := DTOs.TimeBankConfigDTO{
			StartDate:    startDateStr,
			HoursPerWeek: 10, // On attend 10h/semaine
			Offset:       0,
		}

		w := sendRequest(router, "POST", "/user/time-bank/config", configPayload, enums.Employee)
		assertResponse(t, w, http.StatusOK, nil)
	})

	// --- TEST C : Créer des activités (12h travaillées la semaine passée) ---
	// On insère directement en BD pour contourner les règles API éventuelles et forcer les dates
	activities := []DAOs.Activity{
		{
			Name:        "Test TimeBank 1",
			Description: "Auto-generated",
			StartDate:   mondayLastWeek.Add(9 * time.Hour),  // Lundi 9h
			EndDate:     mondayLastWeek.Add(13 * time.Hour), // Lundi 13h (4h)
			UserId:      doNotDeleteUser.Id,
			ProjectId:   doNotDeleteProject.Id,
			CategoryId:  doNotDeleteCategory.Id,
		},
		{
			Name:        "Test TimeBank 2",
			Description: "Auto-generated",
			StartDate:   mondayLastWeek.AddDate(0, 0, 1).Add(9 * time.Hour),  // Mardi 9h
			EndDate:     mondayLastWeek.AddDate(0, 0, 1).Add(13 * time.Hour), // Mardi 13h (4h)
			UserId:      doNotDeleteUser.Id,
			ProjectId:   doNotDeleteProject.Id,
			CategoryId:  doNotDeleteCategory.Id,
		},
		{
			Name:        "Test TimeBank 3",
			Description: "Auto-generated",
			StartDate:   mondayLastWeek.AddDate(0, 0, 2).Add(9 * time.Hour),  // Mercredi 9h
			EndDate:     mondayLastWeek.AddDate(0, 0, 2).Add(13 * time.Hour), // Mercredi 13h (4h)
			UserId:      doNotDeleteUser.Id,
			ProjectId:   doNotDeleteProject.Id,
			CategoryId:  doNotDeleteCategory.Id,
		},
	}

	for _, act := range activities {
		if err := database.DB.Create(&act).Error; err != nil {
			t.Fatalf("Erreur insert activité: %v", err)
		}
	}

	// Nettoyage automatique à la fin du test
	defer database.DB.Where("name LIKE ?", "Test TimeBank%").Delete(&DAOs.Activity{})

	// --- TEST D : Vérifier le solde ---
	t.Run("CheckBalance", func(t *testing.T) {
		w := sendRequest(router, "GET", "/user/time-bank", nil, enums.Employee)
		assertResponse(t, w, http.StatusOK, nil)

		var resp DTOs.TimeBankBalanceDTO
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.True(t, resp.IsConfigured)
		assert.NotNil(t, resp.TimeInBank)

		// Calcul : 12h faites - 10h attendues = +2h
		assert.Equal(t, 2, *resp.TimeInBank, "Le solde devrait être de +2h")
	})
}
