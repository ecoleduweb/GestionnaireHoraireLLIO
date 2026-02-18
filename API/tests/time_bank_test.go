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

func TestTimeBankScenarios(t *testing.T) {
	// --- 1. Préparation des dates (Semaine dernière) ---
	loc := time.Local
	now := time.Now().In(loc)

	// Lundi de la semaine courante
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	mondayCurrentWeek := now.AddDate(0, 0, -(weekday - 1))

	// Lundi de la semaine DERNIÈRE (Début période) - Mise à 00:00:00
	mondayLastWeek := time.Date(mondayCurrentWeek.Year(), mondayCurrentWeek.Month(), mondayCurrentWeek.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -7)
	startDateStr := mondayLastWeek.Format("2006-01-02")

	// --- 2. Nettoyage initial ---
	// On s'assure que l'utilisateur de test est propre
	database.DB.Model(&doNotDeleteUser).Updates(map[string]interface{}{
		"time_bank_start_date":     nil,
		"time_bank_hours_per_week": nil,
		"time_bank_balance_offset": 0,
	})

	// --- 3. Définition du Nettoyage AUTOMATIQUE (Correction du bug de defer) ---
	// t.Cleanup s'exécute TOUJOURS à la fin du test, même en cas de Fatalf ou Panic
	t.Cleanup(func() {
		database.DB.Where("name LIKE ?", "Test TimeBank%").Delete(&DAOs.Activity{})
		// On remet l'user au propre
		database.DB.Model(&doNotDeleteUser).Updates(map[string]interface{}{
			"time_bank_start_date":     nil,
			"time_bank_hours_per_week": nil,
		})
	})

	// --- SCÉNARIO 1 : Date Invalide lors de la config ---
	t.Run("InvalidDate", func(t *testing.T) {
		badConfig := map[string]interface{}{
			"startDate":    "99-99-2024", // Format invalide
			"hoursPerWeek": 35,
		}
		w := sendRequest(router, "POST", "/user/time-bank/config", badConfig, enums.Employee)
		// Doit échouer (Bad Request)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// --- SCÉNARIO 2 : Pas d'activités (Déficit Total) ---
	t.Run("NoActivities_Deficit", func(t *testing.T) {
		// Config : 10h/semaine
		config := DTOs.TimeBankConfigDTO{StartDate: startDateStr, HoursPerWeek: 10, Offset: 0}
		sendRequest(router, "POST", "/user/time-bank/config", config, enums.Employee)

		w := sendRequest(router, "GET", "/user/time-bank", nil, enums.Employee)
		assertResponse(t, w, http.StatusOK, nil)

		var resp DTOs.TimeBankBalanceDTO
		json.Unmarshal(w.Body.Bytes(), &resp)

		// 0h fait - 10h attendu = -10
		assert.Equal(t, -10, *resp.TimeInBank)
	})

	// --- Injection des activités pour la suite (12 heures travaillées) ---
	// On injecte MAINTENANT, après avoir testé le cas "sans activités"
	activities := []DAOs.Activity{
		{
			Name:        "Test TimeBank 1",
			Description: "Auto",
			StartDate:   mondayLastWeek.Add(9 * time.Hour),
			EndDate:     mondayLastWeek.Add(13 * time.Hour), // 4h
			UserId:      doNotDeleteUser.Id,
			ProjectId:   doNotDeleteProject.Id,
			CategoryId:  doNotDeleteCategory.Id,
		},
		{
			Name:        "Test TimeBank 2",
			Description: "Auto",
			StartDate:   mondayLastWeek.AddDate(0, 0, 1).Add(9 * time.Hour),
			EndDate:     mondayLastWeek.AddDate(0, 0, 1).Add(13 * time.Hour), // 4h
			UserId:      doNotDeleteUser.Id,
			ProjectId:   doNotDeleteProject.Id,
			CategoryId:  doNotDeleteCategory.Id,
		},
		{
			Name:        "Test TimeBank 3",
			Description: "Auto",
			StartDate:   mondayLastWeek.AddDate(0, 0, 2).Add(9 * time.Hour),
			EndDate:     mondayLastWeek.AddDate(0, 0, 2).Add(13 * time.Hour), // 4h
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

	// --- SCÉNARIO 3 : Surplus (Fait > Attendu) ---
	t.Run("Surplus", func(t *testing.T) {
		// Config : 10h/semaine (On a fait 12h)
		config := DTOs.TimeBankConfigDTO{StartDate: startDateStr, HoursPerWeek: 10, Offset: 0}
		sendRequest(router, "POST", "/user/time-bank/config", config, enums.Employee)

		w := sendRequest(router, "GET", "/user/time-bank", nil, enums.Employee)
		var resp DTOs.TimeBankBalanceDTO
		json.Unmarshal(w.Body.Bytes(), &resp)

		// 12h - 10h = +2
		assert.Equal(t, 2, *resp.TimeInBank)
	})

	// --- SCÉNARIO 4 : Déficit (Fait < Attendu) ---
	t.Run("Deficit", func(t *testing.T) {
		// Config : 20h/semaine (On a fait 12h)
		config := DTOs.TimeBankConfigDTO{StartDate: startDateStr, HoursPerWeek: 20, Offset: 0}
		sendRequest(router, "POST", "/user/time-bank/config", config, enums.Employee)

		w := sendRequest(router, "GET", "/user/time-bank", nil, enums.Employee)
		var resp DTOs.TimeBankBalanceDTO
		json.Unmarshal(w.Body.Bytes(), &resp)

		// 12h - 20h = -8
		assert.Equal(t, -8, *resp.TimeInBank)
	})

	// --- SCÉNARIO 5 : Solde Nul (Fait == Attendu) ---
	t.Run("ZeroBalance", func(t *testing.T) {
		// Config : 12h/semaine (On a fait 12h)
		config := DTOs.TimeBankConfigDTO{StartDate: startDateStr, HoursPerWeek: 12, Offset: 0}
		sendRequest(router, "POST", "/user/time-bank/config", config, enums.Employee)

		w := sendRequest(router, "GET", "/user/time-bank", nil, enums.Employee)
		var resp DTOs.TimeBankBalanceDTO
		json.Unmarshal(w.Body.Bytes(), &resp)

		// 12h - 12h = 0
		assert.Equal(t, 0, *resp.TimeInBank)
	})

}
